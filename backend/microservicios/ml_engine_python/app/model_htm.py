import numpy as np
import pandas as pd
from sklearn.linear_model import SGDRegressor
from sklearn.preprocessing import StandardScaler
from .utils import rolling_threshold, severity_from_likelihood, description_from_severity, normalize_score_0_100, compute_is_anomaly, preprocess_timeseries
from .config import WINDOW_SIZE as W, ALPHA as alpha, K_ADAPT as k_adapt

def detectar_htm_multivar(df, activo_id):
    """
    Recibe un DataFrame con columnas de sensores y devuelve
    los resultados del análisis de anomalías.
    """
    if 'Timestamp' in df.columns:
        df = df.sort_values('Timestamp').reset_index(drop=True)
    d = preprocess_timeseries(df, timestamp_col="Timestamp", method="linear")

    if len(d) <= W:
        return ValueError(f"Datos insuficientes ({len(d)}) para bomba {activo_id}")

    scaler_X = StandardScaler()
    scaler_y = StandardScaler()
    model = SGDRegressor(max_iter=1, learning_rate='constant', eta0=0.01, warm_start=True)

    scores, probs, preds = [], [], []
    likelihood = 0.0
    err_ref = None

    for i in range(W, len(d)):
        window = d.iloc[i-W:i, 1:].values
        y_true = d.iloc[i, 1:].values

        # Codificación temporal (índices normalizados)
        y_flat = window.flatten()
        X_rep = np.tile(np.arange(W), window.shape[1]).reshape(-1, 1)

        X_scaled = scaler_X.fit_transform(X_rep)
        y_scaled = scaler_y.fit_transform(y_flat.reshape(-1, 1)).ravel()

        # Entrenamiento incremental
        model.partial_fit(X_scaled, y_scaled)

        """# Predicción próximo paso
        X_pred = scaler_X.transform(np.array([[W]]))
        y_pred_scaled = model.predict(X_pred)[0]
        y_pred = scaler_y.inverse_transform([[y_pred_scaled]])[0, 0]"""

        # --- Error normalizado ---
        mu = np.mean(window, axis=0)
        sigma = np.std(window, axis=0) + 1e-6
        err = np.sqrt(np.mean(((y_true - mu) / sigma) ** 2))

        # Escala adaptativa 0–100
        if err_ref is None:
            err_ref = err
        else:
            err_ref = (1 - alpha) * err_ref + alpha * err
        score = np.clip((err / (err_ref + 1e-6)) * 50, 0, 100)

        # Suavizado HTM-like
        likelihood = (1 - alpha) * likelihood + alpha * score

        # preds.append(y_pred)
        scores.append(score)
        probs.append(likelihood)

    # construir DataFrame de resultados alineados con Timestamps desde W en adelante
    res = d.iloc[W:].copy().reset_index(drop=True)
    res['AnomalyScore'] = normalize_score_0_100(pd.Series(scores))
    res['AnomalyLikelihood'] = normalize_score_0_100(pd.Series(probs))

    # umbral adaptativo rolling (sobre AnomalyLikelihood)
    thr = rolling_threshold(res['AnomalyLikelihood'], window=W, k=k_adapt, min_periods=50)
    res['Threshold'] = thr.ffill().fillna(1e6)  # si no hay threshold inicial, usar gran valor para evitar marcar anomaly temprana

    # decisión y texto
    res['is_anomaly'] = compute_is_anomaly(res['AnomalyLikelihood'], res['Threshold'])
    res['Severity'] = res['AnomalyLikelihood'].apply(lambda x: severity_from_likelihood(x))
    res['Description'] = res['Severity'].apply(description_from_severity)

    # preparar salida: lista de dicts + summary (último)
    resultados = res[['Timestamp', 'AnomalyScore', 'AnomalyLikelihood',
                       'Threshold', 'is_anomaly', 'Severity', 'Description']].to_dict(orient='records')

    return {
        "activo_id": activo_id,
        "n_rows": len(res),
        "n_anomalies": int(res['is_anomaly'].sum()),
        "results": resultados,
        "last": resultados[-1] if len(resultados) > 0 else None
    }
