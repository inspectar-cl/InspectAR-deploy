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
    mlURL               string
    notificationURL     string
	httpClient          *http.Client
)

func init() {
	gestionURL = os.Getenv("GESTION_URL")
    parserURL = os.Getenv("PARSER_URL")
    documentacionURL = os.Getenv("DOCUMENTATION_URL")
    oauth2URL = os.Getenv("OAUTH2_URL")
    mlURL = os.Getenv("ML_URL")
    notificationURL = os.Getenv("NOTIFICATION_URL")
    httpClient = &http.Client{
        Timeout: 15 * time.Second,
    }
    // maxAgeCookie := 3600 // 1 hora, 60*60 segundos
}

func obtenerEdificiosUsuario(email string, userType string) []interface{} {
    if email == "" {
        log.Println("Email vacío, no se pueden obtener edificios")
        return nil
    }

    // Construir URL para obtener edificios
    var url string
    if userType != "user-type:Root" {
        url = fmt.Sprintf("%s/usuarios/edificios/%s", gestionURL, email)
    } else {
        url = fmt.Sprintf("%s/edificios", gestionURL)
    }
    
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
        // fmt.Println("Login data:", loginData)
        if !ok {
            if user, exists := response["user"].(map[string]interface{}); exists {
                if userEmail, hasEmail := user["email"].(string); hasEmail {
                    email = userEmail
                }
            }
        }

        userTypeInterface, _ := extractClaimFromToken(response["access_token"].(string), "scope")
        userTypeStr := ""
        if ut, ok := userTypeInterface.(string); ok {
            userTypeStr = ut
        }
        fmt.Printf("Tipo de usuario extraído del token: %v\n", userTypeInterface)

        edificios := obtenerEdificiosUsuario(email, userTypeStr)

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

        userTypeInterface, _ := extractClaimFromToken(response["access_token"].(string), "scope")
        userTypeStr := ""
        if ut, ok := userTypeInterface.(string); ok {
            userTypeStr = ut
        }
        fmt.Printf("Tipo de usuario extraído del token: %v\n", userTypeInterface)

        edificios := obtenerEdificiosUsuario(email, userTypeStr)

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
    authHeader := c.GetHeader("Authorization")
    email, _ := extractClaimFromToken(authHeader, "email")
    if email == "" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Email no encontrado en el token"})
        return
    }
    fmt.Printf("Email usuario: %s\n", email)

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
    authHeader := c.GetHeader("Authorization")
    email, _ := extractClaimFromToken(authHeader, "email")
    if email == "" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Email no encontrado en el token"})
        return
    }
    fmt.Printf("Email usuario: %s\n", email)

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
    authHeader := c.GetHeader("Authorization")
    email, _ := extractClaimFromToken(authHeader, "email")
    if email == "" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Email no encontrado en el token"})
        return
    }
    fmt.Printf("Email usuario: %s\n", email)

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
        "email": email,
        "firma_id": reporteData["firma_id"],
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
    authHeader := c.GetHeader("Authorization")
    email, _ := extractClaimFromToken(authHeader, "email")
    if email == "" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Email no encontrado en el token"})
        return
    }
    fmt.Printf("Email usuario: %s\n", email)

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

func SubirFirmaUsuarioHandler(c *gin.Context) {
    // Extraer el email desde el token JWT
    authHeader := c.GetHeader("Authorization")
    email, _ := extractClaimFromToken(authHeader, "email")
    if email == "" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Email no encontrado en el token"})
        return
    }
    fmt.Printf("Usuario subiendo firma: %s\n", email)

    // Parsear multipart/form-data
    err := c.Request.ParseMultipartForm(10 << 20) // 10 MB
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

    fmt.Printf("Archivo de firma recibido: %s, tamaño: %d bytes\n", fileHeader.Filename, fileHeader.Size)

    // Obtener los demás campos del form
    nombreArchivo := c.Request.FormValue("nombre_archivo")
    esPredeterminadaStr := c.Request.FormValue("es_predeterminada")

    // Validar campos requeridos
    if nombreArchivo == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Campo 'nombre_archivo' es requerido"})
        return
    }

    if esPredeterminadaStr == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Campo 'es_predeterminada' es requerido"})
        return
    }

    fmt.Printf("Datos de la firma - nombre: %s, es_predeterminada: %s\n", nombreArchivo, esPredeterminadaStr)

    // Crear el form-data para el microservicio de gestión
    var requestBody bytes.Buffer
    writer := multipart.NewWriter(&requestBody)

    // Determinar el Content-Type basado en la extensión del archivo
    contentType := fileHeader.Header.Get("Content-Type")
    if contentType == "" {
        // Si no viene el Content-Type, intentar determinarlo por extensión
        ext := strings.ToLower(fileHeader.Filename[strings.LastIndex(fileHeader.Filename, ".")+1:])
        switch ext {
        case "png":
            contentType = "image/png"
        case "jpg", "jpeg":
            contentType = "image/jpeg"
        case "svg":
            contentType = "image/svg+xml"
        default:
            contentType = "application/octet-stream"
        }
    }

    fmt.Printf("Content-Type detectado: %s\n", contentType)

    // Crear el campo del archivo con el Content-Type correcto
    h := make(map[string][]string)
    h["Content-Disposition"] = []string{fmt.Sprintf(`form-data; name="archivo"; filename="%s"`, fileHeader.Filename)}
    h["Content-Type"] = []string{contentType}
    
    part, err := writer.CreatePart(h)
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
    writer.WriteField("email", fmt.Sprintf("%v", email))
    writer.WriteField("nombre_archivo", nombreArchivo)
    writer.WriteField("es_predeterminada", esPredeterminadaStr)

    // Cerrar el writer
    err = writer.Close()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error cerrando writer"})
        return
    }

    // Hacer POST al microservicio de gestión
    url := fmt.Sprintf("%s/firmas/upload", gestionURL)
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
        fmt.Println("Error subiendo firma: ", err)
        c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Error al subir la firma"})
        return
    }
    defer resp.Body.Close()

    fmt.Printf("DEBUG: Status code de respuesta de subida de firma: %d\n", resp.StatusCode)

    // Verificar si la respuesta es exitosa
    if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
        var errorData map[string]interface{}
        if err := json.NewDecoder(resp.Body).Decode(&errorData); err == nil {
            c.JSON(resp.StatusCode, errorData)
            return
        }
        c.JSON(resp.StatusCode, gin.H{"error": "Error al subir la firma"})
        return
    }

    // Leer la respuesta exitosa
    var responseData map[string]interface{}
    if err := json.NewDecoder(resp.Body).Decode(&responseData); err != nil {
        fmt.Println("Error decodificando respuesta de subida de firma: ", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error procesando respuesta"})
        return
    }

    c.JSON(resp.StatusCode, responseData)
}

