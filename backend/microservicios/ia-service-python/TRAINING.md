# 🤖 Sistema de Entrenamiento Automático

## 📋 Descripción

El microservicio `ia-service-python` incluye un sistema de entrenamiento automático del modelo de Machine Learning que se ejecuta periódicamente.

## ⚙️ Configuración

### Archivo config.yaml

La configuración del intervalo de entrenamiento se maneja en `config/config.yaml`:

```yaml
# Training Configuration
training_enabled: true
training_interval_minutes: 60  # Entrenar modelo cada 60 minutos

# Model Storage
model_dir: "/app/models"
model_filename: "anomaly_model.pkl"
scaler_x_filename: "scaler_X.pkl"
scaler_y_filename: "scaler_Y.pkl"
```

**Para cambiar el intervalo de entrenamiento**, edita el archivo `config/config.yaml` o `config/config.docker.yaml` según tu entorno.

### Docker Compose

El servicio está configurado con un volumen persistente para los modelos:

```yaml
volumes:
  - ia_service_models:/app/models
```

## 🔄 Funcionamiento

### Scheduler Automático

1. Al iniciar el servicio, se carga el modelo existente (si existe)
2. Se inicia un scheduler (APScheduler) que ejecuta el entrenamiento cada 60 minutos
3. Durante el entrenamiento:
   - Se registra en logs el inicio del proceso
   - **TODO**: Se recopilan datos de la base de datos
   - **TODO**: Se preprocesan y transforman los datos
   - **TODO**: Se entrena el modelo con `partial_fit()`
   - **TODO**: Se evalúan métricas de performance
   - Se guarda el modelo en disco
4. El modelo se guarda también al cerrar el servicio

### Logs de Entrenamiento

```
======================================================================
🚀 INICIANDO ENTRENAMIENTO DEL MODELO #1
⏰ Timestamp: 2025-10-26T21:00:00
======================================================================
📊 Fase 1: Recolección de datos
   └─ TODO: Consultar base de datos ia_db
🔧 Fase 2: Preprocesamiento
   └─ TODO: Normalizar y transformar datos
🤖 Fase 3: Entrenamiento incremental
   └─ TODO: Ejecutar model.partial_fit()
📈 Fase 4: Evaluación
   └─ TODO: Calcular métricas de performance
💾 Fase 5: Persistencia
   └─ TODO: Guardar modelo si mejora
======================================================================
✅ ENTRENAMIENTO COMPLETADO EN 0.05s
💾 Modelo guardado: True
======================================================================
```

## 📡 Endpoints

### GET /training/status

Obtiene el estado actual del entrenamiento:

```json
{
  "training_enabled": true,
  "interval_minutes": 60,
  "total_trainings": 5,
  "last_training_time": "2025-10-26T21:00:00",
  "last_training_status": "success",
  "is_training": false,
  "next_training": "2025-10-26T22:00:00",
  "scheduler_running": true
}
```

### POST /training/train

Dispara un entrenamiento manual:

```bash
curl -X POST http://localhost:8095/training/train
```

Respuesta:
```json
{
  "status": "success",
  "training_number": 6,
  "start_time": "2025-10-26T21:30:00",
  "end_time": "2025-10-26T21:30:05",
  "duration_seconds": 5.02,
  "model_saved": true,
  "message": "Entrenamiento completado (simulación)"
}
```

### GET /training/model

Obtiene información sobre el modelo actual:

```json
{
  "model_loaded": true,
  "model_exists_on_disk": true,
  "model_path": "/app/models/anomaly_model.pkl",
  "model_type": "SGDRegressor",
  "model_size_mb": 0.15,
  "last_modified": 1729987200.0
}
```

## 🏗️ Arquitectura

