package handlers

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "mime/multipart"
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
            
            // Type assertion para convertir interface{} a string
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

// Obtener Activos por id_edificio
func ObtenerActivosPorEdificio(c *gin.Context) {
    idEdificio := c.Param("id_edificio") // Parámetro opcional para filtrar por edificio
    if idEdificio == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "ID del edificio es requerido"})
        return
    }

    var wg sync.WaitGroup
    var mu sync.Mutex
    activos := []interface{}{}
    estadosMap := make(map[int]string) // Agregar declaración del mapa de estados
    sensoresMap := make(map[int]interface{})

    // Verificar si se solicitan sensores
    incluirSensores := c.Query("sensores") == "true"
    
    wg.Add(1)
    // Goroutine para obtener activos del edificio
    go func() {
        defer wg.Done()
        url := fmt.Sprintf("%s/activos/edificio/%s", gestionURL, idEdificio)

        resp, err := httpClient.Get(url)
        if err != nil {
            fmt.Println("Error obteniendo activos por edificio: ", err)
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

    // Si se solicitan sensores, consultar individualmente cada activo en el parser
    if incluirSensores {
        fmt.Println("DEBUG: Consultando sensores para cada activo")
        
        for _, activo := range activos {
            if activoMap, ok := activo.(map[string]interface{}); ok {
                if idInterface, exists := activoMap["id"]; exists {
                    var activoID int
                    
                    // Convertir el ID a int
                    switch v := idInterface.(type) {
                    case float64:
                        activoID = int(v)
                    case int:
                        activoID = v
                    default:
                        fmt.Printf("Tipo inesperado para id de activo: %T\n", idInterface)
                        continue
                    }

                    // Consultar sensores y datos para este activo específico
                    wg.Add(1)
                    go func(aID int) {
                        defer wg.Done()
                        
                        // Consultar información de sensores (estado, tipo, unidad, etc.)
                        url := fmt.Sprintf("%s/activo/%d", parserURL, aID)
                        resp, err := httpClient.Get(url)
                        if err != nil {
                            fmt.Printf("Error obteniendo sensores para activo %d: %v\n", aID, err)
                            return
                        }
                        defer resp.Body.Close()

                        var activoData map[string]interface{}
                        if err := json.NewDecoder(resp.Body).Decode(&activoData); err != nil {
                            fmt.Printf("Error decodificando datos del activo %d: %v\n", aID, err)
                            return
                        }

                        mu.Lock()
                        // Guardar estado si existe
                        if estado, exists := activoData["estado"]; exists {
                            if estadoStr, ok := estado.(string); ok {
                                estadosMap[aID] = estadoStr
                            }
                        }

                        // Obtener sensores
                        var sensoresArray []interface{}
                        if sensores, exists := activoData["sensores"]; exists {
                            if sensoresArr, ok := sensores.([]interface{}); ok {
                                sensoresArray = sensoresArr
                            }
                        }
                        mu.Unlock()

                        // Consultar datos históricos de sensores
                        urlDatos := fmt.Sprintf("%s/lectura/%d/datos", parserURL, aID)
                        respDatos, err := httpClient.Get(urlDatos)
                        if err != nil {
                            fmt.Printf("Error obteniendo datos para activo %d: %v\n", aID, err)
                            // Guardar sensores sin datos
                            mu.Lock()
                            sensoresMap[aID] = sensoresArray
                            mu.Unlock()
                            return
                        }
                        defer respDatos.Body.Close()

                        var datosResponse map[string]interface{}
                        if err := json.NewDecoder(respDatos.Body).Decode(&datosResponse); err != nil {
                            fmt.Printf("Error decodificando datos históricos del activo %d: %v\n", aID, err)
                            // Guardar sensores sin datos
                            mu.Lock()
                            sensoresMap[aID] = sensoresArray
                            mu.Unlock()
                            return
                        }

                        // Crear un mapa de datos por sensor_id
                        datosMap := make(map[string]interface{})
                        if sensoresDatos, exists := datosResponse["sensores"]; exists {
                            if sensoresDatosArray, ok := sensoresDatos.([]interface{}); ok {
                                for _, sensorDato := range sensoresDatosArray {
                                    if sensorDatoMap, ok := sensorDato.(map[string]interface{}); ok {
                                        if sensorID, exists := sensorDatoMap["sensor_id"]; exists {
                                            if datos, existsDatos := sensorDatoMap["datos"]; existsDatos {
                                                datosMap[fmt.Sprintf("%v", sensorID)] = datos
                                            }
                                        }
                                    }
                                }
                            }
                        }

                        // Agregar los datos a cada sensor
                        for _, sensor := range sensoresArray {
                            if sensorMap, ok := sensor.(map[string]interface{}); ok {
                                if sensorID, exists := sensorMap["sensor_id"]; exists {
                                    sensorIDStr := fmt.Sprintf("%v", sensorID)
                                    if datos, encontrado := datosMap[sensorIDStr]; encontrado {
                                        sensorMap["datos"] = datos
                                    } else {
                                        sensorMap["datos"] = nil
                                    }
                                }
                            }
                        }

                        // Guardar sensores con datos
                        mu.Lock()
                        sensoresMap[aID] = sensoresArray
                        mu.Unlock()

                    }(activoID)
                }
            }
        }

        wg.Wait() // Esperar a que terminen todas las consultas de sensores
    } else {
        // Si NO se solicitan sensores, usar el endpoint para estados del edificio
        wg.Add(1)
        go func() {
            defer wg.Done()
            url := fmt.Sprintf("%s/activo/edificio/%s", parserURL, idEdificio)

            resp, err := httpClient.Get(url)
            if err != nil {
                fmt.Println("Error obteniendo estados de activos: ", err)
                return
            }
            defer resp.Body.Close()

            // Primero intentar decodificar como objeto con clave "activos"
            var responseData map[string]interface{}
            if err := json.NewDecoder(resp.Body).Decode(&responseData); err != nil {
                fmt.Println("Error decodificando respuesta de estados: ", err)
                return
            }

            // Extraer el array de estados (puede estar en diferentes claves)
            var estadosArray []interface{}
            
            // Intentar varias posibles claves
            if activos, exists := responseData["activos"]; exists {
                if arr, ok := activos.([]interface{}); ok {
                    estadosArray = arr
                }
            }

            mu.Lock()
            for _, estadoItem := range estadosArray {
                if estadoData, ok := estadoItem.(map[string]interface{}); ok {
                    if activoID, exists := estadoData["activo_id"]; exists {
                        // Convertir activo_id a int
                        var id int
                        switch v := activoID.(type) {
                        case float64:
                            id = int(v)
                        case int:
                            id = v
                        default:
                            fmt.Printf("Tipo inesperado para activo_id: %T\n", activoID)
                            continue
                        }
                        
                        // Guardar estado
                        if estado, existsEstado := estadoData["estado"]; existsEstado {
                            if estadoStr, ok := estado.(string); ok {
                                estadosMap[id] = estadoStr
                            }
                        }
                    }
                }
            }
            mu.Unlock()
        }()

        wg.Wait() // Esperar a que termine la goroutine de estados
    }

    // Agregar el estado (y sensores si se solicitaron) a cada activo en la lista
    activosConEstado := []interface{}{}
    for _, activo := range activos {
        if activoMap, ok := activo.(map[string]interface{}); ok {
            // Obtener el ID del activo
            if idInterface, exists := activoMap["id"]; exists {
                var activoID int
                
                // Convertir el ID a int
                switch v := idInterface.(type) {
                case float64:
                    activoID = int(v)
                case int:
                    activoID = v
                default:
                    fmt.Printf("Tipo inesperado para id de activo: %T\n", idInterface)
                    activosConEstado = append(activosConEstado, activoMap)
                    continue
                }
                
                // Buscar el estado en el mapa
                if estado, encontrado := estadosMap[activoID]; encontrado {
                    activoMap["estado"] = estado
                } else {
                    // Si no se encuentra estado, poner "Desconocido"
                    activoMap["estado"] = "Desconocido"
                }

                // Si se solicitaron sensores, agregarlos
                if incluirSensores {
                    if sensores, encontrado := sensoresMap[activoID]; encontrado {
                        activoMap["sensores"] = sensores
                    } else {
                        activoMap["sensores"] = []interface{}{}
                    }
                }
            }
            activosConEstado = append(activosConEstado, activoMap)
        }
    }

    result := map[string]interface{}{
        "activos": activosConEstado,
        "total": len(activosConEstado),
    }

    c.JSON(http.StatusOK, result)
}

func ObtenerContactosPorEdificio(c *gin.Context) {
    idEdificio := c.Param("id_edificio")
    if idEdificio == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "ID del edificio es requerido"})
        return
    }

    var wg sync.WaitGroup
    var mu sync.Mutex
    contactos := []interface{}{}
    
    wg.Add(1)
    // Goroutine para obtener contactos
    go func() {
        defer wg.Done()
        url := fmt.Sprintf("%s/tecnicos/edificio/%s", gestionURL, idEdificio)
        resp, err := httpClient.Get(url)
        if err != nil {
            fmt.Println("Error obteniendo contactos: ", err)
            return
        }
        defer resp.Body.Close()
        
        // La respuesta es un array directo, no un objeto con clave "contactos"
        var contactosArray []interface{}
        if err := json.NewDecoder(resp.Body).Decode(&contactosArray); err != nil {
            fmt.Println("Error decodificando contactos: ", err)
            return
        }
        
        mu.Lock()
        contactos = contactosArray
        mu.Unlock()
    }()

    wg.Wait() // Esperar a que termine la goroutine de contactos

    result := map[string]interface{}{
        "contactos": contactos,
        "total": len(contactos),
    }

    c.JSON(http.StatusOK, result)

}

