package persistence

import (
	"context"
	"iam/src/domain/models"
	"iam/src/domain/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PostgresRoleRepository struct {
	db *gorm.DB
}

func NewPostgresRoleRepository(db *gorm.DB) repositories.RoleRepository {
	return &PostgresRoleRepository{db: db}
}

func (r *PostgresRoleRepository) Create(ctx context.Context, role *models.Role) error {
	return r.db.WithContext(ctx).Create(role).Error
}

func (r *PostgresRoleRepository) Update(ctx context.Context, role *models.Role) error {
	return r.db.WithContext(ctx).Save(role).Error
}

func (r *PostgresRoleRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Role{}, id).Error
}

func (r *PostgresRoleRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Role, error) {
	var role models.Role
	err := r.db.WithContext(ctx).First(&role, id).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *PostgresRoleRepository) GetAll(ctx context.Context) ([]models.Role, error) {
	var roles []models.Role
	err := r.db.WithContext(ctx).Find(&roles).Error
	return roles, err
}

func (r *PostgresRoleRepository) GetBySaas(ctx context.Context, saas models.SaasType) ([]models.Role, error) {
	var roles []models.Role
	err := r.db.WithContext(ctx).Where("saas = ?", saas).Find(&roles).Error
	return roles, err
}

func (r *PostgresRoleRepository) GetByName(ctx context.Context, name string) (*models.Role, error) {
	var role models.Role
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}
