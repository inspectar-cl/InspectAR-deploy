from email.mime.text import MIMEText
import smtplib

def enviar_correo(activo):
#Configuracion del correo
    remitente = "inspectar2000@gmail.com"
    contraseña = "kdgz xadt ebwt dsgv"
    destinatario = "diego.morella20@gmail.com"

    #Crear cuerpo del mensaje
    cuerpo = f"Se ha detectado una anomalia critica en los sensores!.\n Activo afectado: {activo}\n"

    mensaje = MIMEText(cuerpo)
    mensaje["Subject"] = "🚨 Alerta de Anomalía Detectada - InspectAR"
    mensaje["From"] = remitente
    mensaje["To"] = destinatario

    #Enviar correo
    try:
        with smtplib.SMTP("smtp.gmail.com", 587) as server:
            server.starttls()
            server.login(remitente, contraseña)
            server.send_message(mensaje)
            print("✅ Correo enviado con éxito a", destinatario)
    except Exception as e:
        print("❌ Error al enviar correo:", e)