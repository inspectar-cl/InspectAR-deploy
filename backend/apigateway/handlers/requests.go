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

	tiposActivos := []string{"bomba de agua", "caldera", "ascensor", "transformador"}
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
			if activosData, exists := data["activos"]; exists {
                if activosArray, ok := activosData.([]interface{}); ok {
                    mu.Lock()
                    activos = append(activos, activosArray...)
                    mu.Unlock()
                }
            }
		}(tipo)
	}

	wg.Wait() // Esperar a que todas las goroutines terminen

	result := map[string]interface{}{
        "activos": activos,
        "total": len(activos),
        "timestamp": time.Now(),
    }

	c.JSON(http.StatusOK, result)
}