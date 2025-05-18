-- Eliminar tabla refresh_tokens
DROP TABLE IF EXISTS refresh_tokens;

-- Eliminar columnas de users
ALTER TABLE users
    DROP COLUMN IF EXISTS provider,
    DROP COLUMN IF EXISTS federated_id,
    ALTER COLUMN password_hash SET NOT NULL;

-- Eliminar tipo auth_provider
DROP TYPE IF EXISTS auth_provider;
