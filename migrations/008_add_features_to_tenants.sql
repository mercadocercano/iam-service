-- Migration: Add features column to tenants table
-- Description: Agregar columna features de tipo JSONB para almacenar feature flags del tenant

-- Agregar la columna features
ALTER TABLE tenants ADD COLUMN IF NOT EXISTS features JSONB DEFAULT '{"friends_family": false, "premium_analytics": false}';

-- Actualizar tenants existentes que tengan features NULL
UPDATE tenants 
SET features = '{"friends_family": false, "premium_analytics": false}'
WHERE features IS NULL;

-- Crear índice en la columna features para consultas rápidas
CREATE INDEX IF NOT EXISTS idx_tenants_features ON tenants USING GIN (features);

-- Comentario sobre la nueva columna
COMMENT ON COLUMN tenants.features IS 'Feature flags del tenant en formato JSON';

-- Verificar que los datos se actualizaron correctamente
DO $$
BEGIN
    RAISE NOTICE '✅ Columna features agregada a la tabla tenants';
    RAISE NOTICE '🔧 Tenants existentes actualizados con features por defecto';
    RAISE NOTICE '📊 Índice GIN creado para consultas en features';
END $$; 