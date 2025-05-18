package persistence

import (
	"context"
	"fmt"
	"iam/src/domain/models"
	"iam/src/domain/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PostgresPlanRepository struct {
	db *gorm.DB
}

func NewPostgresPlanRepository(db *gorm.DB) repositories.PlanRepository {
	return &PostgresPlanRepository{db: db}
}

func (r *PostgresPlanRepository) Create(ctx context.Context, plan *models.Plan) error {
	return r.db.WithContext(ctx).Create(plan).Error
}

func (r *PostgresPlanRepository) Update(ctx context.Context, plan *models.Plan) error {
	// Primero verificar que el plan existe
	var existingPlan models.Plan
	if err := r.db.WithContext(ctx).First(&existingPlan, plan.ID).Error; err != nil {
		return fmt.Errorf("error finding plan: %v", err)
	}

	// Realizar la actualización
	result := r.db.WithContext(ctx).Model(plan).Updates(map[string]interface{}{
		"saas":          plan.Saas,
		"name":          plan.Name,
		"description":   plan.Description,
		"features":      plan.Features,
		"monthly_price": plan.MonthlyPrice,
		"yearly_price":  plan.YearlyPrice,
		"updated_at":    plan.UpdatedAt,
	})

	if result.Error != nil {
		return fmt.Errorf("error updating plan: %v", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("no rows were affected when updating plan ID: %v", plan.ID)
	}

	return nil
}

func (r *PostgresPlanRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Plan{}, id).Error
}

func (r *PostgresPlanRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Plan, error) {
	var plan models.Plan
	err := r.db.WithContext(ctx).First(&plan, id).Error
	if err != nil {
		return nil, err
	}
	return &plan, nil
}

func (r *PostgresPlanRepository) GetAll(ctx context.Context) ([]models.Plan, error) {
	var plans []models.Plan
	err := r.db.WithContext(ctx).Find(&plans).Error
	return plans, err
}

func (r *PostgresPlanRepository) GetBySaas(ctx context.Context, saas models.SaasType) ([]models.Plan, error) {
	var plans []models.Plan
	err := r.db.WithContext(ctx).Where("saas = ?", saas).Find(&plans).Error
	return plans, err
}
