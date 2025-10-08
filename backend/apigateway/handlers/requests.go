package handlers

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "strings"
    "os"
    "sync"
    "time"
    "log"
    "encoding/base64"
	"github.com/gin-gonic/gin"
)

var (
    gestionURL          string
    parserURL           string
    documentacionURL    string
    oauth2URL           string
	httpClient          *http.Client
)

func init() {
	gestionURL = os.Getenv("GESTION_URL")
    parserURL = os.Getenv("PARSER_URL")
    documentacionURL = os.Getenv("DOCUMENTATION_URL")
    oauth2URL = os.Getenv("OAUTH2_URL")
    httpClient = &http.Client{
        Timeout: 15 * time.Second,
    }
    // maxAgeCookie := 3600 // 1 hora, 60*60 segundos
}

func obtenerEdificiosUsuario(email string) []interface{} {
    if email == "" {
        log.Println("Email vacío, no se pueden obtener edificios")
        return nil
    }

    // Construir URL para obtener edificios
    url := fmt.Sprintf("%s/usuarios/edificios/%s", gestionURL, email)
    
    resp, err := http.Get(url)
    if err != nil {
        log.Printf("Error obteniendo edificios del usuario: %v", err)
        return nil
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        log.Printf("Error en respuesta de edificios: status %d", resp.StatusCode)
        return nil
    }

    // Leer respuesta
    var edificiosResponse map[string]interface{}
    if err := json.NewDecoder(resp.Body).Decode(&edificiosResponse); err != nil {
        log.Printf("Error decodificando edificios: %v", err)
        return nil
    }

    // Extraer solo el array de edificios
    if edificios, exists := edificiosResponse["edificios"].([]interface{}); exists {
        return edificios
    }

    return nil
}

func extractClaimFromToken(tokenString string, claimKey string) (interface{}, error) {
    // Un JWT tiene 3 partes separadas por puntos: header.payload.signature
    parts := strings.Split(tokenString, ".")
    if len(parts) != 3 {
        return nil, fmt.Errorf("token inválido")
    }

    // Decodificar el payload (segunda parte)
    payload := parts[1]
    
    // Agregar padding si es necesario
    if l := len(payload) % 4; l > 0 {
        payload += strings.Repeat("=", 4-l)
    }

    // Decodificar de base64
    decoded, err := base64.URLEncoding.DecodeString(payload)
    if err != nil {
        return nil, fmt.Errorf("error decodificando payload: %v", err)
    }

    // Parsear JSON
    var claims map[string]interface{}
    if err := json.Unmarshal(decoded, &claims); err != nil {
        return nil, fmt.Errorf("error parseando claims: %v", err)
    }

    // Extraer el claim solicitado
    if value, ok := claims[claimKey]; ok {
        return value, nil
    }

    return nil, fmt.Errorf("%s no encontrado en token", claimKey)
}

