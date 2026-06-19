-- Down: elimina la función trigger update_updated_at_column y todos los triggers
-- que la usan (CASCADE los dropea automáticamente).
DROP FUNCTION IF EXISTS update_updated_at_column() CASCADE;
