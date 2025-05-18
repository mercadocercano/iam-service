package models

import (
    "time"
    "github.com/google/uuid"
)

type Tenant struct {
    ID           uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
    Saas         SaasType  `gorm:"type:saas_type;not null"`
    Name         string    `gorm:"size:100;not null"`
    PlanID       uuid.UUID `gorm:"type:uuid;not null"`
    EmailUserKey string    `gorm:"size:255;not null;unique"`
    CreatedAt    time.Time
    UpdatedAt    time.Time
    Plan         Plan      `gorm:"foreignKey:PlanID"`
    Users        []User    `gorm:"foreignKey:TenantID"`
}
