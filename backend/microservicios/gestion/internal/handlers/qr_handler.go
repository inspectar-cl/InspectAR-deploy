package handlers

import (
	"encoding/base64"
	"log"
	"net/http"
	"strconv"

	"gestion/internal/repository"
	"gestion/internal/services"

	"github.com/gin-gonic/gin"
)

type QRHandler struct {
	activoRepo *repository.ActivoRepository
	qrService  *services.QRService
}

func NewQRHandler(activoRepo *repository.ActivoRepository, qrService *services.QRService) *QRHandler {
	return &QRHandler{
		activoRepo: activoRepo,
		qrService:  qrService,
	}
}

// ObtenerInfoPorCodigo - GET /qr/:codigo
// Retorna información del activo por su código (descripción y tipo)
func (h *QRHandler) ObtenerInfoPorCodigo(c *gin.Context) {
	codigo := c.Param("codigo")

	if codigo == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Código requerido"})
		return
	}

	// Buscar activo por código
	activo, err := h.activoRepo.GetByCodigo(codigo)
	if err != nil {
		log.Printf("Error buscando activo por código %s: %v", codigo, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Activo no encontrado"})
		return
	}

	// Retornar información completa del activo incluyendo ID
	response := gin.H{
		"id":          activo.ID,
		"codigo":      activo.CodigoActivo,
		"tipo":        activo.Tipo,
		"descripcion": activo.Descripcion,
		"nombre":      activo.Nombre,
		"ubicacion":   activo.Ubicacion,
		"edificio_id": activo.EdificioID,
	}

	c.JSON(http.StatusOK, response)
}

// GenerarYObtenerQR - GET /qr/obtener/:activo_id
// Genera un código QR único para el activo (si no existe) y lo retorna como imagen PNG
func (h *QRHandler) GenerarYObtenerQR(c *gin.Context) {
	idStr := c.Param("activo_id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de activo inválido"})
		return
	}

	// Obtener activo
	activo, err := h.activoRepo.GetByID(id)
	if err != nil {
		log.Printf("Error obteniendo activo ID %d: %v", id, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Activo no encontrado"})
		return
	}

	// Verificar que el activo tenga código (debería generarse automáticamente)
	if activo.CodigoActivo == nil || *activo.CodigoActivo == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "El activo no tiene código generado"})
		return
	}

	var qrBase64 string

	// Si no tiene QR generado, generarlo ahora
	if activo.CodigoQR == nil || *activo.CodigoQR == "" {
		log.Printf("Generando QR para activo ID %d, código: %s", id, *activo.CodigoActivo)

		qrBase64Generado, urlQR, err := h.qrService.GenerarQRActivo(activo.ID, *activo.CodigoActivo)
		if err != nil {
			log.Printf("Error generando QR: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generando código QR"})
			return
		}

		// Guardar QR en base de datos
		if err := h.activoRepo.ActualizarQR(activo.ID, qrBase64Generado, urlQR); err != nil {
			log.Printf("⚠️ Error guardando QR en BD: %v", err)
			// No falla, solo registra el error y devuelve el QR generado
		}

		qrBase64 = qrBase64Generado
	} else {
		qrBase64 = *activo.CodigoQR
	}

	// Decodificar Base64 a bytes
	qrBytes, err := base64.StdEncoding.DecodeString(qrBase64)
	if err != nil {
		log.Printf("Error decodificando QR: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error decodificando código QR"})
		return
	}

	// Retornar imagen PNG
	c.Header("Content-Type", "image/png")
	c.Header("Content-Disposition", `inline; filename="`+*activo.CodigoActivo+`.png"`)
	c.Header("Cache-Control", "public, max-age=86400") // Cache por 24 horas
	c.Data(http.StatusOK, "image/png", qrBytes)
}

