# Tests del Microservicio de Gestión

## Descripción
Esta carpeta contiene todos los tests para el microservicio de gestión, incluyendo pruebas de integración y funcionales.

## Archivos

### `test_gestion_service.sh`
Script completo de testing que valida todas las rutas y funcionalidades del servicio de gestión.

**Funcionalidades testadas:**
- ✅ Health check y rutas básicas
- ✅ Sistema de técnicos (HdU16) - CRUD completo
- ✅ Acciones de mantenimiento (HdU13) - CRUD completo  
- ✅ Reportes automáticos (HdU04) - Generación y consulta
- ✅ **Solicitudes técnicas (HdU16) - FLUJO COMPLETO DE GESTIÓN**
- ✅ Edge cases y validaciones

## Uso

### Prerrequisitos
1. El servicio de gestión debe estar ejecutándose en `http://localhost:8092`
2. La base de datos PostgreSQL debe estar disponible y poblada con datos de prueba

### Ejecución
```bash
# Desde la carpeta test
cd /path/to/backend/microservicios/gestion/test
./test_gestion_service.sh
```

### Resultado
El script ejecuta **69 tests** y proporciona:
- ✅ Resumen de éxito/fallo por categoría
- 📊 Estadísticas detalladas de cobertura
- 🎯 **Flujo completo de solicitudes** con 5 escenarios realistas
- ⚠️ Identificación de edge cases problemáticos

## Flujo de Solicitudes Testado

El script incluye un **flujo completo de gestión de solicitudes** que simula:

1. **📝 CREACIÓN**: 5 solicitudes diferentes
   - Mantenimiento preventivo de caldera (alta prioridad)
   - Emergencia eléctrica (urgente)
   - Inspección sistema HVAC (media prioridad)
   - Consulta técnica especializada (baja prioridad)
   - Reparación compresor industrial (alta prioridad)

2. **📤 ENVÍO**: Asignación inteligente a técnicos según especialidad

3. **🔄 SEGUIMIENTO**: Transiciones de estado realistas
   - `pendiente` → `recibida` → `en_proceso` → `completada`/`cancelada`

4. **✅ RESOLUCIÓN**: Finalización con comentarios técnicos

5. **📊 MÉTRICAS**: Estadísticas de rendimiento del sistema

## APIs Testadas

### Solicitudes (API v1)
- `GET/POST /api/v1/solicitudes` - Crear y listar
- `POST /api/v1/solicitudes/{id}/enviar` - Asignar técnicos  
- `PUT /api/v1/solicitudes/{id}/estado` - Seguimiento
- `GET /api/v1/solicitudes/{id}` - Consulta individual
- `GET /api/v1/solicitudes/estadisticas` - Métricas

### Otras APIs
- Técnicos, acciones de mantenimiento, reportes
- Filtros y parámetros de consulta
- Validaciones y edge cases

## Notas
- El script está diseñado para ser **idempotente** y puede ejecutarse múltiples veces
- Incluye **limpieza automática** de datos de prueba
- Proporciona **output coloreado** para facilitar la lectura
- **Success rate esperado**: ~94% (65/69 tests)
