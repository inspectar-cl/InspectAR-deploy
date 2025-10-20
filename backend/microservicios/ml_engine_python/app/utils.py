# app/utils.py
import numpy as np
import pandas as pd
from typing import Dict, List, Any
from collections import defaultdict
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
        df: DataFrame con columna 'Timestamp' (ya debe estar como string)
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
    
    # Convertir timestamp a datetime - intentar múltiples formatos
    df = df.copy()
    
    # Intentar parsear con formato ISO8601/RFC3339 primero
    df['Timestamp'] = pd.to_datetime(df['Timestamp'], format='ISO8601', errors='coerce')

    # Si aún hay NaT, intentar con inferencia automática
    invalid_mask = df['Timestamp'].isna()
    if invalid_mask.any():
        still_invalid = df.loc[invalid_mask, 'Timestamp']
        parsed = pd.to_datetime(still_invalid, errors='coerce', utc=True)
        # Convertir explícitamente a datetime64[ns] para evitar warnings de dtype
        df.loc[invalid_mask, 'Timestamp'] = parsed.astype('datetime64[ns]')
    
    # Contar timestamps inválidos finales
    invalid_mask = df['Timestamp'].isna()
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
    valid_timestamps = df.loc[~invalid_mask, 'Timestamp']
    
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
        df['Timestamp'] = df['Timestamp'].fillna(method='ffill').fillna(method='bfill')
        
        # Verificar si aún quedan NaN (gaps muy grandes al inicio/fin)
        remaining_nan = df['Timestamp'].isna().sum()
        
        if remaining_nan > 0:
            # Si quedan NaN, eliminarlos
            df_clean = df[~df['Timestamp'].isna()].copy().reset_index(drop=True)
            stats["interpolated_records"] = stats["invalid_timestamps"] - remaining_nan
            stats["dropped_records"] = remaining_nan
            stats["strategy"] = "interpolate_and_drop"
            logger.info(
                f" Estrategia: Interpolar + eliminar extremos "
                f"(interpolados: {stats['interpolated_records']}, "
                f"eliminados: {stats['dropped_records']})"
            )
        else:
            df_clean = df.reset_index(drop=True)
            stats["interpolated_records"] = stats["invalid_timestamps"]
            stats["strategy"] = "interpolate"
            logger.info(
                f" Estrategia: Interpolar todos "
                f"({stats['interpolated_records']} timestamps interpolados)"
            )
        
        return df_clean, stats
    
    # Fallback: si no se puede interpolar, eliminar
    df_clean = df[~invalid_mask].copy().reset_index(drop=True)
    stats["dropped_records"] = stats["invalid_timestamps"]
    stats["strategy"] = "drop_fallback"
    logger.warning(
        f"No se pudo interpolar. Eliminando registros inválidos "
        f"({stats['dropped_records']} registros)"
    )
    
    return df_clean, stats


def transform_sensores_to_records(data: Dict[str, Any]) -> Dict[str, Any]:
    """
    Transforma el JSON de sensores agrupados a formato plano por timestamp,
    y agrega la clave 'pump' con el valor de 'pump'.
    
    Args:
        data: Diccionario con estructura:
            {
                "pump": 2,
                "sensores": [
                    {"sensor_id": "...", "datos": [{"tiempo": "...", "valor": ...}, ...]},
                    ...
                ]
            }
    
    Returns:
        Diccionario con formato:
            {
                "pump": pump,
                "records": [
                    {"timestamp": "...", "sensor1": val, ...},
                    ...
                ]
            }
    """
    sensores = data.get("sensores", [])
    pump = data.get("pump", "UNKNOWN")
    
    # Diccionario temporal agrupado por timestamp
    records_by_time: Dict[str, Dict[str, Any]] = defaultdict(dict)
    
    for sensor in sensores:
        sensor_id = sensor.get("sensor_id")
        for dato in sensor.get("datos", []):
            timestamp = dato.get("tiempo")
            valor = dato.get("valor")
            if timestamp is not None and valor is not None:
                records_by_time[timestamp][sensor_id] = valor
    
    # Convertir a lista de dicts con timestamp
    records = []
    for timestamp, features in records_by_time.items():
        record = {"Timestamp": timestamp}
        record.update(features)
        records.append(record)
    
    # Ordenar por timestamp
    records.sort(key=lambda x: x["Timestamp"])
    
    return {
        "pump": pump,
        "records": records
    }

def preprocess_timeseries(df: pd.DataFrame, timestamp_col: str = "Timestamp", method: str = "linear") -> pd.DataFrame:
    """
    Preprocesa un DataFrame de sensores para análisis de series de tiempo:
      - Ordena por timestamp
      - Interpola valores faltantes (NaN)
      - Opcionalmente rellena los extremos con forward/backward fill si es necesario
    
    Args:
        df: DataFrame con columna de timestamp y columnas de sensores.
        timestamp_col: Nombre de la columna de tiempo.
        method: Método de interpolación ('linear', 'time', 'polynomial', etc.)
        
    Returns:
        DataFrame limpio, ordenado y sin NaN intermedios.
    """
    # Asegurarse de que la columna timestamp sea datetime
    df = df.copy()
    df[timestamp_col] = pd.to_datetime(df[timestamp_col], utc=True)
    
    # Ordenar por tiempo
    df = df.sort_values(timestamp_col).reset_index(drop=True)
    
    # Interpolación de NaNs
    df_interpolated = df.interpolate(method=method, limit_direction='both', axis=0)
    
    return df_interpolated