func ForoEdificioHandler(c *gin.Context) {
    idEdificio := c.Param("id_edificio")
    if idEdificio == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "ID del edificio es requerido"})
        return
    }

    // Obtener parámetro de paginación que puede ser opcional
    pagina := c.Query("pagina")
    
    var wg sync.WaitGroup
    var mu sync.Mutex
    var foroData map[string]interface{}

    wg.Add(1)
    go func() {
        defer wg.Done()
        var url string
        
        // Si no tiene parámetro de paginación, se retornan todos los mensajes del foro (¡No recomendado! ☢️)
        if pagina != "" {
            url = fmt.Sprintf("%s/tipos-falla/edificio/%s?pagina=%s", gestionURL, idEdificio, pagina)
        } else {
            url = fmt.Sprintf("%s/tipos-falla/edificio/%s", gestionURL, idEdificio)
        }

        fmt.Println("DEBUG: URL foro edificio:", url)

        resp, err := httpClient.Get(url)
        if err != nil {
            fmt.Println("Error obteniendo mensajes del foro: ", err)
            return
        }
        defer resp.Body.Close()
        
        var data map[string]interface{}
        if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
            fmt.Println("Error decodificando mensajes del foro: ", err)
            return
        }

        mu.Lock()
        foroData = data
        mu.Unlock()
    }()

    wg.Wait()

    if foroData == nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "No se encontraron mensajes del foro"})
        return
    }

    if mensajes, exists := foroData["data"]; exists {
        foroData["foro"] = mensajes
        delete(foroData, "data") // Eliminar la clave original "data"
    }

    c.JSON(http.StatusOK, foroData)
}

