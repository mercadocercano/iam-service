-- Crear tipo auth_provider
CREATE TYPE auth_provider AS ENUM ('LOCAL', 'GOOGLE');

-- Agregar columnas a users
ALTER TABLE users
    ADD COLUMN provider auth_provider NOT NULL DEFAULT 'LOCAL',
    ADD COLUMN federated_id VARCHAR(255),
    ALTER COLUMN password_hash DROP NOT NULL;

-- Crear tabla refresh_tokens
CREATE TABLE refresh_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token VARCHAR(255) NOT NULL UNIQUE,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_refresh_tokens_user FOREIGN KEY (user_id) REFERENCES users(id)
);

-- Crear índices
CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens(user_id);
CREATE INDEX idx_refresh_tokens_token ON refresh_tokens(token);
