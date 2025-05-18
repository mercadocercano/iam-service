package persistence

import (
	"context"
	"iam/src/domain/models"
	"iam/src/domain/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PostgresTenantRepository struct {
	db *gorm.DB
}

func NewPostgresTenantRepository(db *gorm.DB) repositories.TenantRepository {
	return &PostgresTenantRepository{db: db}
}

func (r *PostgresTenantRepository) Create(ctx context.Context, tenant *models.Tenant) error {
	return r.db.WithContext(ctx).Create(tenant).Error
}

func (r *PostgresTenantRepository) Update(ctx context.Context, tenant *models.Tenant) error {
	return r.db.WithContext(ctx).Save(tenant).Error
}

func (r *PostgresTenantRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Tenant{}, id).Error
}

func (r *PostgresTenantRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Tenant, error) {
	var tenant models.Tenant
	err := r.db.WithContext(ctx).Preload("Plan").First(&tenant, id).Error
	if err != nil {
		return nil, err
	}
	return &tenant, nil
}

func (r *PostgresTenantRepository) GetAll(ctx context.Context) ([]models.Tenant, error) {
	var tenants []models.Tenant
	err := r.db.WithContext(ctx).Preload("Plan").Find(&tenants).Error
	return tenants, err
}

func (r *PostgresTenantRepository) GetBySaas(ctx context.Context, saas models.SaasType) ([]models.Tenant, error) {
	var tenants []models.Tenant
	err := r.db.WithContext(ctx).Preload("Plan").Where("saas = ?", saas).Find(&tenants).Error
	return tenants, err
}

func (r *PostgresTenantRepository) GetByEmailUserKey(ctx context.Context, emailUserKey string) (*models.Tenant, error) {
	var tenant models.Tenant
	err := r.db.WithContext(ctx).Preload("Plan").Where("email_user_key = ?", emailUserKey).First(&tenant).Error
	if err != nil {
		return nil, err
	}
	return &tenant, nil
}
