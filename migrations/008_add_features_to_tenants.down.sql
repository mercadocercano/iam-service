-- Down: elimina la columna features de tenants (y el índice GIN asociado, que
-- Postgres dropea automáticamente al hacer DROP COLUMN).
ALTER TABLE tenants DROP COLUMN IF EXISTS features;
