package models

import (
	"time"

	"github.com/google/uuid"
)

// UserUpdate contiene solo los campos que se pueden actualizar
type UserUpdate struct {
	ID     uuid.UUID  `json:"id" validate:"required"`
	Email  string     `json:"email,omitempty" validate:"omitempty,email"`
	RoleID uuid.UUID  `json:"role_id,omitempty"`
	Status UserStatus `json:"status,omitempty" validate:"omitempty,oneof=PENDING ACTIVE INACTIVE BLOCKED"`
}

type User struct {
	ID           uuid.UUID    `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Email        string       `gorm:"size:255;not null"`
	PasswordHash string       `gorm:"size:255"`
	TenantID     uuid.UUID    `gorm:"type:uuid;not null"`
	RoleID       uuid.UUID    `gorm:"type:uuid;not null"`
	Status       UserStatus   `gorm:"type:user_status;not null;default:'PENDING'"`
	Provider     AuthProvider `gorm:"type:varchar(20);not null;default:'LOCAL'"`
	FederatedID  string       `gorm:"size:255"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Tenant       Tenant `gorm:"foreignKey:TenantID"`
	Role         Role   `gorm:"foreignKey:RoleID"`
}
