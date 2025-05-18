package repositories

import (
    "context"
    "github.com/google/uuid"
    "iam/src/domain/models"
)

type RoleRepository interface {
    Create(ctx context.Context, role *models.Role) error
    Update(ctx context.Context, role *models.Role) error
    Delete(ctx context.Context, id uuid.UUID) error
    GetByID(ctx context.Context, id uuid.UUID) (*models.Role, error)
    GetAll(ctx context.Context) ([]models.Role, error)
    GetBySaas(ctx context.Context, saas models.SaasType) ([]models.Role, error)
    GetByName(ctx context.Context, name string) (*models.Role, error)
}