func EliminarFirmaUsuarioHandler(c *gin.Context) {
    // Extraer el email desde el token JWT
    authHeader := c.GetHeader("Authorization")
    email, _ := extractClaimFromToken(authHeader, "email")
    if email == "" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Email no encontrado en el token"})
        return
    }
    fmt.Printf("Usuario eliminando firma: %s\n", email)

    // Obtener el ID de la firma desde los parámetros de la URL
    idFirma := c.Param("id_firma")
    if idFirma == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "ID de la firma es requerido"})
        return
    }

    fmt.Printf("Eliminando firma ID: %s para usuario: %s\n", idFirma, email)

    // Construir el body con el email
    bodyData := map[string]interface{}{
        "email": email,
    }

    // Convertir a JSON
    jsonData, err := json.Marshal(bodyData)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error procesando datos"})
        return
    }

    // Hacer DELETE al microservicio de gestión
    url := fmt.Sprintf("%s/firmas/%s", gestionURL, idFirma)
    req, err := http.NewRequest("DELETE", url, bytes.NewBuffer(jsonData))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creando request"})
        return
    }

    // Establecer headers
    req.Header.Set("Content-Type", "application/json")

    // Ejecutar la petición
    resp, err := httpClient.Do(req)
    if err != nil {
        fmt.Println("Error eliminando firma: ", err)
        c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Error al eliminar la firma"})
        return
    }
    defer resp.Body.Close()

    fmt.Printf("DEBUG: Status code de eliminación de firma: %d\n", resp.StatusCode)

    // Verificar si la respuesta es exitosa
    if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
        var errorData map[string]interface{}
        if err := json.NewDecoder(resp.Body).Decode(&errorData); err == nil {
            c.JSON(resp.StatusCode, errorData)
            return
        }
        c.JSON(resp.StatusCode, gin.H{"error": "Error al eliminar la firma"})
        return
    }

    // Leer la respuesta exitosa (si existe)
    var responseData map[string]interface{}
    if resp.StatusCode == http.StatusOK {
        if err := json.NewDecoder(resp.Body).Decode(&responseData); err != nil {
            fmt.Println("Error decodificando respuesta de eliminación: ", err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Error procesando respuesta"})
            return
        }
        c.JSON(http.StatusOK, responseData)
    } else {
        // Si es 204 No Content, responder con mensaje de éxito
        c.JSON(http.StatusOK, gin.H{"message": "Firma eliminada exitosamente"})
    }
}

