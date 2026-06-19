-- Migration: Create refresh_tokens table
-- Description: Tabla para almacenar refresh tokens del sistema de autenticación

CREATE TABLE IF NOT EXISTS refresh_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    token VARCHAR(255) NOT NULL UNIQUE,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    -- Foreign key constraint
    CONSTRAINT fk_refresh_tokens_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Indexes para optimizar consultas
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user_id ON refresh_tokens(user_id);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_token ON refresh_tokens(token);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_expires_at ON refresh_tokens(expires_at);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_created_at ON refresh_tokens(created_at);

-- Comentarios sobre la tabla
COMMENT ON TABLE refresh_tokens IS 'Tabla para almacenar refresh tokens del sistema de autenticación JWT';
COMMENT ON COLUMN refresh_tokens.id IS 'Identificador único del refresh token';
COMMENT ON COLUMN refresh_tokens.user_id IS 'ID del usuario propietario del token';
COMMENT ON COLUMN refresh_tokens.token IS 'Valor del refresh token (único)';
COMMENT ON COLUMN refresh_tokens.expires_at IS 'Fecha y hora de expiración del token';
COMMENT ON COLUMN refresh_tokens.created_at IS 'Fecha y hora de creación del token'; 