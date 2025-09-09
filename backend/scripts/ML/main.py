from multiprocessing import Process, Queue
from caldera import run_simulador
from gaussean import run_model
import threading

def escuchar_input(control_queue):
    print("\n💡 Escribe 1 (leve), 2 (extremo), o cualquier otra cosa (normal)")
    while True:
        entrada = input()
        control_queue.put(entrada)

if __name__ == "__main__":
    data_queue = Queue()
    control_queue = Queue()

    # Lanza el hilo que escucha el input del usuario (en el proceso padre)
    input_thread = threading.Thread(target=escuchar_input, args=(control_queue,), daemon=True)
    input_thread.start()

    # Lanza procesos hijos
    simulador_proc = Process(target=run_simulador, args=(data_queue, control_queue))
    modelo_proc = Process(target=run_model, args=(data_queue,))

    simulador_proc.start()
    modelo_proc.start()

    simulador_proc.join()
    modelo_proc.join()
