package handlers

import (
	"gestion/internal/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

type EdificioHandler struct {
	edificioRepo *repository.EdificioRepository
}

func NewEdificioHandler(edificioRepo *repository.EdificioRepository) *EdificioHandler {
	return &EdificioHandler{
		edificioRepo: edificioRepo,
	}
}

// GetAllEdificios obtiene todos los edificios
func (h *EdificioHandler) GetAllEdificios(c *gin.Context) {
	edificios, err := h.edificioRepo.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error al obtener edificios",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"edificios": edificios,
		"total":     len(edificios),
	})
}
