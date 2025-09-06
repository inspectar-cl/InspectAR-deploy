# 🧠 Tutorial: Cómo Agregar Endpoints Compuestos

> **API Gateway InspectAR** - Guía paso a paso para crear endpoints que combinan múltiples microservicios

## 📖 Introducción

Los **endpoints compuestos** son endpoints inteligentes que:
- Realizan múltiples consultas a diferentes microservicios
- Procesan y transforman las respuestas
- Combinan los datos en una sola respuesta unificada
- Manejan consultas concurrentes para mejor rendimiento

**Ejemplo práctico:** El endpoint `/api/calderas-estado` que:
1. Obtiene todas las calderas del servicio de gestión
2. Para cada caldera, consulta el estado de sus sensores
3. Para cada caldera, obtiene los últimos valores de sus sensores
4. Combina todo en una respuesta unificada con `ultimo_valor` en cada sensor

## 🏗️ Arquitectura del Sistema

```
endpoints_config.go → composite_handler.go → transforms.go
    (QUÉ hacer)         (CÓMO ejecutar)      (CÓMO transformar)
```

### Flujo de Ejecución:
1. **Configuración** (`endpoints_config.go`) - Define qué consultas hacer
2. **Ejecución** (`composite_handler.go`) - Motor que ejecuta las consultas automáticamente
3. **Transformación** (`transforms.go`) - Procesa y combina las respuestas

---

## 🚀 Caso de Uso: Endpoint Calderas-Estado

Vamos a analizar el endpoint existente `/api/calderas-estado` como ejemplo:

### 📝 Objetivo del Endpoint
Crear un endpoint que devuelva:
- Todas las calderas del sistema
- Estado de cada sensor de cada caldera
- Último valor registrado de cada sensor
- Todo en una sola respuesta al frontend

### 🔄 Flujo de Consultas
```
1. GET /gestion/activos/tipo/caldera
   ↓
2. Para cada caldera concurrentemente:
   - GET /parser/activo/{activo_id}/sensores/estado
   - GET /parser/lectura/{activo_id}/datos/ultimo
   ↓
3. Combinar respuestas y enviar al frontend
```

---

## 📋 Paso a Paso: Implementación

### **Paso 1: Configurar el Endpoint**

**Archivo:** `handlers/endpoints_config.go`

```go
var CompositeEndpoints = []EndpointConfig{
    // ... otros endpoints existentes ...
    
    {
        Path:        "/api/calderas-estado",
        Method:      "GET", 
        Handler:     "GetCalderasEstado",
        Description: "Obtiene todas las calderas y el estado de sus sensores",
        Queries: []ServiceQuery{
            {
                Name:         "calderas",
                TargetEnvVar: "GESTION_URL",
                Path:         "/activos/tipo/caldera",
                Method:       "GET",
                Transform:    "TransformCalderas",
            },
            // Las consultas adicionales (sensores/lecturas) se manejan 
            // dentro de TransformCalderas por ser dinámicas
        },
    },
}
```

**¿Por qué solo una query?** 
- Las consultas de sensores y lecturas son **dinámicas** (una por cada caldera)
- Se ejecutan dentro de la función de transformación
- No se pueden definir estáticamente en la configuración

### **Paso 2: Crear la Función de Transformación**

**Archivo:** `handlers/transforms.go`