func ActualizarFirmaUsuarioHandler(c *gin.Context) {
    // Extraer el email desde el token JWT
    authHeader := c.GetHeader("Authorization")
    email, _ := extractClaimFromToken(authHeader, "email")
    if email == "" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Email no encontrado en el token"})
        return
    }
    fmt.Printf("Usuario actualizando firma: %s\n", email)

    // Obtener el ID de la firma desde los parámetros de la URL
    idFirma := c.Param("id_firma")
    if idFirma == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "ID de la firma es requerido"})
        return
    }

    // Leer el body de la request
    var firmaData map[string]interface{}
    if err := c.ShouldBindJSON(&firmaData); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
        return
    }

    // Validar que vengan los campos requeridos
    if _, exists := firmaData["nombre_archivo"]; !exists {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Campo 'nombre_archivo' es requerido"})
        return
    }

    if _, exists := firmaData["es_predeterminada"]; !exists {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Campo 'es_predeterminada' es requerido"})
        return
    }

    fmt.Printf("Actualizando firma ID: %s - datos: %+v\n", idFirma, firmaData)

    // Construir el body para el microservicio
    bodyData := map[string]interface{}{
        "nombre_archivo": firmaData["nombre_archivo"],
        "es_predeterminada": firmaData["es_predeterminada"],
    }

    // Convertir a JSON
    jsonData, err := json.Marshal(bodyData)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error procesando datos"})
        return
    }

    fmt.Printf("DEBUG: Enviando actualización de firma: %+v\n", bodyData)

    // Hacer PUT al microservicio de gestión
    url := fmt.Sprintf("%s/firmas/%s", gestionURL, idFirma)
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
        fmt.Println("Error actualizando firma: ", err)
        c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Error al actualizar la firma"})
        return
    }
    defer resp.Body.Close()

    fmt.Printf("DEBUG: Status code de actualización de firma: %d\n", resp.StatusCode)

    // Verificar si la respuesta es exitosa
    if resp.StatusCode != http.StatusOK {
        var errorData map[string]interface{}
        if err := json.NewDecoder(resp.Body).Decode(&errorData); err == nil {
            c.JSON(resp.StatusCode, errorData)
            return
        }
        c.JSON(resp.StatusCode, gin.H{"error": "Error al actualizar la firma"})
        return
    }

    // Leer la respuesta exitosa
    var responseData map[string]interface{}
    if err := json.NewDecoder(resp.Body).Decode(&responseData); err != nil {
        fmt.Println("Error decodificando respuesta de actualización: ", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error procesando respuesta"})
        return
    }

    // Retornar la respuesta del microservicio
    c.JSON(http.StatusOK, responseData)
}

func ObtenerDocumentosHandler(c *gin.Context) {
    // Leer el body con la lista de activos
    var requestData map[string]interface{}
    if err := c.ShouldBindJSON(&requestData); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
        return
    }

    // Validar que venga el array de activos
    activosInterface, exists := requestData["activos"]
    if !exists {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Campo 'activos' es requerido"})
        return
    }

    activos, ok := activosInterface.([]interface{})
    if !ok {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Campo 'activos' debe ser un array"})
        return
    }

    fmt.Printf("Obteniendo documentos para %d activos\n", len(activos))

    var wg sync.WaitGroup
    var mu sync.Mutex
    todosLosDocumentos := []interface{}{}

    // Hacer peticiones concurrentes para cada activo
    for _, activo := range activos {
        activoMap, ok := activo.(map[string]interface{})
        if !ok {
            continue
        }

        // Extraer el ID del activo
        var activoID int
        switch id := activoMap["id"].(type) {
        case float64:
            activoID = int(id)
        case int:
            activoID = id
        default:
            continue
        }

        wg.Add(1)
        go func(id int) {
            defer wg.Done()

            url := fmt.Sprintf("%s/api/v1/documentos/activo/%d", documentacionURL, id)
            resp, err := httpClient.Get(url)
            if err != nil {
                fmt.Printf("Error obteniendo documentos del activo %d: %v\n", id, err)
                return
            }
            defer resp.Body.Close()

            if resp.StatusCode != http.StatusOK {
                fmt.Printf("Error en respuesta de documentos del activo %d: status %d\n", id, resp.StatusCode)
                return
            }

            var data map[string]interface{}
            if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
                fmt.Printf("Error decodificando documentos del activo %d: %v\n", id, err)
                return
            }

            // Extraer solo el array de documentos
            if documentos, exists := data["documentos"].([]interface{}); exists {
                mu.Lock()
                todosLosDocumentos = append(todosLosDocumentos, documentos...)
                mu.Unlock()
            }
        }(activoID)
    }

    wg.Wait()

    // Retornar la respuesta consolidada
    result := map[string]interface{}{
        "documentos": todosLosDocumentos,
    }

    c.JSON(http.StatusOK, result)
}