// ActivosCompletosHandler maneja la petición consolidada de activos con sensores
func ActivosCompletosHandler(c *gin.Context) {
    // tiposActivos := []string{"caldera"}
	tiposActivos := []string{"bomba de agua", "caldera", "ascensor", "transformador"}
    
    var wg sync.WaitGroup
    var mu sync.Mutex
    
    client := &http.Client{Timeout: 15 * time.Second}
    todosLosActivos := []interface{}{}
    fmt.Println("DEBUG: Llegué hasta aquí - iniciando obtención de activos por tipo")
    
    // 1. Obtener activos por tipo (4 peticiones concurrentes)
    for _, tipo := range tiposActivos {
        wg.Add(1)
        go func(t string) {
            defer wg.Done()
            
            url := fmt.Sprintf("%s/activos/tipo/%s", gestionURL, t)
            resp, err := client.Get(url)
            if err != nil {
                return
            }
            defer resp.Body.Close()
            
            var data map[string]interface{}
            if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
                return
            }
            
            if activos, exists := data["activos"]; exists {
                if activosArray, ok := activos.([]interface{}); ok {
                    mu.Lock()
                    todosLosActivos = append(todosLosActivos, activosArray...)
                    mu.Unlock()
                }
            }
        }(tipo)
    }
    
    wg.Wait()
    
    // 2. Para cada activo, hacer las 3 peticiones al parser
    activosCompletos := []interface{}{}
    
    for _, activo := range todosLosActivos {
        if activoMap, ok := activo.(map[string]interface{}); ok {
            activoID := fmt.Sprintf("%v", activoMap["id"])
            
            wg.Add(1)
            go func(aID string, aMap map[string]interface{}) {
                defer wg.Done()
                
                activoCompleto := make(map[string]interface{})
                // Copiar datos base del activo
                for k, v := range aMap {
                    activoCompleto[k] = v
                }
                
                var wgActivo sync.WaitGroup
                var muActivo sync.Mutex
                
                // Petición sensores/estado
                wgActivo.Add(1)
                go func() {
                    defer wgActivo.Done()
                    url := fmt.Sprintf("%s/activo/%s/sensores/estado", parserURL, aID)
                    resp, err := client.Get(url)
                    if err != nil {
                        return
                    }
                    defer resp.Body.Close()
                    
                    var sensoresData interface{}
                    if err := json.NewDecoder(resp.Body).Decode(&sensoresData); err != nil {
                        return
                    }
                    
                    muActivo.Lock()
                    activoCompleto["sensores_estado"] = sensoresData
                    muActivo.Unlock()
                }()
                
                // Petición últimos valores
                wgActivo.Add(1)
                go func() {
                    defer wgActivo.Done()
                    url := fmt.Sprintf("%s/lectura/%s/datos/ultimo", parserURL, aID)
                    resp, err := client.Get(url)
                    if err != nil {
                        return
                    }
                    defer resp.Body.Close()
                    
                    var ultimosData interface{}
                    if err := json.NewDecoder(resp.Body).Decode(&ultimosData); err != nil {
                        return
                    }
                    
                    muActivo.Lock()
                    activoCompleto["ultimos_valores"] = ultimosData
                    muActivo.Unlock()
                }()
                
                // Petición datos históricos
                wgActivo.Add(1)
                go func() {
                    defer wgActivo.Done()
                    url := fmt.Sprintf("%s/lectura/%s/datos", parserURL, aID)
                    resp, err := client.Get(url)
                    if err != nil {
                        return
                    }
                    defer resp.Body.Close()
                    
                    var historicosData interface{}
                    if err := json.NewDecoder(resp.Body).Decode(&historicosData); err != nil {
                        return
                    }
                    
                    muActivo.Lock()
                    activoCompleto["datos_historicos"] = historicosData
                    muActivo.Unlock()
                }()
                
                wgActivo.Wait()
                
                mu.Lock()
                activosCompletos = append(activosCompletos, activoCompleto)
                mu.Unlock()
                
            }(activoID, activoMap)
        }
    }
    
    wg.Wait()
    
    result := map[string]interface{}{
        "activos": activosCompletos,
        "total": len(activosCompletos),
        "timestamp": time.Now(),
    }
    
    c.JSON(http.StatusOK, result)
}

