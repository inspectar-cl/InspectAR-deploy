import os

# ============================================================
# CONFIGURACIÓN GENERAL DEL MODELO HTM-LIKE
# ============================================================

# Tamaño de la ventana temporal (W)
WINDOW_SIZE = int(os.getenv("HTM_WINDOW_SIZE", 180))

# Factor de suavizado (EMA / α)
ALPHA = float(os.getenv("HTM_ALPHA", 0.01))

# Sensibilidad del umbral adaptativo (k)
K_ADAPT = float(os.getenv("HTM_K_ADAPT", 2.0))

# Límite para normalización de scores (por si se escala 0–100)
SCORE_SCALE = float(os.getenv("HTM_SCORE_SCALE", 50.0))

# Mínimo de datos requeridos para calcular rolling threshold
MIN_PERIODS = int(os.getenv("HTM_MIN_PERIODS", 50))

# Longitud de ventana rolling para threshold
ROLLING_WINDOW = int(os.getenv("HTM_ROLLING_WINDOW", 300))

# Modo de operación: "offline" o "stream"
MODE = os.getenv("HTM_MODE", "offline")

# Guardar modelo después de entrenamiento
BASE_DIR = os.path.dirname(os.path.abspath(__file__))
MODEL_DIR = os.path.join(BASE_DIR, "models")

MODEL_FILENAME = "sgd_model.pkl"
SCALER_X_FILENAME = "scaler_X.pkl"
SCALER_Y_FILENAME = "scaler_Y.pkl"

MODEL_PATH = os.path.join(MODEL_DIR, MODEL_FILENAME)
SCALER_X_PATH = os.path.join(MODEL_DIR, SCALER_X_FILENAME)
SCALER_Y_PATH = os.path.join(MODEL_DIR, SCALER_Y_FILENAME)

# ============================================================
# CONFIGURACIÓN DE LOGS Y DEBUG
# ============================================================

LOG_LEVEL = os.getenv("LOG_LEVEL", "INFO")
SAVE_RESULTS = os.getenv("SAVE_RESULTS", "false").lower() == "true"
OUTPUT_PATH = os.getenv("OUTPUT_PATH", "/app/results")  # Carpeta donde guardar logs/resultados

# ============================================================
# CONFIGURACIÓN DE API
# ============================================================

API_TITLE = "Microservicio ML - Detección de Anomalías HTM-like"
API_VERSION = "1.0.0"
API_PORT = int(os.getenv("API_PORT", 8086))

# ============================================================
# CONFIGURACIÓN DE SEED Y REPRODUCIBILIDAD
# ============================================================

RANDOM_SEED = int(os.getenv("RANDOM_SEED", 42))

# ============================================================
# LOG DE CARGA DE CONFIGURACIÓN
# ============================================================

print(f"[CONFIG] Modelo HTM-like cargado en modo {MODE.upper()} con W={WINDOW_SIZE}, α={ALPHA}, k={K_ADAPT}")