```go
// Función principal que orquesta las consultas
func transformCalderas(data interface{}) (interface{}, error) {
    if dataMap, ok := data.(map[string]interface{}); ok {
        // Agregar timestamp de procesamiento
        dataMap["processed_at"] = time.Now()
        
        // Ejecutar consultas adicionales para cada caldera
        return executeCalderasSensoresQueries(dataMap)
    }
    
    return data, nil
}

// Función que ejecuta las consultas concurrentes por caldera
func executeCalderasSensoresQueries(calderasData interface{}) (interface{}, error) {
    // 1. Extraer calderas de la respuesta anterior
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
    
    // 2. Configurar variables para consultas concurrentes
    var wg sync.WaitGroup
    var mu sync.Mutex
    
    calderasConSensores := []map[string]interface{}{}
    errores := []string{}
    
    parserURL := os.Getenv("PARSER_URL")
    if parserURL == "" {
        return nil, fmt.Errorf("variable de entorno PARSER_URL no definida")
    }
    
    client := &http.Client{Timeout: 15 * time.Second}
    
    // 3. Ejecutar consultas concurrentes para cada caldera
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
            
            // 4. HACER DOS CONSULTAS CONCURRENTES POR CALDERA
            var wgCaldera sync.WaitGroup
            var sensoresData, ultimosValores interface{}
            var errSensores, errUltimos error
            
            // Consulta A: Estado de sensores
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
            
            // Consulta B: Últimos valores de sensores
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
            
            // 5. Esperar ambas consultas
            wgCaldera.Wait()
            
            // 6. Verificar errores
            if errSensores != nil && errUltimos != nil {
                mu.Lock()
                errores = append(errores, fmt.Sprintf("Error en ambas consultas para %s: sensores=%v, valores=%v", activoID, errSensores, errUltimos))
                mu.Unlock()
                return
            }
            
            // 7. Combinar datos de caldera con sensores y valores
            calderaConSensores := make(map[string]interface{})
            
            // Copiar todos los datos de la caldera original
            for key, value := range c {
                calderaConSensores[key] = value
            }
            
            // 8. COMBINAR SENSORES CON SUS ÚLTIMOS VALORES
            if errSensores == nil {
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
                                
                                // ⭐ CLAVE: Agregar ultimo_valor a cada sensor
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
            
            // Manejar error solo de valores
            if errUltimos != nil && errSensores == nil {
                calderaConSensores["valores_error"] = errUltimos.Error()
            }
            
            calderaConSensores["consulta_timestamp"] = time.Now()
            
            // 9. Agregar resultado thread-safe
            mu.Lock()
            calderasConSensores = append(calderasConSensores, calderaConSensores)
            mu.Unlock()
            
        }(caldera)
    }
    
    // 10. Esperar todas las goroutines
    wg.Wait()
    
    // 11. Construir respuesta final
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
```

### **Paso 3: Registrar la Función de Transformación**

**Archivo:** `handlers/composite_handler.go`

```go
func NewCompositeHandler() *CompositeHandler {
    return &CompositeHandler{
        httpClient: &http.Client{
            Timeout: 30 * time.Second,
        },
        transforms: map[string]TransformFunc{
            "TransformCalderas":       transformCalderas,        // ⭐ Agregar aquí
            "TransformSensoresEstado": transformSensoresEstado,
            // Agregar más transformaciones aquí...
        },
    }
}
```

---

## 🎯 Resultado Final

### Respuesta del Endpoint `/api/calderas-estado`:

```json
{
    "calderas": {
        "activos": [
            {
                "activo_id": "AC-1001",
                "nombre": "Caldera Principal",
                "estado": "operativo",
                "tipo": "caldera",
                "ubicacion": "Sala de Calderas 1",
                "sensores_info": {
                    "activo_id": "AC-1001",
                    "estado": "operativo",
                    "nombre": "Caldera Principal",
                    "total_sensores": 1,
                    "sensores": [
                        {
                            "sensor_id": "PRES001",
                            "tipo": "presion",
                            "estado": "connected",
                            "is_active": true,
                            "unidad": "bar",
                            "ultimo_valor": {        // ⭐ AQUÍ ESTÁ LA CLAVE
                                "tiempo": "2025-09-04T12:30:47Z",
                                "valor": 85.4
                            }
                        }
                    ],
                    "resumen": {
                        "sensores_activos": 1,
                        "sensores_desconectados": 0,
                        "sensores_nunca_conectados": 0
                    }
                },
                "consulta_timestamp": "2025-09-06T04:26:36.585771347Z"
            }
        ],
        "total_calderas": 1,
        "consultas_exitosas": 1,
        "consultas_fallidas": 0,
        "processed_at": "2025-09-06T04:26:36.585774536Z"
    },
    "timestamp": "2025-09-06T04:26:36.585780541Z"
}
```

