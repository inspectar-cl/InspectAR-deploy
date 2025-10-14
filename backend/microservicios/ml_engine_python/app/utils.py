# app/utils.py
import numpy as np
import pandas as pd

W = 180 # Tamaño de ventana
alpha = 0.01 # Factor de suavizado
k_adapt = 2.0  # sensibilidad del umbral

def rolling_threshold(series: pd.Series, window: int = W, k: float = k_adapt, min_periods: int = 50) -> pd.Series:
    """
    Umbral adaptativo: rolling mean + k * rolling std.
    Devuelve serie con NaN en los primeros valores donde no alcanza min_periods.
    """
    mean_rolling = series.rolling(window=window, min_periods=min_periods).mean()
    std_rolling = series.rolling(window=window, min_periods=min_periods).std()
    thr = (mean_rolling + k * std_rolling).clip(lower=0)
    return thr.fillna(np.nan)

def severity_from_likelihood(likelihood: float) -> str:
    if likelihood < 40:
        return "Baja"
    elif likelihood < 70:
        return "Media"
    else:
        return "Alta"

def description_from_severity(sev: str) -> str:
    d = {
        "Baja": "Funcionamiento dentro del rango esperado.",
        "Media": "Comportamiento irregular detectado. Revisar condiciones operativas.",
        "Alta": "Anomalía crítica detectada. Atención prioritaria requerida."
    }
    return d.get(sev, "")

def normalize_score_0_100(series: pd.Series) -> pd.Series:
    """
    Asegura que la serie quede dentro de 0..100 y normaliza si es necesario.
    Asume que input ya es escala cercana a 0-100, pero recorta y redondea.
    """
    s = series.clip(lower=0, upper=100)
    return s.round(2)

def compute_is_anomaly(likelihood: pd.Series, threshold: pd.Series) -> pd.Series:
    """
    Devuelve 1/0 comparando likelihood > threshold.
    Si threshold es NaN (ventana inicial), devuelve 0.
    """
    is_anom = (likelihood > threshold).astype(int)
    is_anom = is_anom.fillna(0).astype(int)
    return is_anom
