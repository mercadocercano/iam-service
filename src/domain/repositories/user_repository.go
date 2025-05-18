package repositories

import (
	"context"
	"iam/src/domain/models"

	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	Update(ctx context.Context, update *models.UserUpdate) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetAll(ctx context.Context) ([]models.User, error)
	GetByEmail(ctx context.Context, email string, tenantID *uuid.UUID) (*models.User, error)
	GetByTenant(ctx context.Context, tenantID uuid.UUID) ([]models.User, error)
}
