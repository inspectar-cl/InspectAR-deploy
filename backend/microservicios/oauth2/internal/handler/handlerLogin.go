package handler

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"time"

	"oauth2/internal/dto"
	"oauth2/internal/models"
	"oauth2/internal/services"

	"github.com/gin-gonic/gin"
)

var TokenExpiredError = errors.New("token expirado")
var tokenFirst = errors.New("primer login")

type LoginHandler struct {
	AuthService *services.AuthService
}

func NewLoginHandler(authService *services.AuthService) *LoginHandler {
	return &LoginHandler{AuthService: authService}
}

func (h *LoginHandler) Login(c *gin.Context) {
	var loginDto dto.LoginRequest
	if err := c.ShouldBindJSON(&loginDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos de login inválidos"})
		return
	}

	user, err := h.AuthService.VerificarDatos(loginDto.Email, loginDto.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	tokenPayload := &models.TokenPayload{
		Username:   user.Username,
		Email:      user.Email,
		Device:     loginDto.DeviceID,
		Scope:      user.Scope,
		Expiration: time.Now().Add(15 * time.Minute).Unix(), //por ahora no se ocupa
	}

	accessToken, refreshToken, err := h.AuthService.GenerateTokens(tokenPayload)
	h.AuthService.CheckTokens(user.ID) //un usuario no puede tener mas 8 de un tokens refresh
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generando tokens"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

// Refrescar el token de acceso usando el refresh token
func (c LoginHandler) RefreshToken(ctx *gin.Context) {
	var refreshDto dto.RefreshTokenRequest
	if err := ctx.ShouldBindJSON(&refreshDto); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	tokenPayload, err := c.AuthService.GetPayload(refreshDto.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido o expirado"})
		return
	}

	// Obtener el token almacenado usando userID y deviceID
	storedToken, err := c.AuthService.GetRefreshToken(tokenPayload)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "No existe token para este usuario"})
		case TokenExpiredError:
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Token expirado"})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener el token"})
		}
		return
	}

	// Validar que el token coincida con el almacenado
	err = c.AuthService.ValidateRefreshToken(storedToken, refreshDto.RefreshToken)
	if err != nil {
		switch err.Error() {
		case "tokens no pueden estar vacíos":
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case "token invalido":
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Token no coincide con el almacenado"})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error al validar el token"})
		}
		return
	}

	// Generar nuevos tokens
	accessToken, refreshToken, err := c.AuthService.GenerateTokens(tokenPayload)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error generando nuevos tokens"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

// cambio de formato
// Al cambiar la contrasena, invalida todos los refresh tokens del usuario
func (c *LoginHandler) ChangePasswordToken(ctx *gin.Context) {
	email, exists := ctx.Get("email")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "No autorizado"})
		return
	}

	// Elimina los tokens de acceso y refresco al cambiar contrasena
	err := c.AuthService.ChangePasswordToken(email.(string))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error cambiando la contraseña"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Sesiones cerradas"})
}

func (c *LoginHandler) Logout(ctx *gin.Context) {
	email, exists := ctx.Get("email")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "No autorizado"})
		return
	}

	deviceID, exists := ctx.Get("device")
	if !exists {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID de dispositivo no proporcionado"})
		return
	}

	err := c.AuthService.Logout(email.(string), deviceID.(string))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error cerrando sesión"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Sesión cerrada exitosamente"})
}