func ActivosYSensores(c *gin.Context) {
	// c.JSON(http.StatusOK, gin.H{"message": "Sensores endpoint"})
	var wg sync.WaitGroup
    var mu sync.Mutex

	// tiposActivos := []string{"bomba de agua"}
    tiposActivos := []string{"transformador", "bomba de agua", "caldera", "ascensor"}
	activos := []interface{}{}
	
	for _, tipo := range tiposActivos {
		
		fmt.Println("DEBUG: Obteniendo activos de tipo", tipo)
		wg.Add(1) // Aumentar contador del WaitGroup
		go func(t string) {
			defer wg.Done() // Marca la tarea como terminada al finalizar la goroutine
			url := fmt.Sprintf("%s/activos/tipo/%s", gestionURL, t)
			resp, err := httpClient.Get(url)
			if err != nil {
				fmt.Println("Error: ", err)
				return
			}
			defer resp.Body.Close()
			
			var data map[string]interface{} // Mapa para decodificar la respuesta JSON
			if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
				fmt.Println("Error decodificando JSON: ", err)
				return
			}
            // Añadir lo que está en "activos" al slice principal
			if activosData, exists := data["activos"]; exists {
                if activosArray, ok := activosData.([]interface{}); ok {
                    mu.Lock()
                    activos = append(activos, activosArray...)
                    mu.Unlock()
                }
            }
		}(tipo)
	}

	wg.Wait() // Esperar a que termine la goroutine para ejecutar la siguiente parte

    // Hacer la segunda parte: obtener sensores para cada activo
    for _, activo := range activos {
        // activo es una interface del estilo, que hay que convertir para extraer la id:
        // map[edificio_id:1 estado:operativo id:2 nombre:Bomba Centrífuga A tipo:bomba de agua]
        fmt.Println("DEBUG: Obteniendo sensores para activo", activo)
        if activoMap, ok := activo.(map[string]interface{}); ok {
            activoID := fmt.Sprintf("%v", activoMap["id"]) // Convertir a string la id obtenida
            wg.Add(1)
            go func(aMap map[string]interface{}) {
                defer wg.Done() // Ejecutar al finalizar la goroutine
                url := fmt.Sprintf("%s/activo/%s/sensores/estado", parserURL, activoID)
                resp, err := httpClient.Get(url)
                if err != nil {
                    fmt.Println("Error obteniendo sensores: ", err)
                    return
                }
                defer resp.Body.Close()

                body := make(map[string]interface{})
                if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
                    fmt.Println("Error decodificando sensores: ", err)
                    return
                }

                if sensores, exists := body["sensores"]; exists {
                    mu.Lock()
                    // Añadir sensores al mapa del activo
                    // como se pasó como referencia, se modifica el original en el slice.
                    aMap["sensores"] = sensores 
                    mu.Unlock()
                    fmt.Printf("DEBUG: Sensores agregados para activo ID %s\n", activoID)

                    // Obtener datos históricos para el activo completo
                    wg.Add(1)
                    go func(aID string, sensoresList interface{}) {
                        defer wg.Done()
                        fmt.Println("DEBUG: Obteniendo datos históricos para activo id: ", aID)
                        url := fmt.Sprintf("%s/lectura/%s/datos", parserURL, aID)
                        resp, err := httpClient.Get(url)
                        if err != nil {
                            fmt.Println("Error obteniendo datos históricos: ", err)
                            return
                        }
                        defer resp.Body.Close()

                        var datosHistoricos map[string]interface{}
                        if err := json.NewDecoder(resp.Body).Decode(&datosHistoricos); err != nil {
                            fmt.Println("Error decodificando datos históricos: ", err)
                            return
                        }

                        // Distribuir los datos históricos a cada sensor
                        if sensoresData, exists := datosHistoricos["sensores"]; exists {
                            if sensoresArray, ok := sensoresData.([]interface{}); ok {
                                // Crear un mapa para acceso rápido por sensor_id
                                datosPorSensor := make(map[string]interface{})
                                for _, sensorData := range sensoresArray {
                                    if sensorMap, ok := sensorData.(map[string]interface{}); ok {
                                        if sensorID, exists := sensorMap["sensor_id"]; exists {
                                            datosPorSensor[fmt.Sprintf("%v", sensorID)] = sensorMap["datos"]
                                            // fmt.Printf("DEBUG: Datos encontrados para sensor %s\n", sensorID)
                                        }
                                    }
                                }

                                // fmt.Printf("DEBUG: Total sensores con datos: %d\n", len(datosPorSensor))

                                // Agregar los datos a cada sensor en la lista original
                                mu.Lock()
                                if sensoresOriginales, ok := sensoresList.([]interface{}); ok {
                                    fmt.Printf("DEBUG: Distribuyendo datos a sensores del activo id: %s\n", aID)
                                    for _, sensor := range sensoresOriginales {
                                        if sensorMap, ok := sensor.(map[string]interface{}); ok {
                                            // Usar 'sensor_id' en lugar de 'id'
                                            if sensorID, exists := sensorMap["sensor_id"]; exists {
                                                sensorIDStr := fmt.Sprintf("%v", sensorID)
                                                // fmt.Printf("DEBUG: Buscando datos para sensor %s\n", sensorIDStr)
                                                if datos, encontrado := datosPorSensor[sensorIDStr]; encontrado {
                                                    sensorMap["datos"] = datos
                                                    // fmt.Printf("DEBUG: Datos agregados al sensor %s\n", sensorIDStr)
                                                } else {
                                                    fmt.Printf("DEBUG: No se encontraron datos para sensor %s\n", sensorIDStr)
                                                }
                                            }
                                        }
                                    }
                                }
                                mu.Unlock()
                            }
                        }
                    }(activoID, sensores)
                }

            }(activoMap) // Pasar el mapa del activo a la goroutine
        }
    }

    wg.Wait() // Esperar a que terminen todas las goroutines


	result := map[string]interface{}{
        "activos": activos,
        "total": len(activos),
        "timestamp": time.Now(),
    }

	c.JSON(http.StatusOK, result)
}

