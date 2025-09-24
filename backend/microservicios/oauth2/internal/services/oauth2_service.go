package services

import (
	"errors"
	"fmt"
	"log"
	"oauth2/internal/models"
	"oauth2/internal/repository"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	Repo      *repository.UserRepository
	TokenRepo *repository.TokenRepository
	Config    *models.Config // Ahora usamos la configuración cargada desde ConfigService
}

func NewAuthService(repo *repository.UserRepository, tokenRepo *repository.TokenRepository, config *models.Config) *AuthService {
	return &AuthService{
		Repo:      repo,
		TokenRepo: tokenRepo,
		Config:    config,
	}
}

// Autentica el usuario y verifica la contraseña
func (s *AuthService) Authenticate(username, password string, userdb *models.User) error {
	if err := bcrypt.CompareHashAndPassword([]byte(userdb.Password), []byte(password)); err != nil {
		log.Printf("#2 %v", err)
		//log.Printf("User Password: %v", user.Password)
		//log.Printf("Password: %v", password)
		//gener, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
		//log.Printf("generado: %v", string(gener))
		return errors.New("credenciales inválidas #2")
	}

	return nil
}

// VerificarDatos verifica las credenciales y retorna el usuario si son válidas
func (s *AuthService) VerificarDatos(Email, password string) (*models.User, error) {
	user, err := s.Repo.GetUserByEmail(Email)
	if err != nil {
		return nil, errors.New("credenciales inválidas #3")
	}

	if err := s.Authenticate(user.Username, password, user); err != nil {
		return nil, errors.New("credenciales inválidas #4")
	}

	return user, nil
}

// Validar token JWT y extraer claims
func (s *AuthService) ValidateToken(tokenString string) (jwt.MapClaims, error) {
	// Cargar la clave pública desde el archivo
	publicKeyData, err := os.ReadFile("config/jwt_access_public.pem")
	if err != nil {
		return nil, fmt.Errorf("error leyendo clave pública: %v", err)
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyData)
	if err != nil {
		return nil, fmt.Errorf("error parseando clave pública: %v", err)
	}

	// Parsear y validar el token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Asegurar que el método de firma sea RS256
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("método de firma inesperado: %v", token.Header["alg"])
		}
		return publicKey, nil
	})

	if err != nil {
		return nil, err
	}

	// Extraer claims del token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("token inválido")
	}

	return claims, nil
}

func (s *AuthService) GenerateTokens(payload *models.TokenPayload) (string, string, error) {
	// Cargar la clave privada para Access Token
	accessPrivateKeyData, err := os.ReadFile("config/jwt_access_private.pem")
	if err != nil {
		return "", "", fmt.Errorf("error leyendo clave privada access: %v", err)
	}

	accessPrivateKey, err := jwt.ParseRSAPrivateKeyFromPEM(accessPrivateKeyData)
	if err != nil {
		return "", "", fmt.Errorf("error parseando clave privada access: %v", err)
	}

	// Generar Access Token
	accessToken := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"username": payload.Username,
		"email":    payload.Email,
		"empresa":  payload.Empresa,
		"device":   payload.Device,
		"scope":    payload.Scope,
		"exp":      time.Now().Add(time.Duration(s.Config.LifetimeAccess) * time.Minute).Unix(),
	})

	accessTokenString, err := accessToken.SignedString(accessPrivateKey)
	if err != nil {
		return "", "", err
	}

	// Cargar la clave privada para Refresh Token
	refreshPrivateKeyData, err := os.ReadFile("config/jwt_refresh_private.pem")
	if err != nil {
		return "", "", fmt.Errorf("error leyendo clave privada refresh: %v", err)
	}

	refreshPrivateKey, err := jwt.ParseRSAPrivateKeyFromPEM(refreshPrivateKeyData)
	if err != nil {
		return "", "", fmt.Errorf("error parseando clave privada refresh: %v", err)
	}

	user, err := s.Repo.GetUserByEmail(payload.Email)
	// if err != nil {
	// 	return nil, errors.New("credenciales inválidas #3")
	// }

	// Verificar si ya existe un refresh token
	var refreshTokenString string
	refreshTokenString, err = s.TokenRepo.GetRefreshToken(user.ID, payload.Device)
	if err != nil {
		refreshToken := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
			"username": payload.Username,
			"email":    payload.Email,
			"empresa":  payload.Empresa,
			"device":   payload.Device,
			"scope":    payload.Scope,
			"exp":      time.Now().Add(time.Duration(s.Config.LifetimeRefresh) * time.Minute).Unix(),
		})

		refreshTokenString, err = refreshToken.SignedString(refreshPrivateKey)
		if err != nil {
			return "", "", err
		}

		// Guardar refresh token en BD
		err = s.TokenRepo.SaveRefreshToken(user.ID, payload.Device, refreshTokenString, time.Now().Add(time.Duration(s.Config.LifetimeRefresh)*time.Minute))
		if err != nil {
			return "", "", err
		}
	}

	return accessTokenString, refreshTokenString, nil
}