---

## 📝 Template para Nuevos Endpoints

### Ejemplo: Crear `/api/bombas-status`

#### **1. Configurar Endpoint** (`endpoints_config.go`):

```go
{
    Path:        "/api/bombas-status",
    Method:      "GET",
    Handler:     "GetBombasStatus", 
    Description: "Obtiene todas las bombas con datos de rendimiento",
    Queries: []ServiceQuery{
        {
            Name:         "bombas",
            TargetEnvVar: "GESTION_URL",
            Path:         "/activos/tipo/bomba",
            Method:       "GET",
            Transform:    "TransformBombas", // ⭐ Función a crear
        },
    },
},
```

#### **2. Crear Transformación** (`transforms.go`):

```go
// Función principal para bombas
func transformBombas(data interface{}) (interface{}, error) {
    if dataMap, ok := data.(map[string]interface{}); ok {
        dataMap["processed_at"] = time.Now()
        return executeBombasQueries(dataMap)
    }
    return data, nil
}

// Función que ejecuta consultas específicas de bombas
func executeBombasQueries(bombasData interface{}) (interface{}, error) {
    // 1. Extraer bombas
    bombas := []map[string]interface{}{}
    
    // 2. Configurar consultas concurrentes
    // 3. Para cada bomba hacer:
    //    - GET /parser/activo/{activo_id}/rendimiento
    //    - GET /parser/activo/{activo_id}/mantenimiento
    // 4. Combinar resultados
    // 5. Retornar respuesta unificada
    
    // ... implementación similar a executeCalderasSensoresQueries ...
    
    return result, nil
}
```

#### **3. Registrar Transformación** (`composite_handler.go`):

```go
transforms: map[string]TransformFunc{
    "TransformCalderas": transformCalderas,
    "TransformBombas":   transformBombas,    // ⭐ Agregar nueva función
    // ...
},
```

---

## 🔧 Mejores Prácticas

### ✅ **Do's (Hacer):**

1. **Usar consultas concurrentes** para mejor rendimiento
2. **Manejar errores independientes** por cada consulta
3. **Agregar timeouts apropiados** (15-30 segundos)
4. **Usar mutexes** para operaciones thread-safe
5. **Validar datos** antes de procesarlos
6. **Agregar timestamps** para debugging
7. **Normalizar estados** cuando sea necesario
8. **Documentar** qué hace cada función

### ❌ **Don'ts (No hacer):**

1. **No hacer consultas secuenciales** innecesarias
2. **No hardcodear URLs** - usar variables de entorno
3. **No ignorar errores** de los servicios
4. **No olvidar cerrar** response bodies
5. **No asumir** que los datos siempre vienen correctos
6. **No crear endpoints compuestos** para consultas simples
7. **No duplicar lógica** entre diferentes transformaciones

---

## 🧪 Testing

### Probar el Endpoint:

```bash
# Consulta básica
curl -s http://localhost:3500/api/calderas-estado

# Con formato JSON
curl -s http://localhost:3500/api/calderas-estado | python3 -m json.tool

# Verificar campos específicos
curl -s http://localhost:3500/api/calderas-estado | jq '.calderas.activos[0].sensores_info.sensores[0].ultimo_valor'
```

### Health Check:

```bash
curl http://localhost:3500/healthz
# Debe retornar: {"ok": true}
```

---

## 🚀 Conclusión

Los endpoints compuestos son ideales para:

- **Reducir llamadas** del frontend al backend
- **Combinar datos** de múltiples microservicios
- **Mejorar rendimiento** con consultas concurrentes
- **Simplificar** la lógica del frontend
- **Centralizar** transformaciones de datos

**Patrón recomendado:**
1. Una consulta principal en `endpoints_config.go`
2. Consultas adicionales dinámicas en `transforms.go`
3. Combinación de resultados en la misma función
4. Manejo de errores independiente por consulta

**¿Próximos pasos?** Usa este tutorial para crear tus propios endpoints compuestos siguiendo el mismo patrón! 🎯