func PublicacionForoHandler(c *gin.Context) {
    idEdificio := c.Param("id_edificio")
    if idEdificio == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "ID del edificio es requerido"})
        return
    }

    // Extraer el email desde el token JWT
    var email string
    authHeader := c.GetHeader("Authorization")
    if authHeader == "" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
        return
    }
    
    // fmt.Printf("DEBUG: Authorization header: %s\n", authHeader)

    //Bearer eydsdsdsd...
    tokenParts := strings.Split(authHeader, " ")

    emailInterface, err := extractClaimFromToken(tokenParts[1], "email")
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
        return
    }
    
    // Type assertion para convertir interface{} a string
    var ok bool
    email, ok = emailInterface.(string)
    if !ok {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid email format in token"})
        return
    }
    
    fmt.Printf("Email extraído del token: %s\n", email)

    // Leer el body de la request
    var publicacionData map[string]interface{}
    if err := c.ShouldBindJSON(&publicacionData); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
        return
    }

    // Validar que vengan los campos requeridos
    if _, exists := publicacionData["tipo"]; !exists {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Campo 'tipo' es requerido"})
        return
    }

    if _, exists := publicacionData["descripcion"]; !exists {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Campo 'descripcion' es requerido"})
        return
    }

    fmt.Printf("Publicación recibida para edificio %s por usuario %s: %+v\n", idEdificio, email, publicacionData)
    
    // Body para el microservicio
    // Convertir id_edificio a int
    idEdificioInt := 0
    if _, err := fmt.Sscanf(idEdificio, "%d", &idEdificioInt); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "ID del edificio debe ser un número válido"})
        return
    }

    bodyData := map[string]interface{}{
        "tipo": publicacionData["tipo"],
        "descripcion": publicacionData["descripcion"],
        "email": email,
        "id_edificio": idEdificioInt,
    }

    // Convertir a JSON
    jsonData, err := json.Marshal(bodyData)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error procesando datos"})
        return
    }

    fmt.Printf("DEBUG: Enviando publicación al foro: %+v\n", bodyData)

    // Hacer POST al microservicio de gestión
    var wg sync.WaitGroup
    var mu sync.Mutex
    var responseData map[string]interface{}
    var statusCode int

    wg.Add(1)
    go func() {
        defer wg.Done()
        
        url := fmt.Sprintf("%s/tipos-falla", gestionURL)
        resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
        if err != nil {
            fmt.Println("Error enviando publicación al foro: ", err)
            return
        }
        defer resp.Body.Close()

        mu.Lock()
        statusCode = resp.StatusCode
        mu.Unlock()

        var data map[string]interface{}
        if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
            fmt.Println("Error decodificando respuesta de publicación: ", err)
            return
        }

        mu.Lock()
        responseData = data
        mu.Unlock()
    }()

    wg.Wait()

    // Verificar si hubo respuesta
    if responseData == nil {
        c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Error al crear publicación en el foro"})
        return
    }

    // Retornar la respuesta del microservicio
    c.JSON(statusCode, responseData)
}

func ComentariosForoHandler(c *gin.Context) {
    idPublicacion := c.Param("id_publicacion")
    if idPublicacion == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "ID de la publicación es requerido"})
        return
    }

    // Extraer el email desde el token JWT
    var email string
    authHeader := c.GetHeader("Authorization")
    if authHeader == "" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
        return
    }

    //Bearer eydsdsdsd...
    tokenParts := strings.Split(authHeader, " ")

    emailInterface, err := extractClaimFromToken(tokenParts[1], "email")
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
        return
    }
    
    // Type assertion para convertir interface{} a string
    var ok bool
    email, ok = emailInterface.(string)
    if !ok {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid email format in token"})
        return
    }
    
    fmt.Printf("Email extraído del token: %s\n", email)

    // Leer el body de la request
    var comentarioData map[string]interface{}
    if err := c.ShouldBindJSON(&comentarioData); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
        return
    }

    // Validar que venga el campo requerido
    if _, exists := comentarioData["comentario"]; !exists {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Campo 'comentario' es requerido"})
        return
    }

    // Body para el microservicio

    //Convertir id_publicacion a int
    idPublicacionInt := 0
    if _, err := fmt.Sscanf(idPublicacion, "%d", &idPublicacionInt); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "ID de la publicación debe ser un número válido"})
        return
    }
    bodyData := map[string]interface{}{
        "id_falla": idPublicacionInt,
        "email": email,
        "comentario": comentarioData["comentario"],
    }

    // Convertir a JSON
    jsonData, err := json.Marshal(bodyData)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error procesando datos"})
        return
    }
    
    fmt.Printf("DEBUG: Enviando comentario al foro: %+v\n", bodyData)

    // Hacer POST al microservicio de gestión
    var wg sync.WaitGroup
    var mu sync.Mutex
    var responseData map[string]interface{}
    var statusCode int
    
    wg.Add(1)
    go func() {
        defer wg.Done()

        url := fmt.Sprintf("%s/comentarios", gestionURL)
        resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
        if err != nil {
            fmt.Println("Error enviando comentario al foro: ", err)
            return
        }
        defer resp.Body.Close()

        mu.Lock()
        statusCode = resp.StatusCode
        mu.Unlock()

        var data map[string]interface{}
        if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
            fmt.Println("Error decodificando respuesta de comentario: ", err)
            return
        }

        mu.Lock()
        responseData = data
        mu.Unlock()
    }()

    wg.Wait()

    // Verificar si hubo respuesta
    if responseData == nil {
        c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Error al crear comentario en el foro"})
        return
    }

    // Retornar la respuesta del microservicio
    c.JSON(statusCode, responseData)
}