```
┌─────────────────────────────────────────────────────┐
│  FastAPI Application (main.py)                      │
│  - Lifespan events                                  │
│  - Startup: Carga modelo + Inicia scheduler        │
│  - Shutdown: Guarda modelo + Detiene scheduler     │
└────────────────┬────────────────────────────────────┘
                 │
    ┌────────────┴────────────┐
    │                         │
    ▼                         ▼
┌───────────────────┐  ┌──────────────────────┐
│ TrainingService   │  │  ModelManager        │
│ - Scheduler       │  │  - load_model()      │
│ - train_model()   │  │  - save_model()      │
│ - get_stats()     │  │  - get_info()        │
└───────────────────┘  └──────────────────────┘
         │                      │
         │                      ▼
         │            ┌──────────────────────┐
         │            │  Disk Storage        │
         │            │  /app/models/        │
         │            │  - anomaly_model.pkl │
         │            │  - scaler_X.pkl      │
         │            │  - scaler_Y.pkl      │
         │            └──────────────────────┘
         │
         ▼
┌─────────────────────────────────────┐
│  APScheduler                        │
│  - IntervalTrigger (60 minutos)    │
│  - Job: train_model()               │
│  - Max instances: 1                 │
└─────────────────────────────────────┘
```

## 📂 Estructura de Archivos

```
ia-service-python/
├── app/
│   ├── main.py                          # Integración del scheduler
│   ├── config.py                        # Configuración del temporizador
│   ├── models/
│   │   └── model_manager.py            # Gestor del modelo ML
│   └── services/
│       └── training_service.py         # Servicio de entrenamiento
└── Dockerfile
```

## 🚀 Próximos Pasos (TODO)

### 1. Implementar Recolección de Datos

```python
async def collect_training_data(self) -> pd.DataFrame:
    """Recopila datos de anomalías desde la base de datos"""
    # Query a ia_db para obtener histórico de anomalías
    # Incluir: activo_id, sensor_id, timestamp, features, labels
    pass
```

### 2. Implementar Preprocesamiento

```python
def preprocess_data(self, df: pd.DataFrame) -> Tuple[np.ndarray, np.ndarray]:
    """Preprocesa datos para entrenamiento"""
    # Normalización
    # Manejo de valores faltantes
    # Feature engineering
    # Split features/targets
    pass
```

### 3. Implementar Entrenamiento Real

```python
def train_incremental(self, X: np.ndarray, y: np.ndarray):
    """Entrena el modelo con partial_fit"""
    model_manager.model.partial_fit(X, y)
    pass
```

### 4. Implementar Evaluación

```python
def evaluate_model(self, X_test: np.ndarray, y_test: np.ndarray) -> Dict:
    """Evalúa performance del modelo"""
    # Calcular métricas: accuracy, precision, recall, F1
    # Comparar con modelo anterior
    # Decidir si guardar nueva versión
    pass
```

### 5. Agregar Versionado de Modelos

```python
# Guardar con timestamp
model_path = f"/app/models/anomaly_model_{timestamp}.pkl"
# Mantener historial de últimos N modelos
# Permitir rollback a versiones anteriores
```

## 🐛 Debug y Troubleshooting

### Ver logs en tiempo real

```bash
docker logs -f ia-service-python
```

### Verificar estado del scheduler

```bash
curl http://localhost:8095/training/status | jq
```

### Forzar entrenamiento inmediato

```bash
curl -X POST http://localhost:8095/training/train
```

### Verificar archivos del modelo

```bash
docker exec ia-service-python ls -lh /app/models/
```

## 📊 Métricas y Monitoreo

- Total de entrenamientos ejecutados
- Tiempo promedio de entrenamiento
- Última ejecución exitosa
- Estado del scheduler
- Tamaño del modelo en disco
- Performance del modelo (a implementar)

## ⚠️ Consideraciones

1. **Concurrencia**: El scheduler está configurado con `max_instances=1` para evitar entrenamientos simultáneos
2. **Persistencia**: Los modelos se guardan en un volumen Docker que persiste entre reinicios
3. **Rendimiento**: El entrenamiento actualmente es una simulación, el tiempo real dependerá de la cantidad de datos
4. **Memoria**: El modelo SGDRegressor con `warm_start=True` mantiene estado entre entrenamientos
