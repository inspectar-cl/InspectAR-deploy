package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

// transformCalderas procesa la respuesta del servicio de gestión de calderas
// y automáticamente obtiene los datos de sensores para cada caldera
func transformCalderas(data interface{}) (interface{}, error) {
	// La data viene del endpoint /gestion/activos/tipo/caldera
	if dataMap, ok := data.(map[string]interface{}); ok {
		// Agregar timestamp de procesamiento
		dataMap["processed_at"] = time.Now()
		
		// Normalizar el formato si es necesario
		if activos, exists := dataMap["activos"]; exists {
			if activosArray, ok := activos.([]interface{}); ok {
				for i, activo := range activosArray {
					if activoMap, ok := activo.(map[string]interface{}); ok {
						// Normalizar estado
						if estado, exists := activoMap["estado"]; exists {
							activoMap["estado_normalizado"] = normalizeEstado(fmt.Sprintf("%v", estado))
						}
						activosArray[i] = activoMap
					}
				}
				dataMap["activos"] = activosArray
			}
		}
		
		// Ahora obtener los sensores para cada caldera
		return executeCalderasSensoresQueries(dataMap)
	}
	
	return data, nil
}

// transformSensoresEstado procesa la respuesta de sensores y hace las consultas individuales
func transformSensoresEstado(data interface{}) (interface{}, error) {
	// Esta función se ejecuta después de obtener las calderas
	// Necesitamos hacer consultas individuales para cada caldera
	
	// Esta es una implementación especial para el caso de dependencias
	// En lugar de transformar una respuesta única, orquesta múltiples llamadas
	return executeCalderasSensoresQueries(data)
}

