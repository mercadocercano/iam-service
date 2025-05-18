package models

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AuthProvider string

const (
	LocalAuth  AuthProvider = "LOCAL"
	GoogleAuth AuthProvider = "GOOGLE"
)

type LoginRequest struct {
	Email       string       `json:"email" binding:"required,email"`
	Password    string       `json:"password,omitempty"`
	Provider    AuthProvider `json:"provider" binding:"required"`
	GoogleToken string       `json:"google_token,omitempty"`
	TenantID    *uuid.UUID   `json:"tenant_id,omitempty" binding:"-"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
	User         User   `json:"user"`
}

type RefreshToken struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID    uuid.UUID `gorm:"type:uuid;not null"`
	Token     string    `gorm:"size:255;not null;unique"`
	ExpiresAt time.Time `gorm:"not null"`
	CreatedAt time.Time
	User      User `gorm:"foreignKey:UserID"`
}

type TokenClaims struct {
	UserID    uuid.UUID `json:"user_id"`
	Email     string    `json:"email"`
	TenantID  uuid.UUID `json:"tenant_id"`
	RoleID    uuid.UUID `json:"role_id"`
	ExpiresAt int64     `json:"exp"`
}

func (c TokenClaims) GetExpirationTime() (*jwt.NumericDate, error) {
	if c.ExpiresAt == 0 {
		return nil, nil
	}
	return jwt.NewNumericDate(time.Unix(c.ExpiresAt, 0)), nil
}

func (c TokenClaims) GetIssuedAt() (*jwt.NumericDate, error) {
	return nil, nil
}

func (c TokenClaims) GetNotBefore() (*jwt.NumericDate, error) {
	return nil, nil
}

func (c TokenClaims) GetIssuer() (string, error) {
	return "", nil
}

func (c TokenClaims) GetSubject() (string, error) {
	return c.Email, nil
}

func (c TokenClaims) GetAudience() (jwt.ClaimStrings, error) {
	return nil, nil
}

func (c TokenClaims) GetID() (string, error) {
	return c.UserID.String(), nil
}