func ObtenerActivos(c *gin.Context) {
    var wg sync.WaitGroup
    var mu sync.Mutex
    activos := []interface{}{}
    
    wg.Add(1)
    // Goroutine para obtener activos
    go func() {
        defer wg.Done()
        url := fmt.Sprintf("%s/activos", gestionURL)
        resp, err := httpClient.Get(url)
        if err != nil {
            fmt.Println("Error obteniendo activos: ", err)
            return
        }
        defer resp.Body.Close()
        
        var data map[string]interface{}
        if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
            fmt.Println("Error decodificando activos: ", err)
            return
        }
        
        if activosData, exists := data["activos"]; exists {
            if activosArray, ok := activosData.([]interface{}); ok {
                mu.Lock()
                activos = activosArray
                mu.Unlock()
            }
        }
    }()

    wg.Wait() // Esperar a que termine la goroutine de activos

    result := map[string]interface{}{
        "activos": activos,
    }

    c.JSON(http.StatusOK, result)
}

func GenerarReporte(c *gin.Context) {
    // NO terminado
    id := c.Param("id_reporte")
    if id == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "ID del activo es requerido"})
        return
    }

    result := map[string]interface{}{
        "reporte": fmt.Sprintf("Reporte generado para activo ID %s", id),
    }
    c.JSON(http.StatusOK, result)
}

func ObtenerActivoPorID(c *gin.Context) {
    var wg sync.WaitGroup
    var mu sync.Mutex
    id := c.Param("id_activo")
    if id == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "ID del activo es requerido"})
        return
    }
    
    var activo map[string]interface{}
    wg.Add(1)
    go func() {
        defer wg.Done()
        url := fmt.Sprintf("%s/activo/%s", parserURL, id)
        resp, err := httpClient.Get(url)
        if err != nil {
            fmt.Println("Error obteniendo activo: ", err)
            return
        }
        defer resp.Body.Close()

        var data map[string]interface{}
        if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
            fmt.Println("Error decodificando activo: ", err)
            return
        }
        mu.Lock()
        activo = data
        mu.Unlock()
    }()

    wg.Wait()

    // Obtener los datos restantes del activo desde el servicio de gestión
    wg.Add(1)
    go func() {
        defer wg.Done()
        url := fmt.Sprintf("%s/activos/%s", gestionURL, id)
        resp, err := httpClient.Get(url)
        if err != nil {
            fmt.Println("Error obteniendo datos de gestión: ", err)
            return
        }
        defer resp.Body.Close()

        var data map[string]interface{}
        if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
            fmt.Println("Error decodificando datos de gestión: ", err)
            return
        }
        // Extraer solo la ubicación y agregarla al activo
        if ubicacion, exists := data["ubicacion"]; exists {
            mu.Lock()
            if activo != nil {
                activo["ubicacion"] = ubicacion
            }
            mu.Unlock()
        }
    }()

    wg.Wait()

    // Buscar ID ficha técnica
    wg.Add(1)
    go func() {
        defer wg.Done()
        url := fmt.Sprintf("%s/api/v1/documentos/activo/%s/ficha-tecnica", documentacionURL, id)
        fmt.Printf("DEBUG: URL ficha técnica: %s\n", url)
        resp, err := httpClient.Get(url)
        if err != nil {
            fmt.Println("Error obteniendo ficha técnica: ", err)
            return
        }
        defer resp.Body.Close()

        // fmt.Printf("DEBUG: Respuesta ficha técnica: %+v\n", resp)
        
        var data map[string]interface{}
        if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
            fmt.Println("Error decodificando ficha técnica: ", err)
            return
        }

        //Extraer ID de la ficha técnica y agregarla al activo
        if fichaID, exists := data["id"]; exists {
            mu.Lock()
            if activo != nil {
                activo["id_ficha_tecnica"] = fichaID
            }
            mu.Unlock()
        }
    }()

    wg.Wait()    

    // Verificar si se obtuvo el activo
    if activo == nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Activo no encontrado"})
        return
    }

    c.JSON(http.StatusOK, activo)
}

