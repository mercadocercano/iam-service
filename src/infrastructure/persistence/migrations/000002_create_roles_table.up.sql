CREATE TABLE roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    saas saas_type NOT NULL,
    name VARCHAR(50) NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(saas, name)
);

-- Insert default roles
INSERT INTO roles (saas, name, description) VALUES
    ('CRM', 'SUPERADMIN', 'Super administrador con acceso total'),
    ('CRM', 'ADMIN', 'Administrador con acceso limitado'),
    ('CRM', 'MEMBER', 'Usuario regular con acceso básico'),
    ('ERP', 'SUPERADMIN', 'Super administrador con acceso total'),
    ('ERP', 'ADMIN', 'Administrador con acceso limitado'),
    ('ERP', 'MEMBER', 'Usuario regular con acceso básico'),
    ('ECOMMERCE', 'SUPERADMIN', 'Super administrador con acceso total'),
    ('ECOMMERCE', 'ADMIN', 'Administrador con acceso limitado'),
    ('ECOMMERCE', 'MEMBER', 'Usuario regular con acceso básico');
