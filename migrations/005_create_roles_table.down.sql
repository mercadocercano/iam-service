-- Down: elimina la tabla roles y todos sus índices y triggers.
-- Nota: los roles de sistema insertados en el up se pierden con el DROP.
DROP TABLE IF EXISTS roles CASCADE;