func LoginHandler(c *gin.Context) {
    // var wg sync.WaitGroup
    // var mu sync.Mutex

    // Leer el body de la request
    var loginData map[string]interface{}
    if err := c.ShouldBindJSON(&loginData); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
        return
    }
    fmt.Printf("Login data received: %+v\n", loginData)

    // Convertir a JSON para enviar al microservicio OAuth2
    jsonData, err := json.Marshal(loginData)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error processing request"})
        return
    }

    // Hacer la petición POST al microservicio OAuth2
    url := fmt.Sprintf("%s/login", oauth2URL)
    // fmt.Println("DEBUG: URL de login OAuth2:", url, oauth2URL)
    resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
    if err != nil {
        fmt.Println("Error llamando a OAuth2: ", err)
        c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Authentication service unavailable"})
        return
    }
    defer resp.Body.Close()

    // Leer respuesta del microservicio OAuth2
    var response map[string]interface{}
    if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
        fmt.Println("Error decodificando respuesta OAuth2: ", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error processing response"})
        return
    }

    // Si el login es exitoso, obtener los edificios del usuario
    if resp.StatusCode == http.StatusOK {
        email, ok := loginData["email"].(string)
        if !ok {
            if user, exists := response["user"].(map[string]interface{}); exists {
                if userEmail, hasEmail := user["email"].(string); hasEmail {
                    email = userEmail
                }
            }
        }

        edificios := obtenerEdificiosUsuario(email)

        // Filtrar la respuesta
        filteredResponse := make(map[string]interface{})
        
        // Copiar solo los campos que queremos exponer
        if accessToken, exists := response["access_token"]; exists {
            filteredResponse["access_token"] = accessToken
        }

        if refreshToken, exists := response["refresh_token"]; exists {
            filteredResponse["refresh_token"] = refreshToken
        }

        if edificios != nil {
            filteredResponse["edificios"] = edificios
        }

        c.JSON(resp.StatusCode, filteredResponse)
        return
    }

    // En caso de error en login, devolver el mensaje de error
    c.JSON(resp.StatusCode, response)
    return
}

func LogoutHandler(c *gin.Context) {
    // Extraer token del header Authorization
    authHeader := c.GetHeader("Authorization")
    if authHeader == "" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
        return
    }

    // Verificar formato del token Bearer
    tokenParts := strings.Split(authHeader, " ")
    if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
        return
    }

    accessToken := tokenParts[1]

    fmt.Printf("Logout token: %s\n", accessToken)

     // Preparar la request al microservicio OAuth2
    url := fmt.Sprintf("%s/api/logout", oauth2URL)
    
    // Crear request vacía o con body mínimo
    req, err := http.NewRequest("POST", url, nil)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating request"})
        return
    }

    // IMPORTANTE: Enviar el token en el HEADER, no en el body
    req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
    req.Header.Set("Content-Type", "application/json")

    // Ejecutar la petición
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        fmt.Println("Error llamando a OAuth2 para logout: ", err)
        c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Authentication service unavailable"})
        return
    }
    defer resp.Body.Close()

    fmt.Println("DEBUG: Logout request sent to OAuth2, status:", resp.StatusCode)

    // Leer respuesta del microservicio OAuth2
    var response map[string]interface{}
    if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
        fmt.Println("Error decodificando respuesta OAuth2 de logout: ", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error processing response"})
        return
    }

    // Devolver la respuesta del microservicio OAuth2 directamente
    c.JSON(resp.StatusCode, response)
}

