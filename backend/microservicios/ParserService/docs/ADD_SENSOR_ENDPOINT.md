# 🔧 Nueva Funcionalidad: Agregar Sensor a Activo

## Descripción
Se ha implementado una nueva ruta en el ParserService que permite agregar sensores a un activo existente.

## Endpoint

```
POST /activo/:activo_id/sensores
```

## Parámetros

### URL Parameters
- `activo_id` (string, requerido): ID del activo al cual se agregará el sensor

### Request Body
```json
{
    "sensor_id": "string",  // ID único del sensor (requerido)
    "tipo": "string",       // Tipo de sensor (requerido)
    "unidad": "string"      // Unidad de medida (requerido)
}
```

## Ejemplos de Uso

### ✅ Agregar sensor de temperatura
```bash
curl -X POST http://localhost:8090/activo/AC-1001/sensores \
  -H "Content-Type: application/json" \
  -d '{
    "sensor_id": "TEMP_001",
    "tipo": "temperatura",
    "unidad": "°C"
  }'
```

**Respuesta:**
```json
{
  "mensaje": "Sensor agregado correctamente al activo",
  "activo_id": "AC-1001",
  "sensor": {
    "sensor_id": "TEMP_001",
    "tipo": "temperatura",
    "unidad": "°C"
  }
}
```

### ✅ Agregar sensor de presión
```bash
curl -X POST http://localhost:8090/activo/AC-1002/sensores \
  -H "Content-Type: application/json" \
  -d '{
    "sensor_id": "PRESS_001",
    "tipo": "presion",
    "unidad": "bar"
  }'
```

### ✅ Agregar sensor de vibración
```bash
curl -X POST http://localhost:8090/activo/AC-1003/sensores \
  -H "Content-Type: application/json" \
  -d '{
    "sensor_id": "VIBR_001",
    "tipo": "vibracion",
    "unidad": "Hz"
  }'
```

## Respuestas

### 200 OK - Sensor agregado exitosamente
```json
{
  "mensaje": "Sensor agregado correctamente al activo",
  "activo_id": "AC-1001",
  "sensor": {
    "sensor_id": "TEMP_001",
    "tipo": "temperatura",
    "unidad": "°C"
  }
}
```

### 400 Bad Request - Datos inválidos
```json
{
  "error": "Los campos sensor_id, tipo y unidad son requeridos"
}
```

### 404 Not Found - Activo no encontrado
```json
{
  "error": "Activo no encontrado"
}
```

### 500 Internal Server Error - Error interno
```json
{
  "error": "No se pudo agregar el sensor al activo"
}
```

## Validaciones

- ✅ El `activo_id` debe existir en la base de datos
- ✅ El campo `sensor_id` es requerido y no puede estar vacío
- ✅ El campo `tipo` es requerido y no puede estar vacío
- ✅ El campo `unidad` es requerido y no puede estar vacío
- ✅ El JSON debe tener formato válido

## Verificación

Para verificar que el sensor se agregó correctamente, puedes consultar el activo:

```bash
curl http://localhost:8090/activo/AC-1001
```

O verificar el estado de los sensores:

```bash
curl http://localhost:8090/activo/AC-1001/sensores/estado
```

## Código de Implementación

### Handler
```go
// POST /activo/:activo_id/sensores - Agregar sensor a un activo existente
func (h *DataHandler) AddSensorToActivo(c *gin.Context) {
	activoID := c.Param("activo_id")
	
	var sensor models.Sensor
	if err := c.ShouldBindJSON(&sensor); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos del sensor inválidos"})
		return
	}

	// Validar campos requeridos
	if sensor.SensorID == "" || sensor.Tipo == "" || sensor.Unidad == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Los campos sensor_id, tipo y unidad son requeridos"})
		return
	}

	// Validar que el activo existe
	_, err := h.activoService.ObtenerActivo(c.Request.Context(), activoID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Activo no encontrado"})
		return
	}

	// Agregar el sensor al activo
	err = h.activoService.AgregarSensor(c.Request.Context(), activoID, sensor)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo agregar el sensor al activo"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"mensaje": "Sensor agregado correctamente al activo",
		"activo_id": activoID,
		"sensor": sensor,
	})
}
```

### Ruta
```go
r.POST("/activo/:activo_id/sensores", dataHandler.AddSensorToActivo)
```

## Test Automatizado

Se ha creado un test automatizado completo en:
```
/backend/microservicios/ParserService/tests/test_add_sensor_to_activo.sh
```

Para ejecutarlo:
```bash
cd /backend && ./microservicios/ParserService/tests/test_add_sensor_to_activo.sh
```

## Estado de la Implementación

✅ **COMPLETADO** - La funcionalidad está implementada y funcionando correctamente
✅ **VALIDACIONES** - Todas las validaciones implementadas
✅ **TESTS** - Tests automatizados creados y funcionando
✅ **DOCUMENTACIÓN** - Documentación completa disponible

---

**Fecha de implementación:** 4 de Septiembre, 2025  
**Desarrollado por:** GitHub Copilot  
**Microservicio:** ParserService  
**Puerto:** 8090  
