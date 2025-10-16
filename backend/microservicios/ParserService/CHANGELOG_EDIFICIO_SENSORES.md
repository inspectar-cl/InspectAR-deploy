# Changelog - Ruta /activo/edificio/:edificio_id

## 🎉 Versión 2.0 - Información Completa de Sensores (11 de Octubre 2025)

### ✨ Nuevas Características

#### Información Enriquecida de Sensores
La ruta `/activo/edificio/:edificio_id` ahora retorna información completa del estado de cada sensor asociado a los activos del edificio.

**Cambio Principal:**
- **Antes**: Solo retornaba la lista básica de sensores (`sensor_id`, `tipo`, `unidad`)
- **Ahora**: Incluye estado de conexión, métricas y timestamps de cada sensor

### 📋 Estructura de Respuesta Actualizada

```json
{
  "edificio_id": 1,
  "activos": [
    {
      "id": "...",
      "activo_id": 1,
      "estado": "operativo",
      "id_edificio": 1,
      "total_sensores": 2,
      "sensores": [
        {
          "sensor_id": "TEMP_001",
          "tipo": "temperatura",
          "unidad": "°C",
          "estado": "connected",
          "is_active": true,
          "last_seen": "2025-10-11T12:00:00Z",
          "first_seen": "2025-10-01T08:00:00Z",
          "total_reports": 1458,
          "created_at": "2025-10-01T08:00:00Z",
          "updated_at": "2025-10-11T12:00:00Z"
        }
      ]
    }
  ],
  "total": 2
}
```

### 🔍 Campos Agregados por Sensor

| Campo | Tipo | Descripción |
|-------|------|-------------|
| `estado` | string | Estado del sensor: `connected`, `disconnected`, `never_connected` |
| `is_active` | boolean | Indica si el sensor está activo actualmente |
| `last_seen` | timestamp | Última vez que el sensor envió datos |
| `first_seen` | timestamp | Primera vez que el sensor envió datos |
| `total_reports` | int | Número total de reportes enviados por el sensor |
| `created_at` | timestamp | Fecha de creación del registro del sensor |
| `updated_at` | timestamp | Fecha de última actualización del estado |

### 🎯 Estados de Sensores

- **`connected`**: Sensor activo enviando datos regularmente
- **`disconnected`**: Sensor que envió datos anteriormente pero está inactivo (>5 min sin datos)
- **`never_connected`**: Sensor registrado pero que nunca ha enviado datos

### 🛠️ Implementación Técnica

**Archivo modificado:**
- `/internal/handler/dataHandler.go` - Función `GetActivosByEdificio()`

**Lógica implementada:**
1. Obtener activos del edificio desde MongoDB
2. Para cada activo, obtener sus sensores
3. Para cada sensor, consultar su estado en el sistema de monitoreo
4. Enriquecer la respuesta con toda la información de estado
5. Manejar casos donde no hay información de estado (sensores nuevos)

### 📊 Casos de Uso

#### 1. Dashboard de Edificios
Permite visualizar en tiempo real el estado de todos los sensores por edificio:
```bash
curl http://localhost:8090/activo/edificio/1 | jq '.activos[] | {
  activo_id,
  sensores_activos: [.sensores[] | select(.estado == "connected")] | length,
  sensores_totales: .total_sensores
}'
```

#### 2. Detección de Problemas
Identificar sensores desconectados por edificio:
```bash
curl http://localhost:8090/activo/edificio/1 | jq '[
  .activos[].sensores[] | 
  select(.estado == "disconnected")
]'
```

#### 3. Inventario de Sensores
Listar todos los sensores de un edificio con su estado:
```bash
curl http://localhost:8090/activo/edificio/1 | jq '[
  .activos[] | {
    activo: .activo_id,
    sensores: [.sensores[] | {
      id: .sensor_id,
      tipo: .tipo,
      estado: .estado,
      reportes: .total_reports
    }]
  }
]'
```

### ✅ Tests

Se ha creado un test completo para validar la funcionalidad:

**Archivo:** `/tests/test_activos_edificio_con_sensores.sh`

**Pruebas incluidas:**
- ✅ Creación de activos en múltiples edificios
- ✅ Envío de lecturas de sensores
- ✅ Verificación de estado de sensores (connected/disconnected/never_connected)
- ✅ Validación de estructura de respuesta
- ✅ Filtrado correcto por edificio
- ✅ Campos de estado presentes en cada sensor
- ✅ Contadores de sensores correctos

**Ejecutar test:**
```bash
cd /home/joytan/repo/InspectAR/backend/microservicios/ParserService/tests
./test_activos_edificio_con_sensores.sh
```

### 🔄 Compatibilidad

**Retrocompatibilidad:** ✅ Mantenida
- La estructura básica de la respuesta no cambió
- Solo se agregaron campos adicionales
- Clientes existentes funcionarán sin modificaciones
- Los nuevos campos pueden ser ignorados por clientes antiguos

### 📝 Notas de Desarrollo

**Performance:**
- La función ahora realiza consultas adicionales para obtener el estado de sensores
- Optimización: Se consulta el estado de cada sensor individualmente
- Consideración futura: Implementar cache para estados de sensores frecuentemente consultados

**Manejo de Errores:**
- Si no se puede obtener el estado de un sensor, se muestra estado "unknown"
- Si un activo no tiene sensores, retorna array vacío con `total_sensores: 0`
- La ruta es tolerante a fallos parciales

### 🚀 Próximas Mejoras

- [ ] Implementar paginación para edificios con muchos activos
- [ ] Agregar filtros adicionales (por estado de sensor, tipo de sensor)
- [ ] Cache de estados de sensores para mejorar performance
- [ ] Agregar resumen agregado de sensores por edificio
- [ ] WebSocket para updates en tiempo real del estado

### 📚 Documentación Actualizada

- ✅ README.md actualizado con nueva estructura de respuesta
- ✅ Ejemplos de uso con jq agregados
- ✅ Tests documentados
- ✅ Este CHANGELOG creado

---

**Desarrollado por el equipo de InspectAR** 🚀