func ListaActivosHandler(c *gin.Context) {
    idEdificio := c.Param("id_edificio")
    if idEdificio == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "ID del edificio es requerido"})
        return
    }
    // sensores := c.Param("sensores") // Parámetro opcional para incluir sensores

    var wg sync.WaitGroup
    var mu sync.Mutex
    activos := []interface{}{}
    estadosMap := make(map[int]string)
    sensoresMap := make(map[int]interface{})

    wg.Add(1)
    // Goroutine para obtener activos desde Gestión
    go func() {
        defer wg.Done()
        url := fmt.Sprintf("%s/activos/edificio/%s", gestionURL, idEdificio)
        
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

    // Obtener el estado de cada activo mediante el parser
    wg.Add(1)
    go func() {
        defer wg.Done()
        url := fmt.Sprintf("%s/activo/edificio/%s", parserURL, idEdificio)

        resp, err := httpClient.Get(url)
        if err != nil {
            fmt.Println("Error obteniendo estados de activos: ", err)
            return
        }
        defer resp.Body.Close()

        // Primero intentar decodificar como objeto con clave "activos"
        var responseData map[string]interface{}
        if err := json.NewDecoder(resp.Body).Decode(&responseData); err != nil {
            fmt.Println("Error decodificando respuesta de estados: ", err)
            return
        }

        fmt.Printf("DEBUG: Respuesta completa del parser: %+v\n", responseData)

        // Extraer el array de estados (puede estar en diferentes claves)
        var estadosArray []interface{}
        
        // Intentar varias posibles claves
        if activos, exists := responseData["activos"]; exists {
            if arr, ok := activos.([]interface{}); ok {
                estadosArray = arr
            }
        }

        // fmt.Printf("DEBUG: Array de estados extraído: %+v\n", estadosArray)

        mu.Lock()
        for _, estadoItem := range estadosArray {
            if estadoData, ok := estadoItem.(map[string]interface{}); ok {
                if activoID, exists := estadoData["activo_id"]; exists {
                    // Convertir activo_id a int
                    var id int
                    switch v := activoID.(type) {
                    case float64:
                        id = int(v)
                    case int:
                        id = v
                    default:
                        fmt.Printf("Tipo inesperado para activo_id: %T\n", activoID)
                        continue
                    }
                    
                    // Guardar estado
                    if estado, existsEstado := estadoData["estado"]; existsEstado {
                        if estadoStr, ok := estado.(string); ok {
                            estadosMap[id] = estadoStr
                        }
                    }
                    
                    // Guardar sensores si existen
                    if sensores, existsSensores := estadoData["sensores"]; existsSensores {
                        sensoresMap[id] = sensores
                        fmt.Printf("DEBUG: Sensores encontrados para activo %d\n", id)
                    }
                }
            }
        }
        mu.Unlock()
    }()

    wg.Wait() // Esperar a que termine la goroutine de estados

    // fmt.Printf("DEBUG: Estados obtenidos: %+v\n", estadosMap)
    // fmt.Printf("DEBUG: Sensores obtenidos: %+v\n", sensoresMap)

    // Agregar el estado a cada activo en la lista
    activosConEstado := []interface{}{}
    for _, activo := range activos {
        if activoMap, ok := activo.(map[string]interface{}); ok {
            // Obtener el ID del activo
            if idInterface, exists := activoMap["id"]; exists {
                var activoID int
                
                // Convertir el ID a int
                switch v := idInterface.(type) {
                case float64:
                    activoID = int(v)
                case int:
                    activoID = v
                default:
                    fmt.Printf("Tipo inesperado para id de activo: %T\n", idInterface)
                    activosConEstado = append(activosConEstado, activoMap)
                    continue
                }
                
                // Buscar el estado en el mapa
                if estado, encontrado := estadosMap[activoID]; encontrado {
                    activoMap["estado"] = estado
                    fmt.Printf("DEBUG: Activo %d -> Estado: %s\n", activoID, estado)
                } else {
                    // Si no se encuentra estado, poner "Desconocido"
                    activoMap["estado"] = "Desconocido"
                    // fmt.Printf("DEBUG: Activo %d -> Estado no encontrado\n", activoID)
                }

                // Buscar los sensores en el mapa
                if sensores, encontrado := sensoresMap[activoID]; encontrado {
                    activoMap["sensores"] = sensores
                    fmt.Printf("DEBUG: Sensores agregados al activo %d\n", activoID)
                } else {
                    activoMap["sensores"] = []interface{}{} // Array vacío si no hay sensores
                }
            }
            activosConEstado = append(activosConEstado, activoMap)
        }
    }

    result := map[string]interface{}{
        "activos": activosConEstado,
        "total": len(activosConEstado),
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

        if nombre, exists := data["nombre"]; exists {
            mu.Lock()
            if activo != nil {
                activo["nombre"] = nombre
            }
            mu.Unlock()
        }

        if tipo, exists := data["tipo"]; exists {
            mu.Lock()
            if activo != nil {
                activo["tipo"] = tipo
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

func GenerarPDFReporte(c *gin.Context) {
    id := c.Param("id_activo")
    if id == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "ID del activo es requerido"})
        return
    }

    // Extraer el email desde el token JWT
    var email string
    authHeader := c.GetHeader("Authorization")
    if authHeader == "" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
        return
    }

    tokenParts := strings.Split(authHeader, " ")

    emailInterface, err := extractClaimFromToken(tokenParts[1], "email")
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
        return
    }
    
    // Type assertion para convertir interface{} a string
    var ok bool
    email, ok = emailInterface.(string)
    if !ok {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid email format in token"})
        return
    }
    
    fmt.Printf("Email extraído del token: %s\n", email)

    // Leer el body de la request
    var reporteData map[string]interface{}
    if err := c.ShouldBindJSON(&reporteData); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
        return
    }
    
    fmt.Printf("Datos recibidos para generar PDF: %+v\n", reporteData)

    // Construir el body para el microservicio de documentación
    bodyData := map[string]interface{}{
        "campos": reporteData["campos"],
        "usar_firma_predeterminada": true,
        "usuario_id": 1,
    }

    // Convertir a JSON
    jsonData, err := json.Marshal(bodyData)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error procesando datos"})
        return
    }
    
    fmt.Printf("DEBUG: Enviando datos para generar PDF: %+v\n", bodyData)

    // Hacer POST al microservicio de gestión
    url := fmt.Sprintf("%s/reportes/activo/%s", gestionURL, id)
    resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
    if err != nil {
        fmt.Println("Error generando PDF: ", err)
        c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Error al generar el PDF"})
        return
    }
    defer resp.Body.Close()

    fmt.Printf("DEBUG: Status code de respuesta PDF: %d\n", resp.StatusCode)

    // Verificar si la respuesta es exitosa
    if resp.StatusCode != http.StatusOK {
        // Si no es 200, intentar leer como JSON (mensaje de error)
        var errorData map[string]interface{}
        if err := json.NewDecoder(resp.Body).Decode(&errorData); err == nil {
            c.JSON(resp.StatusCode, errorData)
            return
        }
        c.JSON(resp.StatusCode, gin.H{"error": "Error al generar el PDF"})
        return
    }

    // Obtener el Content-Type del microservicio
    contentType := resp.Header.Get("Content-Type")
    if contentType == "" {
        contentType = "application/pdf"
    }

    // Obtener el nombre del archivo si viene en el header
    contentDisposition := resp.Header.Get("Content-Disposition")
    if contentDisposition == "" {
        contentDisposition = fmt.Sprintf("attachment; filename=reporte_activo_%s.pdf", id)
    }

    // Configurar headers para enviar el PDF al cliente
    c.Header("Content-Type", contentType)
    c.Header("Content-Disposition", contentDisposition)
    c.Header("Content-Transfer-Encoding", "binary")

    // Copiar el contenido del PDF directamente a la respuesta
    c.DataFromReader(resp.StatusCode, resp.ContentLength, contentType, resp.Body, nil)
}