// Rutas Sprint 3
func CrearActivoHandler(c *gin.Context) {
    authHeader := c.GetHeader("Authorization")
    email, _ := extractClaimFromToken(authHeader, "email")
    if email == "" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Email no encontrado en el token"})
        return
    }
    fmt.Printf("Usuario creando activo: %s\n", email)

    // Leer el body de la request
    var activoData map[string]interface{}
    if err := c.ShouldBindJSON(&activoData); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
        return
    }
    
    // Validar campos requeridos
    if _, exists := activoData["nombre"]; !exists {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Campo 'nombre' es requerido"})
        return
    }

    if _, exists := activoData["tipo"]; !exists {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Campo 'tipo' es requerido"})
        return
    }

    if _, exists := activoData["descripcion"]; !exists {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Campo 'descripcion' es requerido"})
        return
    }

    if _, exists := activoData["ubicacion"]; !exists {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Campo 'ubicacion' es requerido"})
        return
    }

    if _, exists := activoData["edificio_id"]; !exists {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Campo 'edificio_id' es requerido"})
        return
    }

    fmt.Printf("Datos del activo a crear: %+v\n", activoData)

    // Construir el body para el microservicio de gestión
    bodyData := map[string]interface{}{
        "nombre":      activoData["nombre"],
        "tipo":        activoData["tipo"],
        "descripcion": activoData["descripcion"],
        "ubicacion":   activoData["ubicacion"],
        "edificio_id": activoData["edificio_id"],
        "email":       email,
    }

    // Convertir a JSON
    jsonData, err := json.Marshal(bodyData)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error procesando datos"})
        return
    }

    fmt.Printf("DEBUG: Enviando datos para crear activo: %+v\n", bodyData)

    // Hacer POST al microservicio de gestión
    url := fmt.Sprintf("%s/admin/activos", gestionURL)
    resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
    if err != nil {
        fmt.Println("Error creando activo: ", err)
        c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Error al crear el activo"})
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
        c.JSON(resp.StatusCode, gin.H{"error": "Error al crear el activo"})
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

func EditarActivoHandler(c *gin.Context) {
    id := c.Param("id_activo")
    if id == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "ID del activo es requerido"})
        return
    }

    authHeader := c.GetHeader("Authorization")
    email, _ := extractClaimFromToken(authHeader, "email")
    if email == "" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Email no encontrado en el token"})
        return
    }
    fmt.Printf("Usuario editando activo %s: %s\n", id, email)

    // Leer el body de la request
    var activoData map[string]interface{}
    if err := c.ShouldBindJSON(&activoData); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
        return
    }

	// Validar campos obligatorios
	nombre, existeNombre := activoData["nombre"]
	if !existeNombre || nombre == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El campo 'nombre' es obligatorio"})
		return
	}

	ubicacion, existeUbicacion := activoData["ubicacion"]
	if !existeUbicacion || ubicacion == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El campo 'ubicacion' es obligatorio"})
		return
	}

	descripcion, existeDescripcion := activoData["descripcion"]
	if !existeDescripcion || descripcion == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El campo 'descripcion' es obligatorio"})
		return
	}

	fmt.Printf("Datos del activo a editar: %+v\n", activoData)

	// Construir el body para el microservicio de gestión
	bodyData := map[string]interface{}{
		"email":       email,
		"nombre":      nombre,
		"ubicacion":   ubicacion,
		"descripcion": descripcion,
	}    // Convertir a JSON
    jsonData, err := json.Marshal(bodyData)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error procesando datos"})
        return
    }

    fmt.Printf("DEBUG: Enviando datos para editar activo: %+v\n", bodyData)

    // Hacer PUT al microservicio de gestión
    url := fmt.Sprintf("%s/admin/activos/%s", gestionURL, id)
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
        fmt.Println("Error editando activo: ", err)
        c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Error al editar el activo"})
        return
    }
    defer resp.Body.Close()

    fmt.Printf("DEBUG: Status code de respuesta: %d\n", resp.StatusCode)

    // Verificar si la respuesta es exitosa
    if resp.StatusCode != http.StatusOK {
        var errorData map[string]interface{}
        if err := json.NewDecoder(resp.Body).Decode(&errorData); err == nil {
            c.JSON(resp.StatusCode, errorData)
            return
        }
        c.JSON(resp.StatusCode, gin.H{"error": "Error al editar el activo"})
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

func EliminarActivoHandler(c *gin.Context) {
    id := c.Param("id_activo")
    if id == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "ID del activo es requerido"})
        return
    }

    authHeader := c.GetHeader("Authorization")
    email, _ := extractClaimFromToken(authHeader, "email")
    if email == "" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Email no encontrado en el token"})
        return
    }
    fmt.Printf("Usuario eliminando activo %s: %s\n", id, email)

    // Construir URL con el email como query parameter
    url := fmt.Sprintf("%s/admin/activos/%s?email=%s", gestionURL, id, email)
    
    // Crear request DELETE
    req, err := http.NewRequest("DELETE", url, nil)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creando request"})
        return
    }

    // Establecer headers
    req.Header.Set("Content-Type", "application/json")

    // Ejecutar la petición
    resp, err := httpClient.Do(req)
    if err != nil {
        fmt.Println("Error eliminando activo: ", err)
        c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Error al eliminar el activo"})
        return
    }
    defer resp.Body.Close()

    fmt.Printf("DEBUG: Status code de respuesta: %d\n", resp.StatusCode)

    // Verificar si la respuesta es exitosa
    if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
        var errorData map[string]interface{}
        if err := json.NewDecoder(resp.Body).Decode(&errorData); err == nil {
            c.JSON(resp.StatusCode, errorData)
            return
        }
        c.JSON(resp.StatusCode, gin.H{"error": "Error al eliminar el activo"})
        return
    }

    // Leer la respuesta exitosa (si existe contenido)
    if resp.StatusCode == http.StatusOK {
        var responseData map[string]interface{}
        if err := json.NewDecoder(resp.Body).Decode(&responseData); err != nil {
            fmt.Println("Error decodificando respuesta: ", err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Error procesando respuesta"})
            return
        }
        c.JSON(http.StatusOK, responseData)
    } else {
        // Si es 204 No Content, responder con mensaje de éxito
        c.JSON(http.StatusOK, gin.H{"message": "Activo eliminado exitosamente"})
    }
}

