package services

import (
    "context"
    "github.com/google/uuid"
    "iam/src/domain/models"
    "iam/src/domain/repositories"
)

type PlanService struct {
    planRepo repositories.PlanRepository
}

func NewPlanService(planRepo repositories.PlanRepository) *PlanService {
    return &PlanService{
        planRepo: planRepo,
    }
}

func (s *PlanService) CreatePlan(ctx context.Context, plan *models.Plan) error {
    return s.planRepo.Create(ctx, plan)
}

func (s *PlanService) UpdatePlan(ctx context.Context, plan *models.Plan) error {
    return s.planRepo.Update(ctx, plan)
}

func (s *PlanService) DeletePlan(ctx context.Context, id uuid.UUID) error {
    return s.planRepo.Delete(ctx, id)
}

func (s *PlanService) GetPlanByID(ctx context.Context, id uuid.UUID) (*models.Plan, error) {
    return s.planRepo.GetByID(ctx, id)
}

func (s *PlanService) GetAllPlans(ctx context.Context) ([]models.Plan, error) {
    return s.planRepo.GetAll(ctx)
}

func (s *PlanService) GetPlansBySaas(ctx context.Context, saas models.SaasType) ([]models.Plan, error) {
    return s.planRepo.GetBySaas(ctx, saas)
}