func ObtenerAcciones(c *gin.Context) {
    id := c.Param("id_tecnico")
    if id == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "ID del tecnico es requerido"})
        return
    }

    // Hacer GET al microservicio de gestión
    // url := fmt.Sprintf("%s/acciones/tecnico/%s", gestionURL, email)
    url := fmt.Sprintf("%s/acciones/tecnico/%s", gestionURL, id)
    resp, err := httpClient.Get(url)
    if err != nil {
        fmt.Println("Error obteniendo acciones: ", err)
        c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Error al obtener las acciones"})
        return
    }
    defer resp.Body.Close()

    // Verificar si la respuesta es exitosa
    if resp.StatusCode != http.StatusOK {
        var errorData map[string]interface{}
        if err := json.NewDecoder(resp.Body).Decode(&errorData); err == nil {
            c.JSON(resp.StatusCode, errorData)
            return
        }
        c.JSON(resp.StatusCode, gin.H{"error": "Error al obtener las acciones"})
        return
    }

    // Leer la respuesta como array directo
    var accionesArray []interface{}
    if err := json.NewDecoder(resp.Body).Decode(&accionesArray); err != nil {
        fmt.Println("Error decodificando acciones: ", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error procesando las acciones"})
        return
    }

    // Retornar el array directamente
    c.JSON(http.StatusOK, accionesArray)
}

func SubirDocumentoHandler(c *gin.Context) {
    // // Extraer el email desde el token JWT
    // var email string
    // authHeader := c.GetHeader("Authorization")
    // if authHeader != "" {
    //     tokenParts := strings.Split(authHeader, " ")
    //     if len(tokenParts) == 2 {
    //         emailInterface, err := extractClaimFromToken(tokenParts[1], "email")
    //         if err == nil {
    //             if e, ok := emailInterface.(string); ok {
    //                 email = e
    //             }
    //         }
    //     }
    // }

    // fmt.Printf("Usuario subiendo documento: %s\n", email)

    // Parsear multipart/form-data
    err := c.Request.ParseMultipartForm(32 << 20) // 32 MB
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Error parseando form-data"})
        return
    }

    // Obtener el archivo
    file, fileHeader, err := c.Request.FormFile("archivo")
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Campo 'archivo' es requerido"})
        return
    }
    defer file.Close()

    fmt.Printf("Archivo recibido: %s, tamaño: %d bytes\n", fileHeader.Filename, fileHeader.Size)

    // Obtener los demás campos del form
    activoID := c.Request.FormValue("activo_id")
    categoria := c.Request.FormValue("categoria")
    nombre := c.Request.FormValue("nombre")
    descripcion := c.Request.FormValue("descripcion")
    palabrasClave := c.Request.FormValue("palabras_clave")

    // Validar campos requeridos
    if activoID == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Campo 'activo_id' es requerido"})
        return
    }

    if categoria == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Campo 'categoria' es requerido"})
        return
    }

    if nombre == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Campo 'nombre' es requerido"})
        return
    }

    fmt.Printf("Datos del documento - activo_id: %s, categoria: %s, nombre: %s\n", activoID, categoria, nombre)

    // Crear el form-data para el microservicio de documentación
    var requestBody bytes.Buffer
    writer := multipart.NewWriter(&requestBody)

    // Agregar el archivo
    part, err := writer.CreateFormFile("archivo", fileHeader.Filename)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creando form file"})
        return
    }

    // Copiar el contenido del archivo
    _, err = io.Copy(part, file)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error copiando archivo"})
        return
    }

    // Agregar los demás campos
    writer.WriteField("activo_id", activoID)
    writer.WriteField("categoria", categoria)
    writer.WriteField("nombre", nombre)
    
    if descripcion != "" {
        writer.WriteField("descripcion", descripcion)
    }
    
    if palabrasClave != "" {
        writer.WriteField("palabras_clave", palabrasClave)
    }

    // Cerrar el writer
    err = writer.Close()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error cerrando writer"})
        return
    }

    // Hacer POST al microservicio de documentación
    url := fmt.Sprintf("%s/api/v1/documentos", documentacionURL)
    req, err := http.NewRequest("POST", url, &requestBody)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creando request"})
        return
    }

    // Establecer el Content-Type con el boundary correcto
    req.Header.Set("Content-Type", writer.FormDataContentType())

    // Ejecutar la petición
    client := &http.Client{Timeout: 30 * time.Second} // Mayor timeout para uploads
    resp, err := client.Do(req)
    if err != nil {
        fmt.Println("Error subiendo documento: ", err)
        c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Error al subir el documento"})
        return
    }
    defer resp.Body.Close()

    fmt.Printf("DEBUG: Status code de respuesta de subida: %d\n", resp.StatusCode)

    // Verificar si la respuesta es exitosa
    if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
        var errorData map[string]interface{}
        if err := json.NewDecoder(resp.Body).Decode(&errorData); err == nil {
            c.JSON(resp.StatusCode, errorData)
            return
        }
        c.JSON(resp.StatusCode, gin.H{"error": "Error al subir el documento"})
        return
    }

    // Leer la respuesta exitosa
    var responseData map[string]interface{}
    if err := json.NewDecoder(resp.Body).Decode(&responseData); err != nil {
        fmt.Println("Error decodificando respuesta de subida: ", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error procesando respuesta"})
        return
    }

    c.JSON(resp.StatusCode, responseData)
}

