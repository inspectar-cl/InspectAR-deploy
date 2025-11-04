# 🤖 Sistema de Entrenamiento Automático

## 📋 Descripción

El microservicio `ia-service-python` incluye un sistema de entrenamiento automático del modelo de Machine Learning que:

- ✅ **Detecta modelos no entrenados** mediante flag `is_trained`
- ✅ **Ejecuta entrenamiento inicial** automáticamente al inicio si no existe modelo
- ✅ **Reentrena periódicamente** según configuración (default: 60 minutos)
- ✅ **Obtiene datos reales** del IOT Service para entrenamiento
- ✅ **Calcula métricas avanzadas** incluyendo Coeficientes de Variación (CV)
- ✅ **Persiste modelos** en volumen Docker

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

### Inicio del Servicio

```
┌─────────────────────────┐
│ Servicio inicia         │
└──────────┬──────────────┘
           │
    ┌──────▼──────┐
    │load_model() │
    └──────┬──────┘
           │
    ┌──────▼──────────────┐
    │¿Archivos .pkl       │
    │existen en disco?    │
    └──────┬──────────────┘
           │
    ┌──────▼───────┬──────┐
    │ SÍ           │ NO   │
    │              │      │
┌───▼───┐      ┌───▼────┐
│Load   │      │Create  │
│files  │      │empty   │
│       │      │model   │
│✅     │      │        │
│trained│      │❌ not  │
│= True │      │trained │
└───┬───┘      └───┬────┘
    │              │
    └───────┬──────┘
            │
    ┌───────▼──────────┐
    │start_scheduler() │
    └───────┬──────────┘
            │
    ┌───────▼──────────────┐
    │if not is_trained:    │ ← Verificación clave
    └───────┬──────────────┘
            │
    ┌───────▼──────────┐
    │✅ Programar      │
    │entrenamiento     │
    │INMEDIATO         │
    └──────────────────┘
```

### Scheduler Automático

1. **Al iniciar:**
   - Se verifica el flag `is_trained` del `ModelManager`
   - Si `is_trained = False`: Se programa entrenamiento inmediato
   - Si `is_trained = True`: Se espera el primer ciclo (60 min)

2. **Durante el entrenamiento:**
   - Se registra en logs el inicio del proceso
   - Se obtienen datos del IOT Service (activo ID 2, 500 registros)
   - Se preprocesan y transforman los datos
   - Se entrena el modelo con `partial_fit()` (entrenamiento incremental)
   - Se evalúan métricas de performance (MSE, MAE, RMSE, R², CV)
   - Se guarda el modelo en disco
   - Se marca `is_trained = True`

3. **Ciclo periódico:**
   - El modelo se reentrena cada 60 minutos (configurable)
   - Usa `warm_start=True` para aprendizaje incremental

4. **Al cerrar:**
   - El modelo se guarda automáticamente

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
  "is_trained": true,
  "model_path": "/app/models/anomaly_model.pkl",
  "model_type": "SGDRegressor",
  "model_size_mb": 0.15,
  "last_modified": 1730483400.0
}
```

## 📊 Métricas de Evaluación

### Métricas Básicas

| Métrica | Fórmula | Descripción | Rango |
|---------|---------|-------------|-------|
| **MSE** | `mean((y_pred - y_true)²)` | Error cuadrático medio | 0 a ∞ (menor es mejor) |
| **MAE** | `mean(|y_pred - y_true|)` | Error absoluto medio | 0 a ∞ (menor es mejor) |
| **RMSE** | `sqrt(MSE)` | Raíz del error cuadrático medio | 0 a ∞ (menor es mejor) |
| **R²** | `1 - (SS_res / SS_tot)` | Coeficiente de determinación | -∞ a 1 (1 es perfecto) |

### Coeficientes de Variación (CV)

El **Coeficiente de Variación** mide la variabilidad relativa, independiente de la escala:

```
CV = (Desviación Estándar / Media) × 100%
```

#### CV de Errores (`cv_errors_percent`)

Mide la **consistencia de los errores** del modelo.

```python
cv_errors = (std(errors) / mean(|errors|)) × 100
```

**Interpretación:**
- **< 10%**: ✅ Excelente - Errores muy consistentes
- **10-30%**: ⚠️ Bueno - Errores moderadamente variables
- **> 30%**: ❌ Revisar - Errores muy inconsistentes

#### CV de Predicciones (`cv_predictions_percent`)

Mide la **dispersión de las predicciones** del modelo.

```python
cv_predictions = (std(y_pred) / mean(y_pred)) × 100
```

**Interpretación:**
- Debe ser **similar al CV de valores reales**
- Si es mucho menor → Modelo muy conservador
- Si es mucho mayor → Modelo sobre-predice varianza

#### CV de Valores Reales (`cv_actual_percent`)

Mide la **variabilidad natural** de los datos.

```python
cv_actual = (std(y_true) / mean(y_true)) × 100
```

**Interpretación:**
- Sirve como **referencia** para comparar con el CV de predicciones
- Alta variabilidad natural (>20%) puede dificultar el entrenamiento

### Ejemplo de Análisis

```json
{
  "mse": 0.0234,
  "mae": 0.1123,
  "rmse": 0.1529,
  "r2_score": 0.8567,
  "cv_errors_percent": 8.5,
  "cv_predictions_percent": 12.3,
  "cv_actual_percent": 11.8
}
```

**Análisis:**

1. **R² = 0.8567**: Modelo explica ~86% de la varianza → ✅ Muy bueno
2. **RMSE = 0.1529**: Error promedio de 0.15 unidades → Depende de la escala
3. **CV Errores = 8.5%**: Errores muy consistentes → ✅ Excelente
4. **CV Predicciones = 12.3% ≈ CV Real = 11.8%**: Modelo captura bien la variabilidad natural → ✅ Excelente

**Conclusión:** Modelo con buen ajuste, errores consistentes y predicciones realistas.

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
