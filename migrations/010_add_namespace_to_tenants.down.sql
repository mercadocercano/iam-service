-- Down: revierte la adición de namespace a tenants.
-- El backfill de datos (UPDATE SET namespace = 'mc') no se revierte porque
-- la columna desaparece con DROP COLUMN.

-- 1. Restaurar la constraint UNIQUE original sobre slug (la que existía antes del up).
ALTER TABLE tenants DROP CONSTRAINT IF EXISTS tenants_namespace_slug_unique;
ALTER TABLE tenants ADD CONSTRAINT tenants_slug_key UNIQUE (slug);

-- 2. Eliminar el índice de namespace.
DROP INDEX IF EXISTS idx_tenants_namespace;

-- 3. Eliminar la columna namespace.
ALTER TABLE tenants DROP COLUMN IF EXISTS namespace;