func ActualizarEstadoContactoHandler(c *gin.Context) {
    id := c.Param("id_contacto")
    if id == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "ID del contacto es requerido"})
        return
    }

    // Leer el body de la request
    var autorizadoData map[string]interface{}
    if err := c.ShouldBindJSON(&autorizadoData); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
        return
    }

    // Validar que venga el campo requerido
    if _, exists := autorizadoData["autorizado"]; !exists {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Campo 'autorizado' es requerido"})
        return
    }

    // Body
    bodyData := map[string]interface{}{
        "autorizado": autorizadoData["autorizado"],
    }

    // Convertir a JSON
    jsonData, err := json.Marshal(bodyData)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error procesando datos"})
        return
    }

    // Hacer PUT al microservicio de gestión para actualizar el estado
    url := fmt.Sprintf("%s/tecnicos/%s/autorizado", gestionURL, id)
    req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creando request"})
        return
    }

    // Establecer headers
    req.Header.Set("Content-Type", "application/json")

    // Ejecutar la petición
    resp, err := httpClient.Do(req)
    if err != nil {
        fmt.Println("Error enviando la actualización de estado: ", err)
        c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Error al actualizar el estado del contacto"})
        return
    }
    defer resp.Body.Close()

    // fmt.Printf("DEBUG: Status code de actualización estado: %d\n", resp.StatusCode)

    // Verificar si la respuesta es exitosa
    if resp.StatusCode != http.StatusOK {
        var errorData map[string]interface{}
        if err := json.NewDecoder(resp.Body).Decode(&errorData); err == nil {
            c.JSON(resp.StatusCode, errorData)
            return
        }
        c.JSON(resp.StatusCode, gin.H{"error": "Error al actualizar el estado del contacto"})
        return
    }

    // Leer la respuesta exitosa
    var responseData map[string]interface{}
    if err := json.NewDecoder(resp.Body).Decode(&responseData); err != nil {
        fmt.Println("Error decodificando respuesta: ", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error procesando respuesta"})
        return
    }

    // Retornar la respuesta del microservicio
    c.JSON(http.StatusOK, responseData)
}

