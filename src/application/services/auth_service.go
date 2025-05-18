package services

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"encoding/json"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"iam/src/domain/models"
	"iam/src/domain/repositories"
)

var (
	ErrInvalidCredentials = errors.New("credenciales inválidas")
	ErrUserNotFound       = errors.New("usuario no encontrado")
	ErrInvalidToken       = errors.New("token inválido")
	ErrExpiredToken       = errors.New("token expirado")
)

type AuthConfig struct {
	JWTSecret          string
	AccessTokenExpiry  time.Duration
	RefreshTokenExpiry time.Duration
	GoogleClientID     string
}

type AuthService struct {
	config     AuthConfig
	userRepo   repositories.UserRepository
	authRepo   repositories.AuthRepository
	httpClient *http.Client
}

func NewAuthService(
	config AuthConfig,
	userRepo repositories.UserRepository,
	authRepo repositories.AuthRepository,
) (*AuthService, error) {
	return &AuthService{
		config:     config,
		userRepo:   userRepo,
		authRepo:   authRepo,
		httpClient: &http.Client{},
	}, nil
}

// Login maneja tanto el login local como el federado
func (s *AuthService) Login(ctx context.Context, req *models.LoginRequest) (*models.LoginResponse, error) {
	var user *models.User
	var err error

	switch req.Provider {
	case models.LocalAuth:
		user, err = s.loginLocal(ctx, req)
	case models.GoogleAuth:
		user, err = s.loginGoogle(ctx, req)
	default:
		return nil, fmt.Errorf("proveedor de autenticación no soportado: %s", req.Provider)
	}

	if err != nil {
		return nil, err
	}

	// Generar tokens
	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return nil, fmt.Errorf("error generando access token: %w", err)
	}

	refreshToken, err := s.generateRefreshToken(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("error generando refresh token: %w", err)
	}

	return &models.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int(s.config.AccessTokenExpiry.Seconds()),
		TokenType:    "Bearer",
		User:         *user,
	}, nil
}

// Login local con email y password
func (s *AuthService) loginLocal(ctx context.Context, req *models.LoginRequest) (*models.User, error) {
	user, err := s.userRepo.GetByEmail(ctx, req.Email, req.TenantID)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if user.Provider != models.LocalAuth {
		return nil, fmt.Errorf("este usuario usa autenticación %s", user.Provider)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	// Si no se proporcionó tenant_id, usar el del usuario encontrado
	if req.TenantID == nil {
		req.TenantID = &user.TenantID
	} else if *req.TenantID != user.TenantID {
		return nil, ErrInvalidCredentials
	}

	return user, nil
}

// Login con Google
func (s *AuthService) loginGoogle(ctx context.Context, req *models.LoginRequest) (*models.User, error) {
	// Verificar token de Google usando la API pública
	resp, err := s.httpClient.Get("https://oauth2.googleapis.com/tokeninfo?id_token=" + req.GoogleToken)
	if err != nil {
		return nil, fmt.Errorf("error verificando token de Google: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("token de Google inválido")
	}

	var tokenInfo struct {
		Aud   string `json:"aud"`
		Sub   string `json:"sub"`
		Email string `json:"email"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tokenInfo); err != nil {
		return nil, fmt.Errorf("error decodificando respuesta de Google: %w", err)
	}

	if tokenInfo.Aud != s.config.GoogleClientID {
		return nil, errors.New("client ID de Google inválido")
	}

	// Buscar usuario por ID federado
	user, err := s.authRepo.GetUserByFederatedID(ctx, models.GoogleAuth, tokenInfo.Sub, req.TenantID)
	if err == nil {
		// Si se encontró el usuario y no se proporcionó tenant_id, usar el del usuario
		if req.TenantID == nil {
			req.TenantID = &user.TenantID
		} else if *req.TenantID != user.TenantID {
			return nil, ErrInvalidCredentials
		}
		return user, nil
	}

	// Si no existe, buscar por email
	user, err = s.userRepo.GetByEmail(ctx, tokenInfo.Email, req.TenantID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	// Si no se proporcionó tenant_id, usar el del usuario encontrado
	if req.TenantID == nil {
		req.TenantID = &user.TenantID
	} else if *req.TenantID != user.TenantID {
		return nil, ErrInvalidCredentials
	}

	// Vincular ID federado si el usuario existe
	if err := s.authRepo.LinkFederatedID(ctx, user.ID, models.GoogleAuth, tokenInfo.Sub); err != nil {
		return nil, fmt.Errorf("error vinculando ID federado: %w", err)
	}

	return user, nil
}

// RefreshToken genera un nuevo access token usando un refresh token
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*models.LoginResponse, error) {
	// Obtener refresh token de la base de datos
	token, err := s.authRepo.GetRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, ErrInvalidToken
	}

	if token.ExpiresAt.Before(time.Now()) {
		_ = s.authRepo.DeleteRefreshToken(ctx, refreshToken)
		return nil, ErrExpiredToken
	}

	// Generar nuevo access token
	accessToken, err := s.generateAccessToken(&token.User)
	if err != nil {
		return nil, fmt.Errorf("error generando access token: %w", err)
	}

	// Generar nuevo refresh token
	newRefreshToken, err := s.generateRefreshToken(ctx, &token.User)
	if err != nil {
		return nil, fmt.Errorf("error generando refresh token: %w", err)
	}

	// Eliminar el refresh token anterior
	_ = s.authRepo.DeleteRefreshToken(ctx, refreshToken)

	return &models.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    int(s.config.AccessTokenExpiry.Seconds()),
		TokenType:    "Bearer",
		User:         token.User,
	}, nil
}

// ValidateToken valida un access token y retorna sus claims
func (s *AuthService) ValidateToken(tokenString string) (*models.TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &models.TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de firma inesperado: %v", token.Header["alg"])
		}
		return []byte(s.config.JWTSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("error parseando token: %w", err)
	}

	if claims, ok := token.Claims.(*models.TokenClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidToken
}

// Logout invalida todos los refresh tokens de un usuario
func (s *AuthService) Logout(ctx context.Context, userID uuid.UUID) error {
	return s.authRepo.DeleteAllUserRefreshTokens(ctx, userID)
}

// Funciones auxiliares

func (s *AuthService) generateAccessToken(user *models.User) (string, error) {
	claims := models.TokenClaims{
		UserID:    user.ID,
		Email:     user.Email,
		TenantID:  user.TenantID,
		RoleID:    user.RoleID,
		ExpiresAt: time.Now().Add(s.config.AccessTokenExpiry).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.JWTSecret))
}

func (s *AuthService) generateRefreshToken(ctx context.Context, user *models.User) (string, error) {
	// Generar token aleatorio
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	token := base64.URLEncoding.EncodeToString(b)

	// Crear refresh token en la base de datos
	refreshToken := &models.RefreshToken{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(s.config.RefreshTokenExpiry),
	}

	if err := s.authRepo.CreateRefreshToken(ctx, refreshToken); err != nil {
		return "", err
	}

	return token, nil
}
