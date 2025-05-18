package persistence

import (
	"context"
	"fmt"
	"iam/src/domain/models"
	"iam/src/domain/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PostgresUserRepository struct {
	db *gorm.DB
}

func NewPostgresUserRepository(db *gorm.DB) repositories.UserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) Create(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *PostgresUserRepository) Update(ctx context.Context, update *models.UserUpdate) error {
	// Verificar que el usuario existe
	var exists bool
	sqlDB, err := r.db.DB()
	if err != nil {
		return err
	}

	err = sqlDB.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)", update.ID).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return gorm.ErrRecordNotFound
	}

	// Iniciar transacción
	tx, err := sqlDB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Construir la consulta SQL con solo los campos permitidos
	query := "UPDATE users SET updated_at = NOW()"
	args := make([]interface{}, 0)
	paramCount := 1

	if update.Email != "" {
		query += fmt.Sprintf(", email = $%d", paramCount)
		args = append(args, update.Email)
		paramCount++
	}
	if update.RoleID != uuid.Nil {
		query += fmt.Sprintf(", role_id = $%d", paramCount)
		args = append(args, update.RoleID)
		paramCount++
	}
	if update.Status != "" {
		query += fmt.Sprintf(", status = $%d::user_status", paramCount)
		args = append(args, update.Status)
		paramCount++
	}

	query += fmt.Sprintf(" WHERE id = $%d RETURNING id", paramCount)
	args = append(args, update.ID)

	// Ejecutar la actualización
	var updatedID uuid.UUID
	err = tx.QueryRowContext(ctx, query, args...).Scan(&updatedID)
	if err != nil {
		return err
	}

	// Confirmar la transacción
	return tx.Commit()
}

func (r *PostgresUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.User{}, id).Error
}

func (r *PostgresUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Preload("Tenant").Preload("Role").First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *PostgresUserRepository) GetByEmail(ctx context.Context, email string, tenantID *uuid.UUID) (*models.User, error) {
	var user models.User
	query := r.db.WithContext(ctx).Preload("Tenant").Preload("Role").Where("email = ?", email)
	if tenantID != nil {
		query = query.Where("tenant_id = ?", tenantID)
	}
	err := query.First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *PostgresUserRepository) GetAll(ctx context.Context) ([]models.User, error) {
	var users []models.User
	err := r.db.WithContext(ctx).Preload("Tenant").Preload("Role").Find(&users).Error
	return users, err
}

func (r *PostgresUserRepository) GetByTenant(ctx context.Context, tenantID uuid.UUID) ([]models.User, error) {
	var users []models.User
	err := r.db.WithContext(ctx).Preload("Tenant").Preload("Role").
		Where("tenant_id = ?", tenantID).
		Find(&users).Error
	return users, err
}