func ObtenerTodosActivos(c *gin.Context) {
    var wg sync.WaitGroup
    var mu sync.Mutex
    activos := []interface{}{}
    estadosMap := make(map[int]string)
    sensoresMap := make(map[int]interface{})

    // Verificar si se solicitan sensores
    incluirSensores := c.Query("sensores") == "true"

    wg.Add(1)
    // Goroutine para obtener todos los activos desde Gestión
    go func() {
        defer wg.Done()
        url := fmt.Sprintf("%s/activos", gestionURL)
        
        resp, err := httpClient.Get(url)
        if err != nil {
            fmt.Println("Error obteniendo todos los activos: ", err)
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

    // Si se solicitan sensores, consultar individualmente cada activo en el parser
    if incluirSensores {
        fmt.Println("DEBUG: Consultando sensores para cada activo")
        
        for _, activo := range activos {
            if activoMap, ok := activo.(map[string]interface{}); ok {
                if idInterface, exists := activoMap["id"]; exists {
                    var activoID int
                    
                    // Convertir el ID a int
                    switch v := idInterface.(type) {
                    case float64:
                        activoID = int(v)
                    case int:
                        activoID = v
                    default:
                        fmt.Printf("Tipo inesperado para id de activo: %T\n", idInterface)
                        continue
                    }

                    // Consultar sensores y datos para este activo específico
                    wg.Add(1)
                    go func(aID int) {
                        defer wg.Done()
                        
                        // Consultar información de sensores (estado, tipo, unidad, etc.)
                        url := fmt.Sprintf("%s/activo/%d", parserURL, aID)
                        resp, err := httpClient.Get(url)
                        if err != nil {
                            fmt.Printf("Error obteniendo sensores para activo %d: %v\n", aID, err)
                            return
                        }
                        defer resp.Body.Close()

                        var activoData map[string]interface{}
                        if err := json.NewDecoder(resp.Body).Decode(&activoData); err != nil {
                            fmt.Printf("Error decodificando datos del activo %d: %v\n", aID, err)
                            return
                        }

                        mu.Lock()
                        // Guardar estado si existe
                        if estado, exists := activoData["estado"]; exists {
                            if estadoStr, ok := estado.(string); ok {
                                estadosMap[aID] = estadoStr
                            }
                        }

                        // Obtener sensores
                        var sensoresArray []interface{}
                        if sensores, exists := activoData["sensores"]; exists {
                            if sensoresArr, ok := sensores.([]interface{}); ok {
                                sensoresArray = sensoresArr
                            }
                        }
                        mu.Unlock()

                        // Consultar datos históricos de sensores
                        urlDatos := fmt.Sprintf("%s/lectura/%d/datos", parserURL, aID)
                        respDatos, err := httpClient.Get(urlDatos)
                        if err != nil {
                            fmt.Printf("Error obteniendo datos para activo %d: %v\n", aID, err)
                            // Guardar sensores sin datos
                            mu.Lock()
                            sensoresMap[aID] = sensoresArray
                            mu.Unlock()
                            return
                        }
                        defer respDatos.Body.Close()

                        var datosResponse map[string]interface{}
                        if err := json.NewDecoder(respDatos.Body).Decode(&datosResponse); err != nil {
                            fmt.Printf("Error decodificando datos históricos del activo %d: %v\n", aID, err)
                            // Guardar sensores sin datos
                            mu.Lock()
                            sensoresMap[aID] = sensoresArray
                            mu.Unlock()
                            return
                        }

                        // Crear un mapa de datos por sensor_id
                        datosMap := make(map[string]interface{})
                        if sensoresDatos, exists := datosResponse["sensores"]; exists {
                            if sensoresDatosArray, ok := sensoresDatos.([]interface{}); ok {
                                for _, sensorDato := range sensoresDatosArray {
                                    if sensorDatoMap, ok := sensorDato.(map[string]interface{}); ok {
                                        if sensorID, exists := sensorDatoMap["sensor_id"]; exists {
                                            if datos, existsDatos := sensorDatoMap["datos"]; existsDatos {
                                                datosMap[fmt.Sprintf("%v", sensorID)] = datos
                                            }
                                        }
                                    }
                                }
                            }
                        }

                        // Agregar los datos a cada sensor
                        for _, sensor := range sensoresArray {
                            if sensorMap, ok := sensor.(map[string]interface{}); ok {
                                if sensorID, exists := sensorMap["sensor_id"]; exists {
                                    sensorIDStr := fmt.Sprintf("%v", sensorID)
                                    if datos, encontrado := datosMap[sensorIDStr]; encontrado {
                                        sensorMap["datos"] = datos
                                    } else {
                                        sensorMap["datos"] = nil
                                    }
                                }
                            }
                        }

                        // Guardar sensores con datos
                        mu.Lock()
                        sensoresMap[aID] = sensoresArray
                        mu.Unlock()

                    }(activoID)
                }
            }
        }

        wg.Wait() // Esperar a que terminen todas las consultas de sensores
    } else {
        // Si NO se solicitan sensores, usar el endpoint global para estados
        wg.Add(1)
        go func() {
            defer wg.Done()
            url := fmt.Sprintf("%s/activo", parserURL)

            resp, err := httpClient.Get(url)
            if err != nil {
                fmt.Println("Error obteniendo estados de activos: ", err)
                return
            }
            defer resp.Body.Close()

            // Extraer el array de estados
            var estadosArray []interface{}

            // Intentar decodificar primero como array directo
            if err := json.NewDecoder(resp.Body).Decode(&estadosArray); err != nil {
                // Si falla, intentar como objeto con clave "activos"
                resp.Body.Close()
                resp, err = httpClient.Get(url)
                if err != nil {
                    fmt.Println("Error obteniendo estados de activos (segundo intento): ", err)
                    return
                }
                defer resp.Body.Close()

                var responseData map[string]interface{}
                if err := json.NewDecoder(resp.Body).Decode(&responseData); err != nil {
                    fmt.Println("Error decodificando respuesta de estados: ", err)
                    return
                }

                if activos, exists := responseData["activos"]; exists {
                    if arr, ok := activos.([]interface{}); ok {
                        estadosArray = arr
                    }
                }
            }

            mu.Lock()
            for _, estadoItem := range estadosArray {
                if estadoData, ok := estadoItem.(map[string]interface{}); ok {
                    if activoID, exists := estadoData["activo_id"]; exists {
                        // Convertir activo_id a int
                        var id int
                        switch v := activoID.(type) {
                        case float64:
                            id = int(v)
                        case int:
                            id = v
                        default:
                            fmt.Printf("Tipo inesperado para activo_id: %T\n", activoID)
                            continue
                        }
                        
                        // Guardar estado
                        if estado, existsEstado := estadoData["estado"]; existsEstado {
                            if estadoStr, ok := estado.(string); ok {
                                estadosMap[id] = estadoStr
                            }
                        }
                    }
                }
            }
            mu.Unlock()
        }()

        wg.Wait() // Esperar a que termine la goroutine de estados
    }

    // Agregar el estado (y sensores si se solicitaron) a cada activo en la lista
    activosConEstado := []interface{}{}
    for _, activo := range activos {
        if activoMap, ok := activo.(map[string]interface{}); ok {
            // Obtener el ID del activo
            if idInterface, exists := activoMap["id"]; exists {
                var activoID int
                
                // Convertir el ID a int
                switch v := idInterface.(type) {
                case float64:
                    activoID = int(v)
                case int:
                    activoID = v
                default:
                    fmt.Printf("Tipo inesperado para id de activo: %T\n", idInterface)
                    activosConEstado = append(activosConEstado, activoMap)
                    continue
                }
                
                // Buscar el estado en el mapa
                if estado, encontrado := estadosMap[activoID]; encontrado {
                    activoMap["estado"] = estado
                } else {
                    // Si no se encuentra estado, poner "Desconocido"
                    activoMap["estado"] = "Desconocido"
                }

                // Si se solicitaron sensores, agregarlos
                if incluirSensores {
                    if sensores, encontrado := sensoresMap[activoID]; encontrado {
                        activoMap["sensores"] = sensores
                    } else {
                        activoMap["sensores"] = []interface{}{} // Array vacío si no hay sensores
                    }
                }
            }
            activosConEstado = append(activosConEstado, activoMap)
        }
    }

    result := map[string]interface{}{
        "activos": activosConEstado,
        "total": len(activosConEstado),
    }

    c.JSON(http.StatusOK, result)
}

func AccionMantenimientoHandler(c *gin.Context) {
    id := c.Param("id_activo")
    if id == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "ID del activo es requerido"})
        return
    }

    // Leer el body de la request
    var accionData map[string]interface{}
    if err := c.ShouldBindJSON(&accionData); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
        return
    }

    // Validar que vengan los campos requeridos
    if _, exists := accionData["titulo"]; !exists {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Campo 'titulo' es requerido"})
        return
    }

    if _, exists := accionData["tecnico_id"]; !exists {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Campo 'tecnico_id' es requerido"})
        return
    }

    if _, exists := accionData["tipo"]; !exists {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Campo 'tipo' es requerido"})
        return
    }

    if _, exists := accionData["descripcion"]; !exists {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Campo 'descripcion' es requerido"})
        return
    }

    if _, exists := accionData["prioridad"]; !exists {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Campo 'prioridad' es requerido"})
        return
    }

    fmt.Printf("Acción recibida para activo %s: %+v\n", id, accionData)

    // Convertir id_activo a int
    idActivoInt := 0
    if _, err := fmt.Sscanf(id, "%d", &idActivoInt); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "ID del activo inválido"})
        return
    }

    // Construir el body para el microservicio de gestión
    bodyData := map[string]interface{}{
        "titulo":      accionData["titulo"],
        "descripcion": accionData["descripcion"],
        "tipo":        accionData["tipo"],
        "prioridad":   accionData["prioridad"],
        "activo_id":   idActivoInt,
        "tecnico_id":  accionData["tecnico_id"],
    }

    // Convertir a JSON
    jsonData, err := json.Marshal(bodyData)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error procesando datos"})
        return
    }

    fmt.Printf("DEBUG: Enviando acción de mantenimiento: %+v\n", bodyData)

    // Hacer POST al microservicio de gestión
    url := fmt.Sprintf("%s/acciones", gestionURL)
    resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
    if err != nil {
        fmt.Println("Error creando acción de mantenimiento: ", err)
        c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Error al crear la acción de mantenimiento"})
        return
    }
    defer resp.Body.Close()

    fmt.Printf("DEBUG: Status code de respuesta: %d\n", resp.StatusCode)

    // Verificar si la respuesta es exitosa
    if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
        var errorData map[string]interface{}
        if err := json.NewDecoder(resp.Body).Decode(&errorData); err == nil {
            c.JSON(resp.StatusCode, errorData)
            return
        }
        c.JSON(resp.StatusCode, gin.H{"error": "Error al crear la acción de mantenimiento"})
        return
    }

    // Leer la respuesta exitosa
    var responseData map[string]interface{}
    if err := json.NewDecoder(resp.Body).Decode(&responseData); err != nil {
        fmt.Println("Error decodificando respuesta: ", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error procesando respuesta"})
        return
    }

    // Retornar la respuesta del microservicio
    c.JSON(resp.StatusCode, responseData)
}

