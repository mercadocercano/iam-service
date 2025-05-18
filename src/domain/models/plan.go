package models

import (
    "time"
    "github.com/google/uuid"
    "github.com/lib/pq"
)

type Plan struct {
    ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
    Saas        SaasType  `gorm:"type:saas_type;not null"`
    Name        string    `gorm:"size:100;not null"`
    Description string
    Features    *pq.StringArray `gorm:"type:text[]"`
    MonthlyPrice float64
    YearlyPrice  float64
    CreatedAt   time.Time
    UpdatedAt   time.Time
}
