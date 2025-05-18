package services

import (
	"context"
	"iam/src/domain/models"
	"iam/src/domain/repositories"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) CreateUser(ctx context.Context, user *models.User, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashedPassword)
	return s.userRepo.Create(ctx, user)
}

func (s *UserService) UpdateUser(ctx context.Context, update *models.UserUpdate) error {
	// Validar que el usuario existe antes de actualizar
	_, err := s.userRepo.GetByID(ctx, update.ID)
	if err != nil {
		return err
	}
	return s.userRepo.Update(ctx, update)
}

func (s *UserService) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return s.userRepo.Delete(ctx, id)
}

func (s *UserService) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	return s.userRepo.GetByID(ctx, id)
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string, tenantID *uuid.UUID) (*models.User, error) {
	return s.userRepo.GetByEmail(ctx, email, tenantID)
}

func (s *UserService) GetAllUsers(ctx context.Context) ([]models.User, error) {
	return s.userRepo.GetAll(ctx)
}

func (s *UserService) GetUsersByTenant(ctx context.Context, tenantID uuid.UUID) ([]models.User, error) {
	return s.userRepo.GetByTenant(ctx, tenantID)
}

func (s *UserService) ValidatePassword(user *models.User, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	return err == nil
}
