package models

import (
    "time"
    "github.com/google/uuid"
)

type Role struct {
    ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
    Saas        SaasType  `gorm:"type:saas_type;not null"`
    Name        string    `gorm:"size:50;not null"`
    Description string
    CreatedAt   time.Time
    UpdatedAt   time.Time
}