func RefreshHandler(c *gin.Context) {
    // Leer el body con el refresh_token
    var refreshData map[string]interface{}
    if err := c.ShouldBindJSON(&refreshData); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
        return
    }

    // Verificar que venga el refresh_token
    if _, exists := refreshData["refresh_token"]; !exists {
        c.JSON(http.StatusBadRequest, gin.H{"error": "refresh_token is required"})
        return
    }

    fmt.Printf("Refresh data received: %+v\n", refreshData)

    // Convertir a JSON para enviar al microservicio OAuth2
    jsonData, err := json.Marshal(refreshData)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error processing request"})
        return
    }

    // Hacer la petición POST al microservicio OAuth2
    url := fmt.Sprintf("%s/refresh", oauth2URL)
    resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
    if err != nil {
        fmt.Println("Error llamando a OAuth2 para refresh: ", err)
        c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Authentication service unavailable"})
        return
    }
    defer resp.Body.Close()

    fmt.Println("DEBUG: Refresh request sent to OAuth2, status:", resp.StatusCode)

    // Leer respuesta del microservicio OAuth2, que sería el nuevo access_token y refresh_token
    var response map[string]interface{}
    if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
        fmt.Println("Error decodificando respuesta OAuth2 de refresh: ", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error processing response"})
        return
    }

    // Obtener el email del usuario a partir del access_token recibido
    

    // Si el refresh es exitoso, obtener los edificios del usuario
    if resp.StatusCode == http.StatusOK {
        var email string
        if accessToken, exists := response["access_token"].(string); exists {
            emailInterface, err := extractClaimFromToken(accessToken, "email")
            if err != nil {
                fmt.Println("Error extrayendo email del token: ", err)
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Error processing token"})
                return
            }
            
            // ✅ Type assertion para convertir interface{} a string
            var ok bool
            email, ok = emailInterface.(string)
            if !ok {
                fmt.Println("Error: email no es un string")
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid email format in token"})
                return
            }
        }

        fmt.Printf("Email extraído del token: %s\n", email)

        // Obtener edificios del usuario
        edificios := obtenerEdificiosUsuario(email)

        // Filtrar la respuesta
        filteredResponse := make(map[string]interface{})
        
        // Copiar solo los campos que queremos exponer
        if accessToken, exists := response["access_token"]; exists {
            filteredResponse["access_token"] = accessToken
        }

        if refreshToken, exists := response["refresh_token"]; exists {
            filteredResponse["refresh_token"] = refreshToken
        }

        if edificios != nil {
            filteredResponse["edificios"] = edificios
        }

        // Devolver la respuesta filtrada con edificios
        c.JSON(resp.StatusCode, filteredResponse)
        return
    }

    // Si hubo error en el refresh, retornar el error original
    c.JSON(resp.StatusCode, response)
}

func ObtenerRol(c *gin.Context) {
    // Extraer token del header Authorization
    authHeader := c.GetHeader("Authorization")
    if authHeader == "" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Token requerido"})
        return
    }

    // Verificar formato del token Bearer
    tokenParts := strings.Split(authHeader, " ")
    if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Formato de token inválido"})
        return
    }

    tokenString := tokenParts[1]

    scope, err := extractClaimFromToken(tokenString, "scope")
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
        return
    }

    // Retornar solo el scope
    c.JSON(http.StatusOK, gin.H{
        "scope": scope,
    })
}