func CrearSensorHandler(c *gin.Context) {
    authHeader := c.GetHeader("Authorization")
    email, _ := extractClaimFromToken(authHeader, "email")
    if email == "" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Email no encontrado en el token"})
        return
    }
    fmt.Printf("Usuario creando sensor: %s\n", email)

    // Leer el body de la request
    var sensorData map[string]interface{}
    if err := c.ShouldBindJSON(&sensorData); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
        return
    }

    // Validar campos requeridos
    if _, exists := sensorData["activo_id"]; !exists {
        c.JSON(http.StatusBadRequest, gin.H{"error": "El campo 'activo_id' es obligatorio"})
        return
    }

    if _, exists := sensorData["nombre"]; !exists {
        c.JSON(http.StatusBadRequest, gin.H{"error": "El campo 'nombre' es obligatorio"})
        return
    }

    if _, exists := sensorData["tipo"]; !exists {
        c.JSON(http.StatusBadRequest, gin.H{"error": "El campo 'tipo' es obligatorio"})
        return
    }

    if _, exists := sensorData["unidad"]; !exists {
        c.JSON(http.StatusBadRequest, gin.H{"error": "El campo 'unidad' es obligatorio"})
        return
    }

    fmt.Printf("Datos del sensor a crear: %+v\n", sensorData)

    // Construir el body para el microservicio de gestión
    bodyData := map[string]interface{}{
        "activo_id": sensorData["activo_id"],
        "nombre":    sensorData["nombre"],
        "tipo":      sensorData["tipo"],
        "unidad":    sensorData["unidad"],
        "email":     email,
    }

    // Convertir a JSON
    jsonData, err := json.Marshal(bodyData)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error procesando datos"})
        return
    }

    fmt.Printf("DEBUG: Enviando datos para crear sensor: %+v\n", bodyData)

    // Hacer POST al microservicio de gestión
    url := fmt.Sprintf("%s/admin/sensores", gestionURL)
    resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
    if err != nil {
        fmt.Println("Error creando sensor: ", err)
        c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Error al crear el sensor"})
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
        c.JSON(resp.StatusCode, gin.H{"error": "Error al crear el sensor"})
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

func EditarSensorHandler(c *gin.Context) {
    // SIN FINALIZAR ☢️
    id := c.Param("id_sensor")
    if id == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "ID del sensor es requerido"})
        return
    }

    authHeader := c.GetHeader("Authorization")
    email, _ := extractClaimFromToken(authHeader, "email")
    if email == "" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Email no encontrado en el token"})
        return
    }
    fmt.Printf("Usuario editando sensor %s: %s\n", id, email)

    // Leer el body de la request
    var sensorData map[string]interface{}
    if err := c.ShouldBindJSON(&sensorData); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
        return
    }

    

}

func ObtenerAnomaliaActivoHandler(c *gin.Context) {
    idActivo := c.Param("id_activo")
    if idActivo == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "ID del activo es requerido"})
        return
    }
    fmt.Printf("Obteniendo anomalías para activo ID: %s\n", idActivo)

    // Paso 1: POST a /detect/{id_activo} para iniciar detección
    detectURL := fmt.Sprintf("%s/detect/%s", mlURL, idActivo)
    postReq, err := http.NewRequest("POST", detectURL, nil)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creando solicitud de detección"})
        return
    }

    postResp, err := httpClient.Do(postReq)
    if err != nil {
        c.JSON(http.StatusBadGateway, gin.H{"error": "Error llamando al servicio de ML para detección"})
        return
    }
    defer postResp.Body.Close()

    if postResp.StatusCode != http.StatusOK {
        c.JSON(postResp.StatusCode, gin.H{"error": "El servicio de ML no pudo detectar anomalías"})
        return
    }

    // Paso 2: GET a /anomalies/activo/{id_activo} para obtener resultados
    anomaliesURL := fmt.Sprintf("%s/anomalies/activo/%s", mlURL, idActivo)
    getReq, err := http.NewRequest("GET", anomaliesURL, nil)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creando solicitud de anomalías"})
        return
    }

    getResp, err := httpClient.Do(getReq)
    if err != nil {
        c.JSON(http.StatusBadGateway, gin.H{"error": "Error llamando al servicio de ML para obtener anomalías"})
        return
    }
    defer getResp.Body.Close()

    // Leer la respuesta del GET
    body, err := io.ReadAll(getResp.Body)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error leyendo respuesta del servicio de ML"})
        return
    }

    // Retornar la respuesta del GET al cliente
    c.Data(getResp.StatusCode, getResp.Header.Get("Content-Type"), body)
}

