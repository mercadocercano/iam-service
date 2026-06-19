-- Down: revierte la adición de rate_limits a plans.
-- El seed de valores por tier (UPDATE) no necesita revertirse explícitamente:
-- desaparece al dropear la columna.

-- 1. Eliminar el índice GIN (Postgres lo dropea automáticamente con DROP COLUMN,
--    pero lo explicitamos para claridad y compatibilidad con versiones antiguas).
DROP INDEX IF EXISTS idx_plans_rate_limits;

-- 2. Eliminar la columna rate_limits.
ALTER TABLE plans DROP COLUMN IF EXISTS rate_limits;
