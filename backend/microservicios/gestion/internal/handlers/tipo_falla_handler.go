package handlers

import (
	"gestion/internal/models"
	"gestion/internal/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TipoFallaHandler struct {
	repo *repository.TipoFallaRepository
}

func NewTipoFallaHandler(repo *repository.TipoFallaRepository) *TipoFallaHandler {
	return &TipoFallaHandler{repo: repo}
}

// CrearTipoFalla crea un nuevo reporte de falla
// POST /tipos-falla
func (h *TipoFallaHandler) CrearTipoFalla(c *gin.Context) {
	var req models.CreateTipoFallaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos", "details": err.Error()})
		return
	}

	// Validar tipo de falla
	tiposValidos := []string{"falla agua", "falla ascensor", "falla electricidad", "falla caldera"}
	tipoValido := false
	for _, t := range tiposValidos {
		if req.Tipo == t {
			tipoValido = true
			break
		}
	}

	if !tipoValido {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tipo de falla inválido. Debe ser: falla agua, falla ascensor, falla electricidad o falla caldera"})
		return
	}

	tipoFalla, err := h.repo.CrearTipoFalla(req.Email, req.Tipo, req.Descripcion, req.IDEdificio)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear tipo de falla", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"mensaje": "Tipo de falla creado exitosamente",
		"data":    tipoFalla,
	})
}

// CrearComentario crea un nuevo comentario sobre una falla
// POST /comentarios
func (h *TipoFallaHandler) CrearComentario(c *gin.Context) {
	var req models.CreateComentarioRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos", "details": err.Error()})
		return
	}

	comentario, err := h.repo.CrearComentario(req.Email, req.IDFalla, req.Comentario)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear comentario", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"mensaje": "Comentario creado exitosamente",
		"data":    comentario,
	})
}

// ObtenerTiposFallaPorEdificio obtiene los tipos de falla con sus comentarios
// Si se proporciona el parámetro ?pagina=N, devuelve resultados paginados (10 por página)
// Si NO se proporciona pagina, devuelve TODOS los resultados sin paginación
// GET /tipos-falla/edificio/:edificio_id
// GET /tipos-falla/edificio/:edificio_id?pagina=1
func (h *TipoFallaHandler) ObtenerTiposFallaPorEdificio(c *gin.Context) {
	edificioIDStr := c.Param("edificio_id")
	edificioID, err := strconv.Atoi(edificioIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de edificio inválido"})
		return
	}

	// Verificar si se proporcionó el parámetro de página
	paginaStr := c.Query("pagina")

	// Si NO hay parámetro de página, devolver TODOS los resultados
	if paginaStr == "" {
		tiposFalla, err := h.repo.ObtenerTodosTiposFallaPorEdificio(edificioID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener tipos de falla", "details": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"edificio_id": edificioID,
			"total_items": len(tiposFalla),
			"data":        tiposFalla,
		})
		return
	}

	// Si hay parámetro de página, usar paginación
	pagina, err := strconv.Atoi(paginaStr)
	if err != nil || pagina < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Número de página inválido"})
		return
	}

	tiposFalla, err := h.repo.ObtenerTiposFallaPorEdificio(edificioID, pagina)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener tipos de falla", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"pagina":         pagina,
		"edificio_id":    edificioID,
		"total_items":    len(tiposFalla),
		"items_per_page": 10,
		"data":           tiposFalla,
	})
}