// VerQR - GET /qr/ver/:activo_id
// Muestra el QR en una página HTML para visualización en navegador
func (h *QRHandler) VerQR(c *gin.Context) {
	idStr := c.Param("activo_id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de activo inválido"})
		return
	}

	// Obtener activo
	activo, err := h.activoRepo.GetByID(id)
	if err != nil {
		log.Printf("Error obteniendo activo ID %d: %v", id, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Activo no encontrado"})
		return
	}

	// Verificar que el activo tenga código
	if activo.CodigoActivo == nil || *activo.CodigoActivo == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "El activo no tiene código generado"})
		return
	}

	var qrBase64 string

	// Si no tiene QR generado, generarlo ahora
	if activo.CodigoQR == nil || *activo.CodigoQR == "" {
		log.Printf("Generando QR para visualización - activo ID %d", id)

		qrBase64Generado, urlQR, err := h.qrService.GenerarQRActivo(activo.ID, *activo.CodigoActivo)
		if err != nil {
			log.Printf("Error generando QR: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generando código QR"})
			return
		}

		// Guardar QR en base de datos
		if err := h.activoRepo.ActualizarQR(activo.ID, qrBase64Generado, urlQR); err != nil {
			log.Printf("⚠️ Error guardando QR en BD: %v", err)
		}

		qrBase64 = qrBase64Generado
	} else {
		qrBase64 = *activo.CodigoQR
	}

	// Renderizar página HTML con el QR
	htmlContent := `<!DOCTYPE html>
<html lang="es">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Código QR - ` + activo.Nombre + `</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
            display: flex;
            justify-content: center;
            align-items: center;
            padding: 20px;
        }
        .container {
            background: white;
            border-radius: 20px;
            box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
            padding: 40px;
            max-width: 500px;
            width: 100%;
            text-align: center;
        }
        .header {
            margin-bottom: 30px;
        }
        .header h1 {
            color: #2d3748;
            font-size: 28px;
            margin-bottom: 10px;
        }
        .codigo {
            background: #f7fafc;
            padding: 10px 20px;
            border-radius: 10px;
            display: inline-block;
            font-family: 'Courier New', monospace;
            font-size: 18px;
            color: #667eea;
            font-weight: bold;
            margin-bottom: 20px;
        }
        .info {
            background: #edf2f7;
            padding: 20px;
            border-radius: 10px;
            margin-bottom: 30px;
            text-align: left;
        }
        .info-row {
            display: flex;
            justify-content: space-between;
            padding: 8px 0;
            border-bottom: 1px solid #cbd5e0;
        }
        .info-row:last-child {
            border-bottom: none;
        }
        .info-label {
            font-weight: 600;
            color: #4a5568;
        }
        .info-value {
            color: #2d3748;
        }
        .qr-container {
            background: #f7fafc;
            padding: 20px;
            border-radius: 15px;
            margin-bottom: 20px;
        }
        .qr-image {
            max-width: 100%;
            height: auto;
            border-radius: 10px;
            border: 2px solid #667eea;
        }
        .buttons {
            display: flex;
            gap: 10px;
            justify-content: center;
            flex-wrap: wrap;
        }
        .btn {
            padding: 12px 24px;
            border: none;
            border-radius: 8px;
            font-size: 14px;
            font-weight: 600;
            cursor: pointer;
            transition: all 0.3s ease;
            text-decoration: none;
            display: inline-block;
        }
        .btn-primary {
            background: #667eea;
            color: white;
        }
        .btn-primary:hover {
            background: #5568d3;
            transform: translateY(-2px);
            box-shadow: 0 4px 12px rgba(102, 126, 234, 0.4);
        }
        .btn-secondary {
            background: #48bb78;
            color: white;
        }
        .btn-secondary:hover {
            background: #38a169;
            transform: translateY(-2px);
            box-shadow: 0 4px 12px rgba(72, 187, 120, 0.4);
        }
        .footer {
            margin-top: 20px;
            color: #718096;
            font-size: 12px;
        }
        @media print {
            body {
                background: white;
            }
            .container {
                box-shadow: none;
            }
            .buttons {
                display: none;
            }
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>📱 Código QR del Activo</h1>
            <div class="codigo">` + *activo.CodigoActivo + `</div>
        </div>
        
        <div class="info">
            <div class="info-row">
                <span class="info-label">Nombre:</span>
                <span class="info-value">` + activo.Nombre + `</span>
            </div>
            <div class="info-row">
                <span class="info-label">Tipo:</span>
                <span class="info-value">` + activo.Tipo + `</span>
            </div>`

	if activo.Descripcion != nil && *activo.Descripcion != "" {
		htmlContent += `
            <div class="info-row">
                <span class="info-label">Descripción:</span>
                <span class="info-value">` + *activo.Descripcion + `</span>
            </div>`
	}

	if activo.Ubicacion != "" {
		htmlContent += `
            <div class="info-row">
                <span class="info-label">Ubicación:</span>
                <span class="info-value">` + activo.Ubicacion + `</span>
            </div>`
	}

	htmlContent += `
        </div>

        <div class="qr-container">
            <img src="data:image/png;base64,` + qrBase64 + `" alt="Código QR" class="qr-image">
        </div>

        <div class="buttons">
            <a href="/qr/obtener/` + strconv.Itoa(activo.ID) + `" download="QR_` + *activo.CodigoActivo + `.png" class="btn btn-primary">
                💾 Descargar QR
            </a>
            <button onclick="window.print()" class="btn btn-secondary">
                🖨️ Imprimir
            </button>
        </div>

        <div class="footer">
            <p>InspectAR © 2025 - Sistema de Gestión de Activos</p>
        </div>
    </div>
</body>
</html>`

	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, htmlContent)
}
