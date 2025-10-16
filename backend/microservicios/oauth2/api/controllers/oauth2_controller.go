package controllers

import (
	"net/http"
	"oauth2/internal/dto"
	"oauth2/internal/services"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	Service *services.AuthService
}

func NewAuthController(service *services.AuthService) *AuthController {
	return &AuthController{Service: service}
}

func (c *AuthController) Login(ctx *gin.Context) {
	var loginDto dto.LoginRequest
	if err := ctx.ShouldBindJSON(&loginDto); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	// token, err := c.Service.Authenticate(loginDto.Username, loginDto.Password)
	// if err != nil {
	// 	ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciales inválidas #4"})
	// 	return
	// }

	// ctx.JSON(http.StatusOK, gin.H{"token": token})
}
