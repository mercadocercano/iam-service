-- Migration: Seed system data
-- Description: Datos iniciales del sistema (usuario admin)

-- Crear tenant de sistema si la tabla existe y no está ya creado
DO $$
BEGIN
    -- Verificar si la tabla tenants existe
    IF EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'tenants') THEN
        -- Insertar tenant de sistema si no existe
        INSERT INTO tenants (id, name, slug, description, type, status, owner_id, max_users) 
        SELECT 
            '123e4567-e89b-12d3-a456-426614174003',
            'System',
            'system',
            'Tenant de sistema para administración',
            'ENTERPRISE',
            'ACTIVE',
            '123e4567-e89b-12d3-a456-426614174004',
            999
        WHERE NOT EXISTS (SELECT 1 FROM tenants WHERE slug = 'system');
        
        -- Insertar tenant de demo/testing (hardcodeado en backoffice y otros servicios)
        INSERT INTO tenants (id, name, slug, description, type, status, owner_id, max_users) 
        SELECT 
            '9a4c3eb9-2471-4688-bfc8-973e5b3e4ce8',
            'Demo Company',
            'demo-company',
            'Tenant de demostración para testing',
            'BUSINESS',
            'ACTIVE',
            '123e4567-e89b-12d3-a456-426614174005',
            100
        WHERE NOT EXISTS (SELECT 1 FROM tenants WHERE id = '9a4c3eb9-2471-4688-bfc8-973e5b3e4ce8');
        
        RAISE NOTICE 'Sistema: Tenants verificados/creados';
    ELSE
        RAISE NOTICE 'Sistema: Tabla tenants no existe, continuando sin tenant';
    END IF;
END $$;

-- Crear usuario admin del sistema si no existe
-- Password: 123456 (hash bcrypt: $2a$10$yMecfP7H7mT0m9VHNnMFIezB06ihuIqUVpD12sa34UHhSfCRIQdje)
DO $$
DECLARE
    system_role_id UUID;
    tenant_admin_role_id UUID;
    system_tenant_id UUID := '123e4567-e89b-12d3-a456-426614174003';
    demo_tenant_id UUID := '9a4c3eb9-2471-4688-bfc8-973e5b3e4ce8';
BEGIN
    -- Obtener los IDs de los roles
    SELECT id INTO system_role_id 
    FROM roles 
    WHERE name = 'System Administrator' AND type = 'SYSTEM_ADMIN'
    LIMIT 1;
    
    SELECT id INTO tenant_admin_role_id 
    FROM roles 
    WHERE name = 'Tenant Administrator' AND type = 'TENANT_ADMIN'
    LIMIT 1;
    
    IF system_role_id IS NULL THEN
        RAISE EXCEPTION 'Sistema: Rol System Administrator no encontrado';
    END IF;
    
    IF tenant_admin_role_id IS NULL THEN
        RAISE EXCEPTION 'Sistema: Rol Tenant Administrator no encontrado';
    END IF;
    
    -- Si no existe tabla tenants, usar un tenant_id fijo
    IF NOT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'tenants') THEN
        system_tenant_id := '00000000-0000-0000-0000-000000000000';
        demo_tenant_id := '00000000-0000-0000-0000-000000000001';
    END IF;
    
    -- Insertar usuario admin de sistema si no existe
    INSERT INTO users (id, email, password_hash, tenant_id, role_id, status, provider) 
    VALUES (
        '123e4567-e89b-12d3-a456-426614174004',
        'admin@saasadmin.com',
        '$2a$10$yMecfP7H7mT0m9VHNnMFIezB06ihuIqUVpD12sa34UHhSfCRIQdje',
        system_tenant_id,
        system_role_id,
        'ACTIVE',
        'LOCAL'
    )
    ON CONFLICT (email, tenant_id) DO NOTHING;
    
    -- Insertar usuario admin de demo company si no existe  
    INSERT INTO users (id, email, password_hash, tenant_id, role_id, status, provider) 
    VALUES (
        '123e4567-e89b-12d3-a456-426614174005',
        'admin@democompany.com',
        '$2a$10$yMecfP7H7mT0m9VHNnMFIezB06ihuIqUVpD12sa34UHhSfCRIQdje',
        demo_tenant_id,
        tenant_admin_role_id,
        'ACTIVE',
        'LOCAL'
    )
    ON CONFLICT (email, tenant_id) DO NOTHING;
    
    -- También crear el usuario admin@saasadmin.com en el tenant de demo para compatibilidad
    INSERT INTO users (id, email, password_hash, tenant_id, role_id, status, provider) 
    VALUES (
        '123e4567-e89b-12d3-a456-426614174006',
        'admin@saasadmin.com',
        '$2a$10$yMecfP7H7mT0m9VHNnMFIezB06ihuIqUVpD12sa34UHhSfCRIQdje',
        demo_tenant_id,
        tenant_admin_role_id,
        'ACTIVE',
        'LOCAL'
    )
    ON CONFLICT (email, tenant_id) DO NOTHING;
    
    -- Verificar que los usuarios se crearon
    IF EXISTS (SELECT 1 FROM users WHERE email = 'admin@saasadmin.com') THEN
        RAISE NOTICE 'Sistema: Usuarios admin creados/verificados exitosamente';
    ELSE
        RAISE EXCEPTION 'Sistema: Error creando usuarios admin';
    END IF;
END $$;

-- Actualizar los owner_id de los tenants si existen
DO $$
BEGIN
    IF EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'tenants') THEN
        UPDATE tenants 
        SET owner_id = '123e4567-e89b-12d3-a456-426614174004' 
        WHERE slug = 'system' AND owner_id != '123e4567-e89b-12d3-a456-426614174004';
        
        UPDATE tenants 
        SET owner_id = '123e4567-e89b-12d3-a456-426614174005' 
        WHERE id = '9a4c3eb9-2471-4688-bfc8-973e5b3e4ce8' AND owner_id != '123e4567-e89b-12d3-a456-426614174005';
    END IF;
END $$;

-- Mensaje final
DO $$
BEGIN
    RAISE NOTICE '🎉 Sistema: Datos de seed creados exitosamente';
    RAISE NOTICE '👤 Sistema: Usuario admin sistema: admin@saasadmin.com / 123456';
    RAISE NOTICE '👤 Sistema: Usuario admin demo: admin@democompany.com / 123456';
    RAISE NOTICE '👤 Sistema: Usuario admin (compatibilidad): admin@saasadmin.com / 123456 en tenant demo';
    RAISE NOTICE '🔑 Sistema: Hash de contraseña: $2a$10$yMecfP7H7mT0m9VHNnMFIezB06ihuIqUVpD12sa34UHhSfCRIQdje';
    RAISE NOTICE '🏢 Sistema: Tenant ID demo (hardcoded): 9a4c3eb9-2471-4688-bfc8-973e5b3e4ce8';
END $$; 