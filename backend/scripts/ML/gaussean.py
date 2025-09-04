import pickle
import os
from river import naive_bayes
from multiprocessing import Queue
import time

from escribirJson import escribir_alerta
from cambiarEstado import put_estado
from enviarCorreo import enviar_correo

activo_id = "BombaDeAgua1"

def run_model(queue: Queue):
    modelo_path = "modelo_gaussian_nb.pkl"

    if os.path.exists(modelo_path):
        with open(modelo_path, "rb") as f:
            modelo_sup = pickle.load(f)
        print("✅ Modelo cargado")
    else:
        modelo_sup = naive_bayes.GaussianNB()
        print("🆕 Modelo creado")

    try:
        while True:
            if not queue.empty():
                data = queue.get()
                x = {
                    "caudal": data["caudal"],
                    "presion": data["presion"],
                    "temperatura": data["temperatura"]
                }
                etiqueta = data.get("etiqueta", None)

                if etiqueta is None:
                    # No hay etiqueta → solo predecir
                    pred = modelo_sup.predict_one(x)
                    probas = modelo_sup.predict_proba_one(x)
                    print(f"🤖 Predicción sin etiqueta: {'Anómalo' if pred == 1 else 'Normal'} - Probabilidad: {probas.get(1, 0.0):.2f} | Datos: {x}")
                    enviar_correo(activo_id)
                else:
                    # Aprendizaje supervisado
                    y = 1 if etiqueta == "Anormal" else 0
                    modelo_sup.learn_one(x, y)

                    pred = modelo_sup.predict_one(x)
                    probas = modelo_sup.predict_proba_one(x)


                    anom_prob = probas.get(1, 0.0)
                    anom_score = round(anom_prob * 100, 2)

                    # Determinar el estado
                    if anom_score < 30:
                        estado = "OK"
                    elif anom_score < 70:
                        estado = "Medio"
                    else:
                        estado = "Crítico"
                        enviar_correo(activo_id)

                    put_estado(estado)

                    print(f"🔥 Predicción: {'Anómalo' if pred == 1 else 'Normal'} | Score: {anom_score} | Estado: {estado} | Datos: {x}")


                    print(f"📚 Aprendido: {etiqueta} → Predicción: {'Anómalo' if pred == 1 else 'Normal'} - Prob: {probas.get(1, 0.0):.2f} | Datos: {x}")

                    if  x["caudal"] >= 200:
                        escribir_alerta(activo_id + "1", "Caudal Alto")
                    if x["presion"] >= 13 * 14.5038:
                        escribir_alerta(activo_id + "2", "Presion Alta")
                    if x["temperatura"] >= 170:
                        escribir_alerta(activo_id, "Temperatura Alta")

                # Guardar el modelo
                with open(modelo_path, "wb") as f:
                    pickle.dump(modelo_sup, f)

            else:
                time.sleep(0.1)

    except KeyboardInterrupt:
        print("\n🛑 Listener detenido.")