func ActualizarEstadoAccionHandler(c *gin.Context) {
    id := c.Param("id_accion")
    if id == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "ID de la acción es requerido"})
        return
    }

    // Leer el body de la request
    var estadoData map[string]interface{}
    if err := c.ShouldBindJSON(&estadoData); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
        return
    }

    // Validar que venga el campo requerido
    if _, exists := estadoData["estado"]; !exists {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Campo 'estado' es requerido"})
        return
    }

    fmt.Printf("Actualización de estado recibida para acción %s: %+v\n", id, estadoData)

    // Body para el microservicio
    bodyData := map[string]interface{}{
        "estado": estadoData["estado"],
    }

    // Convertir a JSON
    jsonData, err := json.Marshal(bodyData)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error procesando datos"})
        return
    }

    fmt.Printf("DEBUG: Enviando actualización de estado: %+v\n", bodyData)

    // Hacer PUT al microservicio de gestión
    url := fmt.Sprintf("%s/acciones/%s/estado", gestionURL, id)
    req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creando request"})
        return
    }

    // Establecer headers
    req.Header.Set("Content-Type", "application/json")

    // Ejecutar la petición
    resp, err := httpClient.Do(req)
    if err != nil {
        fmt.Println("Error actualizando estado de acción: ", err)
        c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Error al actualizar el estado de la acción"})
        return
    }
    defer resp.Body.Close()

    fmt.Printf("DEBUG: Status code de actualización estado: %d\n", resp.StatusCode)

    // Verificar si la respuesta es exitosa
    if resp.StatusCode != http.StatusOK {
        var errorData map[string]interface{}
        if err := json.NewDecoder(resp.Body).Decode(&errorData); err == nil {
            c.JSON(resp.StatusCode, errorData)
            return
        }
        c.JSON(resp.StatusCode, gin.H{"error": "Error al actualizar el estado de la acción"})
        return
    }

    // Leer la respuesta exitosa
    var responseData map[string]interface{}
    if err := json.NewDecoder(resp.Body).Decode(&responseData); err != nil {
        fmt.Println("Error decodificando respuesta: ", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error procesando respuesta"})
        return
    }

    // Retornar la respuesta del microservicio
    c.JSON(http.StatusOK, responseData)
}

func ObtenerFirmasUsuarioHandler(c *gin.Context) {
    // Extraer el email desde el token JWT
    var email string
    authHeader := c.GetHeader("Authorization")
    if authHeader == "" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
        return
    }

    tokenParts := strings.Split(authHeader, " ")

    emailInterface, err := extractClaimFromToken(tokenParts[1], "email")
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
        return
    }
    
    // Type assertion para convertir interface{} a string
    var ok bool
    email, ok = emailInterface.(string)
    if !ok {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid email format in token"})
        return
    }
    
    fmt.Printf("Email extraído del token: %s\n", email)

    // Hacer GET al microservicio de gestión
    url := fmt.Sprintf("%s/firmas/usuario/%s", gestionURL, email)
    resp, err := httpClient.Get(url)
    if err != nil {
        fmt.Println("Error obteniendo firmas del usuario: ", err)
        c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Error al obtener las firmas"})
        return
    }
    defer resp.Body.Close()

    // Verificar si la respuesta es exitosa
    if resp.StatusCode != http.StatusOK {
        var errorData map[string]interface{}
        if err := json.NewDecoder(resp.Body).Decode(&errorData); err == nil {
            c.JSON(resp.StatusCode, errorData)
            return
        }
        c.JSON(resp.StatusCode, gin.H{"error": "Error al obtener las firmas"})
        return
    }

    // Leer la respuesta del microservicio
    var firmasData interface{}
    if err := json.NewDecoder(resp.Body).Decode(&firmasData); err != nil {
        fmt.Println("Error decodificando firmas: ", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error procesando las firmas"})
        return
    }

    // Retornar la respuesta del microservicio directamente
    c.JSON(http.StatusOK, firmasData)
}