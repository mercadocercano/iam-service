-- Migration: Add slug to roles + seed cashier/supervisor (RBAC en JWT)
-- Description: El slug es el identificador de autorización estable que viaja en el
-- claim `roles` del JWT (independiente del Name visible y del RoleType grueso).
-- Aditiva y retrocompatible: la columna es nullable y el seed es idempotente.

-- 1. Columna slug (nullable para no romper filas existentes)
ALTER TABLE roles ADD COLUMN IF NOT EXISTS slug VARCHAR(50);

-- 1.5. Dedup defensivo de roles de SISTEMA duplicados (idempotente / no-op si no hay).
--    Causa raíz: en algunos entornos el init-db.sh seedeó los roles de sistema DOS veces
--    (2 filas por type con tenant_id NULL), lo que rompe el índice único parcial del paso 3
--    (idx_roles_slug_system). Resolución: conservar el más viejo por type, repuntar
--    users.role_id de los duplicados al conservado, y borrar los duplicados. Se archivan
--    los borrados en roles_dup_archive para auditoría/rollback. Acotado a los 4 tipos de
--    sistema estándar (no toca CUSTOM como cashier/supervisor). Set-based, atómico con la
--    migración (golang-migrate envuelve en transacción): si algo falla, revierte todo.
CREATE TABLE IF NOT EXISTS roles_dup_archive AS SELECT * FROM roles WHERE false;

INSERT INTO roles_dup_archive
SELECT r.* FROM roles r
JOIN (
  SELECT id, row_number() OVER (PARTITION BY type ORDER BY created_at ASC, id ASC) AS rn
  FROM roles
  WHERE tenant_id IS NULL AND type IN ('SYSTEM_ADMIN','TENANT_ADMIN','USER','READ_ONLY')
) d ON r.id = d.id AND d.rn > 1;

UPDATE users u SET role_id = k.keeper_id
FROM (
  SELECT id,
         first_value(id) OVER (PARTITION BY type ORDER BY created_at ASC, id ASC) AS keeper_id,
         row_number()    OVER (PARTITION BY type ORDER BY created_at ASC, id ASC) AS rn
  FROM roles
  WHERE tenant_id IS NULL AND type IN ('SYSTEM_ADMIN','TENANT_ADMIN','USER','READ_ONLY')
) k
WHERE u.role_id = k.id AND k.rn > 1;

DELETE FROM roles r
USING (
  SELECT id, row_number() OVER (PARTITION BY type ORDER BY created_at ASC, id ASC) AS rn
  FROM roles
  WHERE tenant_id IS NULL AND type IN ('SYSTEM_ADMIN','TENANT_ADMIN','USER','READ_ONLY')
) d
WHERE r.id = d.id AND d.rn > 1;

-- 2. Backfill de los roles de sistema existentes
UPDATE roles SET slug = 'system_admin' WHERE type = 'SYSTEM_ADMIN' AND slug IS NULL;
UPDATE roles SET slug = 'tenant_admin' WHERE type = 'TENANT_ADMIN' AND slug IS NULL;
UPDATE roles SET slug = 'user'         WHERE type = 'USER'         AND slug IS NULL;
UPDATE roles SET slug = 'read_only'    WHERE type = 'READ_ONLY'    AND slug IS NULL;

-- 3. Unicidad de slug por tenant.
--    OJO (condición @architect A1 / @security O4): en un UNIQUE normal Postgres trata
--    dos NULL como distintos, así que dos roles de SISTEMA (tenant_id IS NULL) con el
--    mismo slug podrían coexistir. Se usan DOS índices parciales para garantizar
--    unicidad real en ambos lados:
CREATE UNIQUE INDEX IF NOT EXISTS idx_roles_slug_system
    ON roles(slug) WHERE tenant_id IS NULL;
CREATE UNIQUE INDEX IF NOT EXISTS idx_roles_slug_tenant
    ON roles(slug, tenant_id) WHERE tenant_id IS NOT NULL;

-- 4. Seed de roles de caja como roles de SISTEMA (tenant_id NULL), reutilizables por
--    todos los tenants. Idempotente vía el índice parcial de sistema sobre slug.
--    El permiso `sales:cash_session:approve_review` SOLO lo tiene supervisor:
--    materializa la separación de funciones (un cajero no aprueba su propio descuadre).
INSERT INTO roles (name, description, type, tenant_id, slug, permissions, is_active)
VALUES
  ('Cashier', 'Operador de caja POS', 'CUSTOM', NULL, 'cashier',
   ARRAY['sales:pos:sell',
         'sales:cash_session:open',
         'sales:cash_session:close',
         'sales:cash_session:read',
         'sales:cash_movement:create'],
   true)
ON CONFLICT DO NOTHING;

INSERT INTO roles (name, description, type, tenant_id, slug, permissions, is_active)
VALUES
  ('Supervisor', 'Supervisor de caja y arqueos', 'CUSTOM', NULL, 'supervisor',
   ARRAY['sales:pos:sell',
         'sales:cash_session:open',
         'sales:cash_session:close',
         'sales:cash_session:read',
         'sales:cash_movement:create',
         'sales:cash_session:approve_review'],
   true)
ON CONFLICT DO NOTHING;
