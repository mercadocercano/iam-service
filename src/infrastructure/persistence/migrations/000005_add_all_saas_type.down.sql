-- Eliminar el usuario admin
DELETE FROM users WHERE username = 'adminall';

-- Eliminar el tenant
DELETE FROM tenants WHERE saas = 'ALL';

-- Eliminar el rol SUPERADMIN
DELETE FROM roles WHERE name = 'SUPERADMIN' AND saas = 'ALL';

-- Eliminar el plan SAAS ADMIN
DELETE FROM plans WHERE name = 'SAAS ADMIN';

-- No se puede eliminar un valor de un ENUM en PostgreSQL directamente.
-- Se tendría que recrear el tipo sin ese valor si fuera necesario.