func executeCalderasSensoresQueries(calderasData interface{}) (interface{}, error) {
	// Extraer las calderas de la respuesta anterior
	calderas := []map[string]interface{}{}
	
	if dataMap, ok := calderasData.(map[string]interface{}); ok {
		if activosArray, exists := dataMap["activos"]; exists {
			if activos, ok := activosArray.([]interface{}); ok {
				for _, activo := range activos {
					if activoMap, ok := activo.(map[string]interface{}); ok {
						calderas = append(calderas, activoMap)
					}
				}
			}
		}
	}
	
	if len(calderas) == 0 {
		return map[string]interface{}{
			"calderas_con_sensores": []interface{}{},
			"total_calderas": 0,
		}, nil
	}
	
	// Hacer consultas concurrentes para cada caldera
	var wg sync.WaitGroup
	var mu sync.Mutex
	
	calderasConSensores := []map[string]interface{}{}
	errores := []string{}
	
	parserURL := os.Getenv("PARSER_URL")
	if parserURL == "" {
		return nil, fmt.Errorf("variable de entorno PARSER_URL no definida")
	}
	
	client := &http.Client{Timeout: 15 * time.Second}
	
	for _, caldera := range calderas {
		wg.Add(1)
		go func(c map[string]interface{}) {
			defer wg.Done()
			
			activoID, ok := c["activo_id"].(string)
			if !ok {
				mu.Lock()
				errores = append(errores, fmt.Sprintf("activo_id no válido para caldera %v", c["nombre"]))
				mu.Unlock()
				return
			}
			
			// 🆕 HACER DOS CONSULTAS CONCURRENTES POR CALDERA
			var wgCaldera sync.WaitGroup
			var sensoresData, ultimosValores interface{}
			var errSensores, errUltimos error
			
			// Consulta 1: Estado de sensores
			wgCaldera.Add(1)
			go func() {
				defer wgCaldera.Done()
				sensorURL := fmt.Sprintf("%s/activo/%s/sensores/estado", parserURL, activoID)
				resp, err := client.Get(sensorURL)
				if err != nil {
					errSensores = fmt.Errorf("error consultando sensores: %v", err)
					return
				}
				defer resp.Body.Close()
				
				if resp.StatusCode >= 400 {
					errSensores = fmt.Errorf("error del servicio sensores: %d", resp.StatusCode)
					return
				}
				
				if err := json.NewDecoder(resp.Body).Decode(&sensoresData); err != nil {
					errSensores = fmt.Errorf("error decodificando sensores: %v", err)
					return
				}
			}()
			
			// Consulta 2: Últimos valores de sensores
			wgCaldera.Add(1)
			go func() {
				defer wgCaldera.Done()
				ultimosURL := fmt.Sprintf("%s/lectura/%s/datos/ultimo", parserURL, activoID)
				resp, err := client.Get(ultimosURL)
				if err != nil {
					errUltimos = fmt.Errorf("error consultando últimos valores: %v", err)
					return
				}
				defer resp.Body.Close()
				
				if resp.StatusCode >= 400 {
					errUltimos = fmt.Errorf("error del servicio últimos valores: %d", resp.StatusCode)
					return
				}
				
				if err := json.NewDecoder(resp.Body).Decode(&ultimosValores); err != nil {
					errUltimos = fmt.Errorf("error decodificando últimos valores: %v", err)
					return
				}
			}()
			
			// Esperar ambas consultas
			wgCaldera.Wait()
			
			// Verificar errores
			if errSensores != nil && errUltimos != nil {
				mu.Lock()
				errores = append(errores, fmt.Sprintf("Error en ambas consultas para %s: sensores=%v, valores=%v", activoID, errSensores, errUltimos))
				mu.Unlock()
				return
			}
			
			// Combinar datos de caldera con sensores y valores
			calderaConSensores := make(map[string]interface{})
			
			// Copiar todos los datos de la caldera
			for key, value := range c {
				calderaConSensores[key] = value
			}
			
			// Combinar sensores con sus últimos valores
			if errSensores == nil {
				// Obtener los datos de sensores
				sensoresInfo := sensoresData
				
				// Si también tenemos los últimos valores, combinarlos
				if errUltimos == nil {
					if sensoresMap, ok := sensoresInfo.(map[string]interface{}); ok {
						if sensoresArray, exists := sensoresMap["sensores"]; exists {
							if sensores, ok := sensoresArray.([]interface{}); ok {
								// Convertir últimos valores a map para fácil búsqueda
								ultimosMap := make(map[string]interface{})
								if ultimosData, ok := ultimosValores.(map[string]interface{}); ok {
									ultimosMap = ultimosData
								}
								
								// Agregar ultimo_valor a cada sensor
								for i, sensor := range sensores {
									if sensorMap, ok := sensor.(map[string]interface{}); ok {
										if sensorID, exists := sensorMap["sensor_id"]; exists {
											if sensorIDStr, ok := sensorID.(string); ok {
												if valorData, exists := ultimosMap[sensorIDStr]; exists {
													sensorMap["ultimo_valor"] = valorData
												} else {
													sensorMap["ultimo_valor"] = nil
												}
											}
										}
										sensores[i] = sensorMap
									}
								}
								sensoresMap["sensores"] = sensores
							}
						}
					}
				}
				
				calderaConSensores["sensores_info"] = sensoresInfo
			} else {
				calderaConSensores["sensores_error"] = errSensores.Error()
			}
			
			// Solo agregar error de valores si falló la consulta de valores pero sensores funcionó
			if errUltimos != nil && errSensores == nil {
				calderaConSensores["valores_error"] = errUltimos.Error()
			}
			
			calderaConSensores["consulta_timestamp"] = time.Now()
			
			mu.Lock()
			calderasConSensores = append(calderasConSensores, calderaConSensores)
			mu.Unlock()
			
		}(caldera)
	}
	
	wg.Wait()
	
	result := map[string]interface{}{
		"activos": calderasConSensores,
		"total_calderas": len(calderas),
		"consultas_exitosas": len(calderasConSensores),
		"consultas_fallidas": len(errores),
		"processed_at": time.Now(),
	}
	
	if len(errores) > 0 {
		result["errores"] = errores
	}
	
	return result, nil
}

func normalizeEstado(estado string) string {
	switch strings.ToLower(estado) {
	case "ok", "operativo", "active", "running":
		return "operativo"
	case "error", "fallido", "failed", "stopped":
		return "fallido"
	case "mantenimiento", "maintenance":
		return "mantenimiento"
	case "desconectado", "disconnected", "offline":
		return "desconectado"
	default:
		return "desconocido"
	}
}