// GetPayload extrae la información del token sin validación y la devuelve como TokenPayload
func (s *AuthService) GetPayload(tokenString string) (*models.TokenPayload, error) {
	// Parsear el token sin validar
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		log.Printf("Error al parsear el token: %v", err)
		return nil, fmt.Errorf("error al parsear el token: %w", err)
	}

	// Obtener claims sin validación
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("no se pudieron obtener los claims del token")
	}

	// Construcción del payload asegurando la conversión correcta de tipos
	payload := &models.TokenPayload{
		Username: claims["username"].(string),
		Email:    claims["email"].(string),
		Empresa:  claims["empresa"].(string),
		Device:   claims["device"].(string),
		Scope:    claims["scope"].(string),
	}

	// Manejo seguro del campo Expiration como interface{}
	if exp, exists := claims["exp"]; exists {
		payload.Expiration = exp
	}

	return payload, nil
}

// Valida un refresh token comparándolo con el almacenado en la base de datos
func (s *AuthService) ValidateRefreshToken(storedToken, requestToken string) error {
	if storedToken == "" || requestToken == "" {
		return errors.New("tokens no pueden estar vacíos")
	}

	if storedToken != requestToken {
		return errors.New("token invalido")
	}

	return nil
}

// Al cambiar la contrasena, elimina y invalida todos los refresh tokens del usuario
func (s *AuthService) ChangePasswordToken(email string) error {
	// Antes de entrar a esta funcion, debe de haberse actualizado la contrasena desde otro servicio

	user, err := s.Repo.GetUserByEmail(email)
	if err != nil {
		return errors.New("credenciales inválidas #3")
	}

	// Invalidar todos los refresh tokens del usuario
	err = s.TokenRepo.DeleteAllTokensForUser(user.ID)
	if err != nil {
		return err
	}

	return nil
}

// Cierra sesión eliminando el refresh token correspondiente al usuario y dispositivo
func (s *AuthService) Logout(email, deviceID string) error {
	//log.Printf("Cerrando sesión para el usuario con ID: %d", userID)

	user, err := s.Repo.GetUserByEmail(email)
	if err != nil {
		return errors.New("credenciales inválidas #3")
	}


	err = s.TokenRepo.DeleteRefreshToken(user.ID, deviceID)
	if err != nil {
		//log.Printf("Error al eliminar tokens para el usuario con ID: %d, error: %v", userID, err)
		return err
	}
	//log.Printf("Sesión cerrada correctamente para el usuario con ID: %d", userID)
	return nil
}

// Cuenta la cantidad de tokens de refresh de un usuario
func (s *AuthService) CountTokens(userID uint) (int, error) {
	var count int
	err := s.TokenRepo.DB.Get(&count, "SELECT COUNT(*) FROM refresh_tokens WHERE user_id = $1", userID)
	return count, err
}

// Verifica y limpia tokens si exceden el límite
func (s *AuthService) CheckTokens(userID uint) error {
	count, err := s.CountTokens(userID)
	if err != nil {
		return fmt.Errorf("error al contar tokens: %v", err)
	}

	if count >= 8 {
		err = s.TokenRepo.DeleteRefreshTokenOld(userID)
		if err != nil {
			return fmt.Errorf("error al eliminar token antiguo: %v", err)
		}
		log.Printf("Token antiguo eliminado para usuario %d por exceder límite", userID)
	}

	return nil
}

// GetRefreshToken obtiene el token de refresco usando el payload del token
func (s *AuthService) GetRefreshToken(payload *models.TokenPayload) (string, error) {
	// Buscar el usuario por email para obtener su ID
	user, err := s.Repo.GetUserByEmail(payload.Email)
	if err != nil {
		return "", fmt.Errorf("error al obtener usuario: %w", err)
	}

	// Obtener el token de refresco usando el ID del usuario y el device
	refreshToken, err := s.TokenRepo.GetRefreshToken(user.ID, payload.Device)
	if err != nil {
		return "", fmt.Errorf("error al obtener refresh token: %w", err)
	}

	return refreshToken, nil
}
