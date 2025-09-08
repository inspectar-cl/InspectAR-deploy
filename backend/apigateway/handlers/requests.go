package handlers

import (
    "encoding/json"
    "fmt"
    "net/http"
    "os"
    "sync"
    "time"

	"github.com/gin-gonic/gin"
)

var (
    gestionURL string
    parserURL  string
	httpClient *http.Client
)

func init() {
	gestionURL = os.Getenv("GESTION_URL")
    parserURL = os.Getenv("PARSER_URL")
    httpClient = &http.Client{
        Timeout: 15 * time.Second,
    }
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

	tiposActivos := []string{"caldera"}
    // tiposActivos := []string{"bomba de agua", "caldera", "ascensor", "transformador"}
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
    id := c.Param("id")
    if id == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "ID del activo es requerido"})
        return
    }

    result := map[string]interface{}{
        "reporte": fmt.Sprintf("Reporte generado para activo ID %s", id),
    }
    c.JSON(http.StatusOK, result)
}