func CrearEdificioHandler(c *gin.Context) {
    authHeader := c.GetHeader("Authorization")
    email, _ := extractClaimFromToken(authHeader, "email")
    if email == "" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Email no encontrado en el token"})
        return
    }
    fmt.Printf("Usuario creando edificio: %s\n", email)

    // Leer el body de la request
    var edificioData map[string]interface{}
    if err := c.ShouldBindJSON(&edificioData); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
        return
    }

    // Validar campos requeridos
    if _, exists := edificioData["nombre"]; !exists {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Campo 'nombre' es requerido"})
        return
    }

    if _, exists := edificioData["direccion"]; !exists {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Campo 'direccion' es requerido"})
        return
    }

    if _, exists := edificioData["latitud"]; !exists {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Campo 'latitud' es requerido"})
        return
    }

    if _, exists := edificioData["longitud"]; !exists {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Campo 'longitud' es requerido"})
        return
    }

    fmt.Printf("Datos del edificio a crear: %+v\n", edificioData)

    // Construir el body para el microservicio de gestión
    bodyData := map[string]interface{}{
        "nombre":    edificioData["nombre"],
        "direccion": edificioData["direccion"],
        "latitud":   edificioData["latitud"],
        "longitud":  edificioData["longitud"],
        "email":     email,
    }

    // Convertir a JSON
    jsonData, err := json.Marshal(bodyData)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error procesando datos"})
        return
    }

    fmt.Printf("DEBUG: Enviando datos para crear edificio: %+v\n", bodyData)

    // Hacer POST al microservicio de gestión
    url := fmt.Sprintf("%s/admin/edificios", gestionURL)
    resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
    if err != nil {
        fmt.Println("Error creando edificio: ", err)
        c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Error al crear el edificio"})
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
        c.JSON(resp.StatusCode, gin.H{"error": "Error al crear el edificio"})
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

func EliminarEdificioHandler(c *gin.Context) {
    id := c.Param("id_edificio")
    if id == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "ID del edificio es requerido"})
        return
    }

    authHeader := c.GetHeader("Authorization")
    email, _ := extractClaimFromToken(authHeader, "email")
    if email == "" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Email no encontrado en el token"})
        return
    }
    fmt.Printf("Usuario eliminando edificio %s: %s\n", id, email)

    // Construir URL con el email como query parameter
    url := fmt.Sprintf("%s/admin/edificios/%s?email=%s", gestionURL, id, email)
    
    // Crear request DELETE
    req, err := http.NewRequest("DELETE", url, nil)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creando request"})
        return
    }

    // Establecer headers
    req.Header.Set("Content-Type", "application/json")

    // Ejecutar la petición
    resp, err := httpClient.Do(req)
    if err != nil {
        fmt.Println("Error eliminando edificio: ", err)
        c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Error al eliminar el edificio"})
        return
    }
    defer resp.Body.Close()

    fmt.Printf("DEBUG: Status code de respuesta: %d\n", resp.StatusCode)

    // Verificar si la respuesta es exitosa
    if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
        var errorData map[string]interface{}
        if err := json.NewDecoder(resp.Body).Decode(&errorData); err == nil {
            c.JSON(resp.StatusCode, errorData)
            return
        }
        c.JSON(resp.StatusCode, gin.H{"error": "Error al eliminar el edificio"})
        return
    }

    // Leer la respuesta exitosa (si existe contenido)
    if resp.StatusCode == http.StatusOK {
        var responseData map[string]interface{}
        if err := json.NewDecoder(resp.Body).Decode(&responseData); err != nil {
            fmt.Println("Error decodificando respuesta: ", err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Error procesando respuesta"})
            return
        }
        c.JSON(http.StatusOK, responseData)
    } else {
        // Si es 204 No Content, responder con mensaje de éxito
        c.JSON(http.StatusOK, gin.H{"message": "Edificio eliminado exitosamente"})
    }
}

