-- Down no soportado: esta migración inserta datos de sistema (tenant system, tenant demo,
-- usuarios admin) con UUIDs fijos que están hardcodeados en otros servicios del ecosistema.
-- Eliminarlos rompe contratos inter-servicio. Revertir requiere restaurar desde backup.
-- Ver ADR-001 (go-shared).
DO $$ BEGIN
  RAISE EXCEPTION 'irreversible migration: down no soportado (ver ADR-001)';
END $$;
