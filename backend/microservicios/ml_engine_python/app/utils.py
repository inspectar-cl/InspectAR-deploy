# app/utils.py
import numpy as np
import pandas as pd
from .config import WINDOW_SIZE as W, K_ADAPT as k_adapt

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


def clean_timestamps(df: pd.DataFrame, max_invalid_percent: float = 5.0) -> tuple[pd.DataFrame, dict]:
    """
    Limpia y maneja timestamps inválidos en el DataFrame.
    
    Estrategia adaptativa basada en porcentaje de timestamps inválidos:
    - Si < max_invalid_percent: Elimina registros inválidos
    - Si >= max_invalid_percent pero < 20%: Intenta interpolar
    - Si >= 20%: Lanza excepción (calidad de datos muy baja)
    
    Args:
        df: DataFrame con columna 'timestamp' (ya debe estar como string)
        max_invalid_percent: Porcentaje máximo permitido para eliminar registros (default: 5%)
    
    Returns:
        tuple: (DataFrame limpio, diccionario con estadísticas de limpieza)
    
    Raises:
        ValueError: Si la calidad de datos es muy baja (>20% inválidos)
    """
    import logging
    logger = logging.getLogger("ml_service")
    
    stats = {
        "total_records": len(df),
        "invalid_timestamps": 0,
        "dropped_records": 0,
        "interpolated_records": 0,
        "strategy": None,
        "invalid_percent": 0.0
    }
    
    if len(df) == 0:
        return df, stats
    
    # Convertir timestamp a datetime
    df = df.copy()
    df['timestamp'] = pd.to_datetime(df['timestamp'], errors='coerce')
    
    # Contar timestamps inválidos
    invalid_mask = df['timestamp'].isna()
    stats["invalid_timestamps"] = int(invalid_mask.sum())
    stats["invalid_percent"] = (stats["invalid_timestamps"] / stats["total_records"]) * 100
    
    if stats["invalid_timestamps"] == 0:
        stats["strategy"] = "no_cleaning_needed"
        return df, stats
    
    logger.warning(
        f"Timestamps inválidos: {stats['invalid_timestamps']} de {stats['total_records']} "
        f"({stats['invalid_percent']:.2f}%)"
    )
    
    # Estrategia 1: Calidad muy baja - rechazar
    if stats["invalid_percent"] >= 20.0:
        stats["strategy"] = "rejected"
        raise ValueError(
            f"Calidad de datos insuficiente: {stats['invalid_percent']:.1f}% de timestamps inválidos. "
            f"Se requiere al menos 80% de timestamps válidos para análisis confiable."
        )
    
    # Estrategia 2: Eliminar registros (pocos inválidos)
    if stats["invalid_percent"] < max_invalid_percent:
        df_clean = df[~invalid_mask].copy().reset_index(drop=True)
        stats["dropped_records"] = stats["invalid_timestamps"]
        stats["strategy"] = "drop_invalid"
        logger.info(
            f"✓ Estrategia: Eliminar registros inválidos "
            f"({stats['dropped_records']} registros, {stats['invalid_percent']:.2f}%)"
        )
        return df_clean, stats
    
    # Estrategia 3: Interpolar (porcentaje moderado de inválidos)
    logger.info(
        f"Intentando interpolar timestamps "
        f"({stats['invalid_percent']:.2f}% inválidos > {max_invalid_percent}%)"
    )
    
    # Ordenar por índice original para mantener secuencia
    df = df.sort_index()
    
    # Verificar si hay suficientes timestamps válidos para interpolar
    valid_timestamps = df.loc[~invalid_mask, 'timestamp']
    
    if len(valid_timestamps) < 2:
        stats["strategy"] = "insufficient_valid_data"
        raise ValueError(
            f"No hay suficientes timestamps válidos para interpolar "
            f"(solo {len(valid_timestamps)} válidos de {len(df)})"
        )
    
    # Calcular frecuencia promedio entre timestamps válidos consecutivos
    time_diffs = valid_timestamps.diff().dropna()
    
    if len(time_diffs) > 0:
        median_freq = time_diffs.median()
        mean_freq = time_diffs.mean()
        
        logger.debug(
            f"Frecuencia detectada - Media: {mean_freq}, Mediana: {median_freq}"
        )
        
        # Usar forward fill y backward fill para rellenar gaps pequeños
        df['timestamp'] = df['timestamp'].fillna(method='ffill').fillna(method='bfill')
        
        # Verificar si aún quedan NaN (gaps muy grandes al inicio/fin)
        remaining_nan = df['timestamp'].isna().sum()
        
        if remaining_nan > 0:
            # Si quedan NaN, eliminarlos
            df_clean = df[~df['timestamp'].isna()].copy().reset_index(drop=True)
            stats["interpolated_records"] = stats["invalid_timestamps"] - remaining_nan
            stats["dropped_records"] = remaining_nan
            stats["strategy"] = "interpolate_and_drop"
            logger.info(
                f"✓ Estrategia: Interpolar + eliminar extremos "
                f"(interpolados: {stats['interpolated_records']}, "
                f"eliminados: {stats['dropped_records']})"
            )
        else:
            df_clean = df.reset_index(drop=True)
            stats["interpolated_records"] = stats["invalid_timestamps"]
            stats["strategy"] = "interpolate"
            logger.info(
                f"✓ Estrategia: Interpolar todos "
                f"({stats['interpolated_records']} timestamps interpolados)"
            )
        
        return df_clean, stats
    
    # Fallback: si no se puede interpolar, eliminar
    df_clean = df[~invalid_mask].copy().reset_index(drop=True)
    stats["dropped_records"] = stats["invalid_timestamps"]
    stats["strategy"] = "drop_fallback"
    logger.warning(
        f"⚠ No se pudo interpolar. Eliminando registros inválidos "
        f"({stats['dropped_records']} registros)"
    )
    
    return df_clean, stats
