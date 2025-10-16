# Base de Datos IA - PostgreSQL

Base de datos PostgreSQL dedicada al almacenamiento de datos relacionados con Machine Learning e Inteligencia Artificial.

## 📊 Tablas

### `anomalias`
Almacena las anomalías detectadas por el sistema de IA en los datos de sensores.

**Campos:**
- `id` (SERIAL PRIMARY KEY): Identificador único
- `activo_id` (INTEGER): ID del activo relacionado
- `sensor_id` (VARCHAR(100)): Identificador del sensor
- `timestamp` (TIMESTAMP): Momento de la detección
- `anomaly_score` (FLOAT): Puntuación de anomalía (0-1)
- `anomaly_likelihood` (FLOAT): Probabilidad de anomalía real (0-1)
- `severidad` (VARCHAR(20)): Nivel de gravedad
  - `baja`: Anomalías menores
  - `media`: Requieren atención
  - `alta`: Requieren acción inmediata
  - `critica`: Emergencia crítica
- `descripcion` (TEXT): Descripción detallada
- `threshold` (FLOAT): Umbral usado para detección
- `is_anomaly` (BOOLEAN): Confirmación de anomalía
- `created_at` (TIMESTAMP): Fecha de creación
- `updated_at` (TIMESTAMP): Última actualización

**Índices:**
- `idx_anomalias_activo_id`: Por activo
- `idx_anomalias_sensor_id`: Por sensor
- `idx_anomalias_timestamp`: Por fecha (descendente)
- `idx_anomalias_is_anomaly`: Solo anomalías confirmadas
- `idx_anomalias_severidad`: Por nivel de severidad
- `idx_anomalias_activo_timestamp`: Combinado activo + fecha

## 🔌 Conexión

**Puerto externo:** 5436
**Usuario:** ia_user
**Contraseña:** ia_pass
**Base de datos:** ia_db

**String de conexión:**
```
postgresql://ia_user:ia_pass@localhost:5436/ia_db
```

**Desde otros contenedores:**
```
postgresql://ia_user:ia_pass@ia-db:5432/ia_db
```

## 🚀 Uso

### Iniciar la base de datos
```bash
docker-compose up -d ia-db
```

### Conectarse a la base de datos
```bash
docker exec -it ia-db psql -U ia_user -d ia_db
```

### Verificar tablas
```sql
\dt
\d+ anomalias
```

### Consultar anomalías
```sql
-- Anomalías recientes
SELECT * FROM anomalias 
WHERE is_anomaly = true 
ORDER BY timestamp DESC 
LIMIT 10;

-- Anomalías críticas
SELECT * FROM anomalias 
WHERE severidad = 'critica' 
  AND is_anomaly = true 
ORDER BY timestamp DESC;

-- Estadísticas por sensor
SELECT 
    sensor_id,
    COUNT(*) as total_anomalias,
    COUNT(*) FILTER (WHERE severidad = 'critica') as criticas,
    AVG(anomaly_score) as avg_score
FROM anomalias 
WHERE is_anomaly = true
GROUP BY sensor_id;
```

## 📝 Notas

- Los scripts SQL en `docker-entrypoint-initdb.d/` se ejecutan automáticamente al crear el contenedor
- Trigger automático actualiza `updated_at` en cada modificación
- Datos de ejemplo incluidos para testing
- Volumen persistente: `ia_data`

## 🔄 Mantenimiento

### Backup
```bash
docker exec ia-db pg_dump -U ia_user ia_db > backup_ia_$(date +%Y%m%d).sql
```

### Restore
```bash
cat backup_ia_20251014.sql | docker exec -i ia-db psql -U ia_user -d ia_db
```

### Ver logs
```bash
docker logs ia-db
```
