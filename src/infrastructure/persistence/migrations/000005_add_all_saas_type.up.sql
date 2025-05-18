-- Modificar el tipo ENUM saas_type para incluir 'ALL'
ALTER TYPE saas_type ADD VALUE 'ALL';

-- Crear el plan SAAS ADMIN
INSERT INTO plans (saas, name, description, monthly_price, yearly_price) VALUES
    ('ALL', 'SAAS ADMIN', 'Plan para administradores del sistema', 0, 0);

-- Crear el rol SUPERADMIN
INSERT INTO roles (saas, name, description) VALUES
    ('ALL', 'SUPERADMIN', 'Super administrador con acceso total');

-- Crear el tenant con saas_type ALL
WITH new_tenant AS (
    INSERT INTO tenants (saas, name, plan_id, email_user_key) 
    SELECT 'ALL', 'SAAS Admin Tenant', plans.id, 'admin@saasadmin.com'
    FROM plans 
    WHERE plans.name = 'SAAS ADMIN'
    RETURNING id
)
-- Crear el usuario admin
INSERT INTO users (tenant_id, role_id, email, password_hash, status)
SELECT 
    new_tenant.id,
    roles.id,
    'admin@saasadmin.com',
    '$2a$10$m9TStlKCnapXLFlhzwqXuuzlNRBamHcuhLGOuQEHMsxEG5GCoz37K',
    'ACTIVE'
FROM new_tenant, roles 
WHERE roles.name = 'SUPERADMIN' AND roles.saas = 'ALL';
