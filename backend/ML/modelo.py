# ---------------------------------------------------------------
# Función principal HTM-like
# ---------------------------------------------------------------
def detectar_htm_multivar(df, features, prefix):
    """
    df: dataset de una bomba
    features: columnas relevantes (ya limpias)
    prefix: etiqueta de bomba ("A", "B", "C")
    """
    d = df[['Timestamp'] + features].dropna().sort_values('Timestamp').reset_index(drop=True)
    if len(d) <= W:
        print(f"Muy pocos datos para bomba {prefix}")
        return None

    scaler_X = StandardScaler()
    scaler_y = StandardScaler()
    model = SGDRegressor(max_iter=1, learning_rate='constant', eta0=0.01, warm_start=True)

    scores, probs, preds = [], [], []
    likelihood = 0.0

    # Simulación streaming (como MQTT)
    for i in range(W, len(d)):
        window = d.iloc[i-W:i, 1:].values
        y_true = d.iloc[i, 1:].values

        # Codificación temporal (índices normalizados)
        X = np.arange(W).reshape(-1, 1)
        y_flat = window.flatten()
        X_rep = np.tile(np.arange(W), window.shape[1]).reshape(-1, 1)

        X_scaled = scaler_X.fit_transform(X_rep)
        y_scaled = scaler_y.fit_transform(y_flat.reshape(-1, 1)).ravel()

        # Entrenamiento incremental
        model.partial_fit(X_scaled, y_scaled)

        # Predicción próximo paso
        X_pred = scaler_X.transform(np.array([[W]]))
        y_pred_scaled = model.predict(X_pred)[0]
        y_pred = scaler_y.inverse_transform([[y_pred_scaled]])[0, 0]

        # --- Error y normalización adaptativa ---
        mu = np.mean(window, axis=0)
        sigma = np.std(window, axis=0) + 1e-6
        err = np.sqrt(np.mean(((y_true - mu) / sigma) ** 2))  # error normalizado por feature

        # Escala adaptativa de 0–100 (percentil dinámico)
        if i == W:
            err_ref = err
        else:
            err_ref = (0.99 * err_ref + 0.01 * err)  # adaptación suave

        score = np.clip((err / (err_ref + 1e-6)) * 50, 0, 100)

        # Suavizado tipo HTM (EMA)
        likelihood = (1 - alpha) * likelihood + alpha * score

        preds.append(y_pred)
        scores.append(score)
        probs.append(likelihood)

    res = d.iloc[W:].copy()
    res['Pump'] = prefix
    res['AnomalyScore'] = scores
    res['AnomalyLikelihood'] = probs
    return res