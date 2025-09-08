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

type CompositeHandler struct {
	httpClient *http.Client
	transforms map[string]TransformFunc
}

type TransformFunc func(interface{}) (interface{}, error)

func NewCompositeHandler() *CompositeHandler {
	return &CompositeHandler{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		transforms: map[string]TransformFunc{
			"TransformCalderas":       transformCalderas,
			"TransformSensoresEstado": transformSensoresEstado,
		},
	}
}

// Handler genérico que ejecuta las consultas definidas en la configuración
func (h *CompositeHandler) ExecuteCompositeEndpoint(config EndpointConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		results := make(map[string]interface{})
		errors := make(map[string]error)

		// Ejecutar queries con dependencias
		err := h.executeQueriesWithDependencies(config.Queries, results, errors)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   true,
				"message": "Error ejecutando consultas",
				"details": err.Error(),
			})
			return
		}

		// Procesar resultado final según el handler específico
		finalResult := h.processResults(config.Handler, results, errors)

		// Determinar código de estado
		statusCode := http.StatusOK
		if len(errors) > 0 {
			if len(results) == 0 {
				statusCode = http.StatusServiceUnavailable
			} else {
				statusCode = http.StatusPartialContent
			}
		}

		c.JSON(statusCode, finalResult)
	}
}

func (h *CompositeHandler) executeQueriesWithDependencies(queries []ServiceQuery, results map[string]interface{}, errors map[string]error) error {
	// Separar queries independientes de las dependientes
	independent := []ServiceQuery{}
	dependent := []ServiceQuery{}

	for _, query := range queries {
		if query.DependsOn == "" {
			independent = append(independent, query)
		} else {
			dependent = append(dependent, query)
		}
	}

	// Ejecutar queries independientes concurrentemente
	if len(independent) > 0 {
		var wg sync.WaitGroup
		var mu sync.Mutex

		for _, query := range independent {
			wg.Add(1)
			go func(q ServiceQuery) {
				defer wg.Done()
				result, err := h.executeQuery(q, results)

				mu.Lock()
				defer mu.Unlock()

				if err != nil {
					errors[q.Name] = err
				} else {
					// Aplicar transformación si está definida
					if q.Transform != "" {
						if transformFunc, exists := h.transforms[q.Transform]; exists {
							if transformed, err := transformFunc(result); err == nil {
								result = transformed
							}
						}
					}
					results[q.Name] = result
				}
			}(query)
		}
		wg.Wait()
	}

	// Ejecutar queries dependientes secuencialmente
	for _, query := range dependent {
		if _, hasError := errors[query.DependsOn]; hasError {
			continue // Skip si la dependencia falló
		}

		result, err := h.executeQuery(query, results)
		if err != nil {
			errors[query.Name] = err
			continue
		}

		// Aplicar transformación si está definida
		if query.Transform != "" {
			if transformFunc, exists := h.transforms[query.Transform]; exists {
				if transformed, err := transformFunc(result); err == nil {
					result = transformed
				}
			}
		}
		results[query.Name] = result
	}

	return nil
}

func (h *CompositeHandler) executeQuery(query ServiceQuery, currentResults map[string]interface{}) (interface{}, error) {
	target := os.Getenv(query.TargetEnvVar)
	if target == "" {
		return nil, fmt.Errorf("variable de entorno %s no definida", query.TargetEnvVar)
	}

	// Construir la URL
	var url string
	if query.PathTemplate != "" {
		// Es una query dependiente con template
		url = h.buildURLFromTemplate(target, query.PathTemplate, query.DependsOn, currentResults)
	} else {
		// Es una query simple
		url = fmt.Sprintf("%s%s", target, query.Path)
	}

	req, err := http.NewRequest(query.Method, url, nil)
	if err != nil {
		return nil, err
	}

	// Agregar headers personalizados
	for key, value := range query.Headers {
		req.Header.Set(key, value)
	}

	resp, err := h.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("error del servicio: %d", resp.StatusCode)
	}

	var result interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

func (h *CompositeHandler) buildURLFromTemplate(baseURL, template, dependsOn string, results map[string]interface{}) string {
	// Para el caso específico de sensores que depende de calderas
	if dependsOn == "calderas" {
		// Esta función se llamará múltiples veces, una por cada caldera
		// Por ahora, construiremos las URLs dinámicamente en la función de transformación
		return fmt.Sprintf("%s%s", baseURL, template)
	}

	return fmt.Sprintf("%s%s", baseURL, template)
}

func (h *CompositeHandler) processResults(handlerName string, results map[string]interface{}, errors map[string]error) interface{} {
	switch handlerName {
	case "GetCalderasEstado":
		return h.processCalderasEstado(results, errors)
	default:
		// Resultado genérico
		response := map[string]interface{}{
			"data": results,
		}
		if len(errors) > 0 {
			response["errors"] = errors
		}
		return response
	}
}

func (h *CompositeHandler) processCalderasEstado(results map[string]interface{}, errors map[string]error) interface{} {
	response := map[string]interface{}{
		"calderas": make([]interface{}, 0),
		"timestamp": time.Now(),
	}

	if calderas, ok := results["calderas"]; ok {
		response["calderas"] = calderas
	}

	if sensoresData, ok := results["sensores_estado"]; ok {
		response["sensores_por_caldera"] = sensoresData
	}

	if len(errors) > 0 {
		response["errors"] = errors
	}

	return response
}
