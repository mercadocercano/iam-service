-- Down: revierte la adición de slug a roles.
-- Los roles Cashier y Supervisor seedeados en el up NO se eliminan: son datos
-- que pueden estar referenciados en users.role_id de producción. Eliminarlos
-- rompería integridad referencial. Revertirlos requiere validación manual.
-- Los backfills de slug en roles de sistema quedan sin efecto al dropear la columna.

-- 1. Eliminar los índices únicos parciales (deben ir antes del DROP COLUMN).
DROP INDEX IF EXISTS idx_roles_slug_system;
DROP INDEX IF EXISTS idx_roles_slug_tenant;

-- 2. Eliminar la columna slug.
ALTER TABLE roles DROP COLUMN IF EXISTS slug;
