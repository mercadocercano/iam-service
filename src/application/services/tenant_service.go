package services

import (
    "context"
    "github.com/google/uuid"
    "iam/src/domain/models"
    "iam/src/domain/repositories"
)

type TenantService struct {
    tenantRepo repositories.TenantRepository
}

func NewTenantService(tenantRepo repositories.TenantRepository) *TenantService {
    return &TenantService{
        tenantRepo: tenantRepo,
    }
}

func (s *TenantService) CreateTenant(ctx context.Context, tenant *models.Tenant) error {
    return s.tenantRepo.Create(ctx, tenant)
}

func (s *TenantService) UpdateTenant(ctx context.Context, tenant *models.Tenant) error {
    return s.tenantRepo.Update(ctx, tenant)
}

func (s *TenantService) DeleteTenant(ctx context.Context, id uuid.UUID) error {
    return s.tenantRepo.Delete(ctx, id)
}

func (s *TenantService) GetTenantByID(ctx context.Context, id uuid.UUID) (*models.Tenant, error) {
    return s.tenantRepo.GetByID(ctx, id)
}

func (s *TenantService) GetAllTenants(ctx context.Context) ([]models.Tenant, error) {
    return s.tenantRepo.GetAll(ctx)
}

func (s *TenantService) GetTenantsBySaas(ctx context.Context, saas models.SaasType) ([]models.Tenant, error) {
    return s.tenantRepo.GetBySaas(ctx, saas)
}

func (s *TenantService) GetTenantByEmailUserKey(ctx context.Context, emailUserKey string) (*models.Tenant, error) {
    return s.tenantRepo.GetByEmailUserKey(ctx, emailUserKey)
}
