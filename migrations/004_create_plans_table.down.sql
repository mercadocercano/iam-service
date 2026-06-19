-- Down: elimina la tabla plans y todos sus índices y triggers.
-- Nota: los planes por defecto insertados en el up se pierden con el DROP.
DROP TABLE IF EXISTS plans CASCADE;
