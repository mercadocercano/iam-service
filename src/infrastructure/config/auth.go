package config

import (
    "os"
    "time"
)

type AuthConfig struct {
    JWTSecret           string
    AccessTokenExpiry   time.Duration
    RefreshTokenExpiry  time.Duration
    GoogleClientID      string
}

func NewAuthConfig() AuthConfig {
    // En producción, estas variables deberían venir de variables de entorno
    return AuthConfig{
        JWTSecret:           getEnvOrDefault("JWT_SECRET", "your-secret-key"),
        AccessTokenExpiry:   time.Hour,      // 1 hora
        RefreshTokenExpiry:  time.Hour * 24, // 24 horas
        GoogleClientID:      getEnvOrDefault("GOOGLE_CLIENT_ID", ""),
    }
}

func getEnvOrDefault(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}
