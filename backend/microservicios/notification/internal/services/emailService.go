package services

import (
	"fmt"
	"log"
	"notification/internal/models"

	"github.com/spf13/viper"
	"gopkg.in/gomail.v2"
)

type EmailService struct {
	smtpHost     string
	smtpPort     int
	smtpUsername string
	smtpPassword string
	fromEmail    string
}

func NewEmailService() *EmailService {
	// Leer EXCLUSIVAMENTE desde variables de entorno (.env)
	smtpHost := viper.GetString("EMAIL_SMTP_HOST")
	smtpPort := viper.GetInt("EMAIL_SMTP_PORT")
	smtpUsername := viper.GetString("EMAIL_SMTP_USERNAME")
	smtpPassword := viper.GetString("EMAIL_SMTP_PASSWORD")
	fromEmail := viper.GetString("EMAIL_FROM_EMAIL")

	// Valores por defecto solo si no están en variables de entorno
	if smtpPort == 0 {
		smtpPort = 587
	}

	// Log de la configuración real que se está usando
	log.Printf("📧 Configuración de correo desde .env:")
	log.Printf("   Host: %s", smtpHost)
	log.Printf("   Port: %d", smtpPort)
	log.Printf("   Username: %s", smtpUsername)
	log.Printf("   From Email: %s", fromEmail)

	if smtpHost == "" || smtpUsername == "" || smtpPassword == "" {
		log.Printf("❌ ADVERTENCIA: Faltan variables de entorno de correo en .env")
		log.Printf("   Verifica que tengas configurado: EMAIL_SMTP_HOST, EMAIL_SMTP_USERNAME, EMAIL_SMTP_PASSWORD, EMAIL_FROM_EMAIL")
	}

	return &EmailService{
		smtpHost:     smtpHost,
		smtpPort:     smtpPort,
		smtpUsername: smtpUsername,
		smtpPassword: smtpPassword,
		fromEmail:    fromEmail,
	}
}

func (e *EmailService) SendNotificationEmail(activo *models.Activo, usuarios []models.Usuario, mensaje string) error {
	if len(usuarios) == 0 {
		log.Println("No hay usuarios para enviar el correo")
		return nil
	}

	// Crear el mensaje de correo
	m := gomail.NewMessage()
	m.SetHeader("From", e.fromEmail)

	// Agregar todos los correos de los usuarios del edificio
	var toEmails []string
	for _, usuario := range usuarios {
		toEmails = append(toEmails, usuario.Correo)
	}
	m.SetHeader("To", toEmails...)

	// Asunto del correo
	subject := fmt.Sprintf("🚨 Alerta de Activo: %s", activo.Nombre)
	m.SetHeader("Subject", subject)

	// Cuerpo del correo en HTML
	body := e.buildEmailHTML(activo, mensaje)
	m.SetBody("text/html", body)

	// Configurar el dialer SMTP
	d := gomail.NewDialer(e.smtpHost, e.smtpPort, e.smtpUsername, e.smtpPassword)

	// Enviar el correo
	if err := d.DialAndSend(m); err != nil {
		log.Printf("Error enviando correo: %v", err)
		return err
	}

	log.Printf("Correo enviado exitosamente a %d usuarios del edificio %s", len(usuarios), activo.Edificio.Direccion)
	return nil
}

func (e *EmailService) buildEmailHTML(activo *models.Activo, mensaje string) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Alerta de Activo</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            margin: 0;
            padding: 20px;
            background-color: #f4f4f4;
        }
        .container {
            max-width: 600px;
            margin: 0 auto;
            background-color: white;
            padding: 30px;
            border-radius: 10px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
        }
        .header {
            text-align: center;
            background-color: #ff4444;
            color: white;
            padding: 20px;
            border-radius: 10px 10px 0 0;
            margin: -30px -30px 20px -30px;
        }
        .alert-icon {
            font-size: 48px;
            margin-bottom: 10px;
        }
        .content {
            margin: 20px 0;
        }
        .info-box {
            background-color: #f8f9fa;
            padding: 15px;
            border-left: 4px solid #007bff;
            margin: 15px 0;
        }
        .footer {
            text-align: center;
            color: #666;
            font-size: 12px;
            margin-top: 30px;
            padding-top: 20px;
            border-top: 1px solid #eee;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <div class="alert-icon">🚨</div>
            <h1>Alerta de Activo</h1>
        </div>
        
        <div class="content">
            <h2>Se ha detectado una alerta en el siguiente activo:</h2>
            
            <div class="info-box">
                <strong>🏢 Edificio:</strong> %s<br>
                <strong>📡 Activo:</strong> %s<br>
                <strong>📝 Mensaje:</strong> %s<br>
                <strong>🕒 Fecha:</strong> %s
            </div>
            
            <p>Por favor, revise el estado del activo y tome las medidas necesarias.</p>
            
            <p><strong>Nota:</strong> Este correo ha sido enviado automáticamente por el sistema InspectAR.</p>
        </div>
        
        <div class="footer">
            <p>Sistema de Monitoreo InspectAR</p>
            <p>Este es un correo automático, por favor no responder.</p>
        </div>
    </div>
</body>
</html>
    `, activo.Edificio.Direccion, activo.Nombre, mensaje, activo.CreatedAt.Format("2006-01-02 15:04:05"))
}
