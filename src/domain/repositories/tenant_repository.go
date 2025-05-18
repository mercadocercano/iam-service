package repositories

import (
    "context"
    "github.com/google/uuid"
    "iam/src/domain/models"
)

type TenantRepository interface {
    Create(ctx context.Context, tenant *models.Tenant) error
    Update(ctx context.Context, tenant *models.Tenant) error
    Delete(ctx context.Context, id uuid.UUID) error
    GetByID(ctx context.Context, id uuid.UUID) (*models.Tenant, error)
    GetAll(ctx context.Context) ([]models.Tenant, error)
    GetBySaas(ctx context.Context, saas models.SaasType) ([]models.Tenant, error)
    GetByEmailUserKey(ctx context.Context, emailUserKey string) (*models.Tenant, error)
}