func ActualizarEdificioHandler(c *gin.Context) {
    id := c.Param("id_edificio")
    if id == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "ID del edificio es requerido"})
        return
    }

    authHeader := c.GetHeader("Authorization")
    email, _ := extractClaimFromToken(authHeader, "email")
    if email == "" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Email no encontrado en el token"})
        return
    }
    fmt.Printf("Usuario actualizando edificio %s: %s\n", id, email)

    // Leer el body de la request
    var edificioData map[string]interface{}
    if err := c.ShouldBindJSON(&edificioData); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
        return
    }

    // Validar que al menos uno de los campos actualizables esté presente
    hasNombre := false
    hasDireccion := false
    hasLatitud := false
    hasLongitud := false

    if _, exists := edificioData["nombre"]; exists {
        hasNombre = true
    }
    if _, exists := edificioData["direccion"]; exists {
        hasDireccion = true
    }
    if _, exists := edificioData["latitud"]; exists {
        hasLatitud = true
    }
    if _, exists := edificioData["longitud"]; exists {
        hasLongitud = true
    }

    if !hasNombre && !hasDireccion && !hasLatitud && !hasLongitud {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Debe proporcionar al menos un campo para actualizar (nombre, direccion, latitud, longitud)"})
        return
    }

    fmt.Printf("Datos del edificio a actualizar: %+v\n", edificioData)

    // Construir el body para el microservicio de gestión
    bodyData := map[string]interface{}{
        "email": email,
    }

    // Agregar solo los campos que vienen en la request
    if hasNombre {
        bodyData["nombre"] = edificioData["nombre"]
    }
    if hasDireccion {
        bodyData["direccion"] = edificioData["direccion"]
    }
    if hasLatitud {
        bodyData["latitud"] = edificioData["latitud"]
    }
    if hasLongitud {
        bodyData["longitud"] = edificioData["longitud"]
    }

    // Convertir a JSON
    jsonData, err := json.Marshal(bodyData)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error procesando datos"})
        return
    }

    fmt.Printf("DEBUG: Enviando datos para actualizar edificio: %+v\n", bodyData)

    // Hacer PUT al microservicio de gestión
    url := fmt.Sprintf("%s/admin/edificios/%s", gestionURL, id)
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
        fmt.Println("Error actualizando edificio: ", err)
        c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Error al actualizar el edificio"})
        return
    }
    defer resp.Body.Close()

    fmt.Printf("DEBUG: Status code de respuesta: %d\n", resp.StatusCode)

    // Verificar si la respuesta es exitosa
    if resp.StatusCode != http.StatusOK {
        var errorData map[string]interface{}
        if err := json.NewDecoder(resp.Body).Decode(&errorData); err == nil {
            c.JSON(resp.StatusCode, errorData)
            return
        }
        c.JSON(resp.StatusCode, gin.H{"error": "Error al actualizar el edificio"})
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

func ObtenerInfoQRHandler(c *gin.Context) {
	codigoActivo := c.Param("codigo_activo")
	if codigoActivo == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Código del activo es requerido"})
		return
	}

	// Construir URL del microservicio de gestión
	url := fmt.Sprintf("%s/qr/%s", gestionURL, codigoActivo)

	// Crear request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Error creando request a gestión: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno del servidor"})
		return
	}

	// Ejecutar request
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Printf("Error llamando a gestión: %v", err)
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Servicio no disponible"})
		return
	}
	defer resp.Body.Close()

	// Leer respuesta
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error leyendo respuesta: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error procesando respuesta"})
		return
	}

	// Si la respuesta no es exitosa, retornarla tal cual
	if resp.StatusCode != http.StatusOK {
		c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), body)
		return
	}

	// Parsear la respuesta para obtener el ID del activo
	var activoInfo map[string]interface{}
	if err := json.Unmarshal(body, &activoInfo); err != nil {
		log.Printf("Error parseando respuesta: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error procesando información del activo"})
		return
	}

	// Extraer el ID del activo
	activoID, ok := activoInfo["id"]
	if !ok {
		log.Printf("ID del activo no encontrado en la respuesta")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ID del activo no encontrado"})
		return
	}

	// Convertir el ID a string para la URL
	var activoIDStr string
	switch v := activoID.(type) {
	case float64:
		activoIDStr = fmt.Sprintf("%.0f", v)
	case int:
		activoIDStr = fmt.Sprintf("%d", v)
	default:
		activoIDStr = fmt.Sprintf("%v", v)
	}

	// Obtener el estado del activo desde el Parser
	parserURL := fmt.Sprintf("%s/activo/%s", parserURL, activoIDStr)
	parserReq, err := http.NewRequest("GET", parserURL, nil)
	if err != nil {
		log.Printf("Error creando request al parser: %v", err)
		// Si falla, retornar la info sin estado
		c.JSON(http.StatusOK, activoInfo)
		return
	}

	parserResp, err := httpClient.Do(parserReq)
	if err != nil {
		log.Printf("Error llamando al parser: %v", err)
		// Si falla, retornar la info sin estado
		c.JSON(http.StatusOK, activoInfo)
		return
	}
	defer parserResp.Body.Close()

	// Si el parser responde exitosamente, extraer el estado
	if parserResp.StatusCode == http.StatusOK {
		var parserData map[string]interface{}
		if err := json.NewDecoder(parserResp.Body).Decode(&parserData); err == nil {
			if estado, exists := parserData["estado"]; exists {
				activoInfo["estado"] = estado
			}
		}
	}

	// Retornar la información completa con el estado
	c.JSON(http.StatusOK, activoInfo)
}

