package repositories

import (
	"context"
	"iam/src/domain/models"

	"github.com/google/uuid"
)

type AuthRepository interface {
	// Refresh Tokens
	CreateRefreshToken(ctx context.Context, token *models.RefreshToken) error
	GetRefreshToken(ctx context.Context, token string) (*models.RefreshToken, error)
	DeleteRefreshToken(ctx context.Context, token string) error
	DeleteAllUserRefreshTokens(ctx context.Context, userID uuid.UUID) error

	// Federated Auth
	GetUserByFederatedID(ctx context.Context, provider models.AuthProvider, federatedID string, tenantID *uuid.UUID) (*models.User, error)
	LinkFederatedID(ctx context.Context, userID uuid.UUID, provider models.AuthProvider, federatedID string) error
}
