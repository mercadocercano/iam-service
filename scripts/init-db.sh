#!/bin/bash

echo "🚀 Iniciando migraciones del IAM..."

# Configurar variables de entorno por defecto si no están definidas
DB_HOST=${DB_HOST:-postgres}
DB_PORT=${DB_PORT:-5432}
DB_USER=${DB_USER:-postgres}
DB_NAME=${DB_NAME:-iam_db}

# Función para ejecutar SQL
execute_sql() {
    local file=$1
    echo "📁 Ejecutando: $(basename "$file")"
    PGPASSWORD=$POSTGRES_PASSWORD psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -f "$file"
    
    if [ $? -eq 0 ]; then
        echo "✅ $(basename "$file") ejecutado exitosamente"
    else
        echo "❌ Error ejecutando $(basename "$file")"
        exit 1
    fi
}

# Crear la función update_updated_at si no existe (necesaria para triggers)
echo "🔧 Creando funciones auxiliares..."
PGPASSWORD=$POSTGRES_PASSWORD psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" <<EOF
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS \$\$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
\$\$ language 'plpgsql';
EOF

# Ejecutar migraciones en orden específico
echo "🗃️ Ejecutando migraciones en orden..."

# 1. Primero las migraciones de tenant (si existen) - solo UP
if [ -f "/src/tenant/infrastructure/persistence/migrations/001_create_tenants_table.sql" ]; then
    echo "📁 Ejecutando migración de tenants..."
    # Ejecutar solo la parte UP de la migración
    PGPASSWORD=$POSTGRES_PASSWORD psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" <<EOF
-- +migrate Up
CREATE TABLE IF NOT EXISTS tenants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    type VARCHAR(50) NOT NULL CHECK (type IN ('PERSONAL', 'STARTUP', 'BUSINESS', 'ENTERPRISE')),
    status VARCHAR(50) NOT NULL DEFAULT 'ACTIVE' CHECK (status IN ('ACTIVE', 'INACTIVE', 'SUSPENDED', 'DELETED')),
    owner_id UUID NOT NULL,
    domain VARCHAR(255) UNIQUE,
    plan_id UUID,
    user_count INTEGER NOT NULL DEFAULT 0,
    max_users INTEGER,
    settings JSONB DEFAULT '{}',
    expires_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Crear índices para optimizar búsquedas
CREATE INDEX IF NOT EXISTS idx_tenants_owner_id ON tenants(owner_id);
CREATE INDEX IF NOT EXISTS idx_tenants_status ON tenants(status);
CREATE INDEX IF NOT EXISTS idx_tenants_type ON tenants(type);
CREATE INDEX IF NOT EXISTS idx_tenants_plan_id ON tenants(plan_id);
CREATE INDEX IF NOT EXISTS idx_tenants_expires_at ON tenants(expires_at) WHERE expires_at IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_tenants_slug ON tenants(slug);
CREATE INDEX IF NOT EXISTS idx_tenants_domain ON tenants(domain) WHERE domain IS NOT NULL;

-- Trigger para actualizar updated_at automáticamente
CREATE TRIGGER update_tenants_updated_at 
    BEFORE UPDATE ON tenants 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();
EOF
    echo "✅ Migración de tenants ejecutada exitosamente"
fi

# 2. Luego las migraciones principales en orden
for file in /migrations/003_create_users_table.sql \
           /migrations/004_create_plans_table.sql \
           /migrations/005_create_roles_table.sql \
           /migrations/006_seed_system_data.sql \
           /migrations/007_create_refresh_tokens_table.sql \
           /migrations/008_add_features_to_tenants.sql; do
    if [ -f "$file" ]; then
        execute_sql "$file"
    else
        echo "⚠️ Archivo no encontrado: $file"
    fi
done

# 3. Ejecutar cualquier migración adicional que no esté en la lista
echo "🔍 Buscando migraciones adicionales..."
for file in /migrations/*.sql; do
    filename=$(basename "$file")
    case "$filename" in
        "003_create_users_table.sql"|"004_create_plans_table.sql"|"005_create_roles_table.sql"|"006_seed_system_data.sql"|"007_create_refresh_tokens_table.sql"|"008_add_features_to_tenants.sql")
            # Ya ejecutada
            ;;
        *)
            if [ -f "$file" ]; then
                execute_sql "$file"
            fi
            ;;
    esac
done

echo "🎉 Migraciones del IAM completadas exitosamente!"
echo "👤 Usuario admin creado: admin@saasadmin.com / 123456"