func CrearTicketHandler(c *gin.Context) {
	// Extraer el email desde el token JWT
	authHeader := c.GetHeader("Authorization")
	email, _ := extractClaimFromToken(authHeader, "email")
	if email == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email no encontrado en el token"})
		return
	}
	fmt.Printf("Usuario creando ticket: %s\n", email)

	// Leer el body de la request
	var ticketData map[string]interface{}
	if err := c.ShouldBindJSON(&ticketData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	// Validar campos obligatorios
	tipoEntidad, existeTipo := ticketData["tipo_entidad"]
	if !existeTipo || tipoEntidad == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Campo 'tipo_entidad' es requerido"})
		return
	}

	tipoOperacion, existeOperacion := ticketData["tipo_operacion"]
	if !existeOperacion || tipoOperacion == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Campo 'tipo_operacion' es requerido"})
		return
	}

	justificacion, existeJustificacion := ticketData["justificacion"]
	if !existeJustificacion || justificacion == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Campo 'justificacion' es requerido"})
		return
	}

	// Validar valores permitidos
	tipoEntidadStr := fmt.Sprintf("%v", tipoEntidad)
	tipoOperacionStr := fmt.Sprintf("%v", tipoOperacion)

	if tipoEntidadStr != "edificio" && tipoEntidadStr != "activo" && tipoEntidadStr != "tecnico" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tipo_entidad debe ser 'edificio', 'activo' o 'tecnico'"})
		return
	}

	if tipoOperacionStr != "ingreso" && tipoOperacionStr != "modificacion" && tipoOperacionStr != "eliminacion" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tipo_operacion debe ser 'ingreso', 'modificacion' o 'eliminacion'"})
		return
	}

	fmt.Printf("Creando ticket de %s para %s\n", tipoOperacionStr, tipoEntidadStr)

	// Construir el body para el microservicio de notificaciones
	// Siempre incluir los campos base
	bodyData := map[string]interface{}{
		"tipo_entidad":    tipoEntidad,
		"tipo_operacion":  tipoOperacion,
		"usuario_email":   email,
		"justificacion":   justificacion,
	}

	// Agregar campos específicos según la entidad y operación
	// Para modificación y eliminación, el ID es obligatorio
	if tipoOperacionStr == "modificacion" || tipoOperacionStr == "eliminacion" {
		var idKey string
		switch tipoEntidadStr {
		case "edificio":
			idKey = "edificio_id"
		case "activo":
			idKey = "activo_id"
		case "tecnico":
			idKey = "tecnico_id"
		}

		if _, exists := ticketData[idKey]; !exists {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Campo '%s' es requerido para %s", idKey, tipoOperacionStr)})
			return
		}
		bodyData[idKey] = ticketData[idKey]
	}

	// Agregar todos los demás campos que vengan en el request
	// (campos específicos de edificio, activo o técnico)
	for key, value := range ticketData {
		// Saltar los campos que ya agregamos
		if key == "tipo_entidad" || key == "tipo_operacion" || key == "justificacion" {
			continue
		}
		// Agregar el resto de campos (edificio_nombre, activo_tipo, etc.)
		bodyData[key] = value
	}

	// Convertir a JSON
	jsonData, err := json.Marshal(bodyData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error procesando datos"})
		return
	}

	fmt.Printf("DEBUG: Enviando ticket: %+v\n", bodyData)

	// Hacer POST al microservicio de notificaciones
	url := fmt.Sprintf("%s/tickets", notificationURL)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creando ticket: ", err)
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Error al crear el ticket"})
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
		c.JSON(resp.StatusCode, gin.H{"error": "Error al crear el ticket"})
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

func ResolverTicketHandler(c *gin.Context) {
	// Obtener el ID del ticket desde los parámetros de la URL
	idTicket := c.Param("id_ticket")
	if idTicket == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID del ticket es requerido"})
		return
	}

	// Extraer el email desde el token JWT
	authHeader := c.GetHeader("Authorization")
	email, _ := extractClaimFromToken(authHeader, "email")
	if email == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email no encontrado en el token"})
		return
	}
	fmt.Printf("Usuario resolviendo ticket %s: %s\n", idTicket, email)

	// Leer el body de la request
	var resolverData map[string]interface{}
	if err := c.ShouldBindJSON(&resolverData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	// Validar campo obligatorio
	comentarioAdmin, existeComentario := resolverData["comentario_admin"]
	if !existeComentario || comentarioAdmin == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Campo 'comentario_admin' es requerido"})
		return
	}

	fmt.Printf("Resolviendo ticket con comentario: %v\n", comentarioAdmin)

	// Construir el body para el microservicio de notificaciones
	bodyData := map[string]interface{}{
		"resuelto_por_email": email,
		"comentario_admin":   comentarioAdmin,
	}

	// Convertir a JSON
	jsonData, err := json.Marshal(bodyData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error procesando datos"})
		return
	}

	fmt.Printf("DEBUG: Enviando resolución de ticket: %+v\n", bodyData)

	// Hacer PUT al microservicio de notificaciones
	url := fmt.Sprintf("%s/tickets/%s/resolver", notificationURL, idTicket)
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
		fmt.Println("Error resolviendo ticket: ", err)
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Error al resolver el ticket"})
		return
	}
	defer resp.Body.Close()

	fmt.Printf("DEBUG: Status code de respuesta: %d\n", resp.StatusCode)

	// Verificar si la respuesta es exitosa
	if resp.StatusCode != http.StatusOK {
		var errorData map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&errorData); err == nil {
			c.JSON(resp.StatusCode, errorData)
			return
		}
		c.JSON(resp.StatusCode, gin.H{"error": "Error al resolver el ticket"})
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
