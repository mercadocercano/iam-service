package usecase_test

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"iam/src/auth/application/usecase"
	"iam/src/auth/domain/value_object"
	tenant_vo "iam/src/tenant/domain/value_object"
)

func generateTestToken(secret string, claims *value_object.TokenClaims) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, _ := token.SignedString([]byte(secret))
	return tokenStr
}

func TestValidateTokenUseCase_Execute_ValidToken_ReturnsClaims(t *testing.T) {
	// Arrange
	secret := "test-secret-key-for-testing-purposes"
	config := usecase.AuthConfig{
		JWTSecret: secret,
	}

	validateUseCase := usecase.NewValidateTokenUseCase(config)

	userID := uuid.New()
	tenantID := uuid.New()
	roleID := uuid.New()

	claims := value_object.NewTokenClaims(
		userID,
		tenantID,
		roleID,
		"test@example.com",
		tenant_vo.NewTenantFeatures(),
		time.Now().Add(15*time.Minute),
	)

	tokenStr := generateTestToken(secret, claims)

	// Act
	result, err := validateUseCase.Execute(tokenStr)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, userID, result.UserID)
	assert.Equal(t, tenantID, result.TenantID)
	assert.Equal(t, roleID, result.RoleID)
	assert.Equal(t, "test@example.com", result.Email)
}

func TestValidateTokenUseCase_Execute_ExpiredToken_ReturnsError(t *testing.T) {
	// Arrange
	secret := "test-secret-key-for-testing-purposes"
	config := usecase.AuthConfig{
		JWTSecret: secret,
	}

	validateUseCase := usecase.NewValidateTokenUseCase(config)

	claims := value_object.NewTokenClaims(
		uuid.New(),
		uuid.New(),
		uuid.New(),
		"test@example.com",
		tenant_vo.NewTenantFeatures(),
		time.Now().Add(-1*time.Hour), // Ya expirado
	)

	tokenStr := generateTestToken(secret, claims)

	// Act
	result, err := validateUseCase.Execute(tokenStr)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "token")
}

func TestValidateTokenUseCase_Execute_InvalidFormat_ReturnsError(t *testing.T) {
	// Arrange
	config := usecase.AuthConfig{
		JWTSecret: "test-secret",
	}

	validateUseCase := usecase.NewValidateTokenUseCase(config)

	// Act
	result, err := validateUseCase.Execute("not-a-valid-jwt-token")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestValidateTokenUseCase_Execute_WrongSecret_ReturnsError(t *testing.T) {
	// Arrange
	config := usecase.AuthConfig{
		JWTSecret: "correct-secret",
	}

	validateUseCase := usecase.NewValidateTokenUseCase(config)

	claims := value_object.NewTokenClaims(
		uuid.New(),
		uuid.New(),
		uuid.New(),
		"test@example.com",
		tenant_vo.NewTenantFeatures(),
		time.Now().Add(15*time.Minute),
	)

	// Firmar con un secret diferente
	tokenStr := generateTestToken("wrong-secret", claims)

	// Act
	result, err := validateUseCase.Execute(tokenStr)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestValidateTokenUseCase_Execute_EmptyToken_ReturnsError(t *testing.T) {
	// Arrange
	config := usecase.AuthConfig{
		JWTSecret: "test-secret",
	}

	validateUseCase := usecase.NewValidateTokenUseCase(config)

	// Act
	result, err := validateUseCase.Execute("")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
}
