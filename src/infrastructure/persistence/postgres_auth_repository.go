package persistence

import (
	"context"
	"iam/src/domain/models"
	"iam/src/domain/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PostgresAuthRepository struct {
	db *gorm.DB
}

func NewPostgresAuthRepository(db *gorm.DB) repositories.AuthRepository {
	return &PostgresAuthRepository{db: db}
}

// Refresh Tokens
func (r *PostgresAuthRepository) CreateRefreshToken(ctx context.Context, token *models.RefreshToken) error {
	return r.db.WithContext(ctx).Create(token).Error
}

func (r *PostgresAuthRepository) GetRefreshToken(ctx context.Context, token string) (*models.RefreshToken, error) {
	var refreshToken models.RefreshToken
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("User.Tenant").
		Preload("User.Role").
		Where("token = ? AND expires_at > NOW()", token).
		First(&refreshToken).Error
	if err != nil {
		return nil, err
	}
	return &refreshToken, nil
}

func (r *PostgresAuthRepository) DeleteRefreshToken(ctx context.Context, token string) error {
	return r.db.WithContext(ctx).
		Where("token = ?", token).
		Delete(&models.RefreshToken{}).Error
}

func (r *PostgresAuthRepository) DeleteAllUserRefreshTokens(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Delete(&models.RefreshToken{}).Error
}

// Federated Auth
func (r *PostgresAuthRepository) GetUserByFederatedID(ctx context.Context, provider models.AuthProvider, federatedID string, tenantID *uuid.UUID) (*models.User, error) {
	var user models.User
	query := r.db.WithContext(ctx).
		Preload("Tenant").
		Preload("Role").
		Where("provider = ? AND federated_id = ?", provider, federatedID)
	if tenantID != nil {
		query = query.Where("tenant_id = ?", tenantID)
	}
	err := query.First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *PostgresAuthRepository) LinkFederatedID(ctx context.Context, userID uuid.UUID, provider models.AuthProvider, federatedID string) error {
	return r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"provider":     provider,
			"federated_id": federatedID,
		}).Error
}
