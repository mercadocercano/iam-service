package repositories

import (
    "context"
    "github.com/google/uuid"
    "iam/src/domain/models"
)

type PlanRepository interface {
    Create(ctx context.Context, plan *models.Plan) error
    Update(ctx context.Context, plan *models.Plan) error
    Delete(ctx context.Context, id uuid.UUID) error
    GetByID(ctx context.Context, id uuid.UUID) (*models.Plan, error)
    GetAll(ctx context.Context) ([]models.Plan, error)
    GetBySaas(ctx context.Context, saas models.SaasType) ([]models.Plan, error)
}
