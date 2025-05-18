package services

import (
	"context"
	"iam/src/domain/models"
	"iam/src/domain/repositories"

	"github.com/google/uuid"
)

type RoleService struct {
	roleRepository repositories.RoleRepository
}

func NewRoleService(roleRepository repositories.RoleRepository) *RoleService {
	return &RoleService{
		roleRepository: roleRepository,
	}
}

func (s *RoleService) CreateRole(ctx context.Context, role *models.Role) error {
	return s.roleRepository.Create(ctx, role)
}

func (s *RoleService) UpdateRole(ctx context.Context, role *models.Role) error {
	return s.roleRepository.Update(ctx, role)
}

func (s *RoleService) DeleteRole(ctx context.Context, id uuid.UUID) error {
	return s.roleRepository.Delete(ctx, id)
}

func (s *RoleService) GetRoleByID(ctx context.Context, id uuid.UUID) (*models.Role, error) {
	return s.roleRepository.GetByID(ctx, id)
}

func (s *RoleService) GetAllRoles(ctx context.Context) ([]models.Role, error) {
	return s.roleRepository.GetAll(ctx)
}

func (s *RoleService) GetRolesBySaas(ctx context.Context, saas models.SaasType) ([]models.Role, error) {
	return s.roleRepository.GetBySaas(ctx, saas)
}

func (s *RoleService) GetRoleByName(ctx context.Context, name string) (*models.Role, error) {
	return s.roleRepository.GetByName(ctx, name)
}
