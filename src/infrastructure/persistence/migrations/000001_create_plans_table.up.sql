CREATE TYPE saas_type AS ENUM ('CRM', 'ERP', 'ECOMMERCE');

CREATE TABLE plans (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    saas saas_type NOT NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    features TEXT[] DEFAULT ARRAY[]::text[],
    monthly_price DECIMAL(10,2) NOT NULL,
    yearly_price DECIMAL(10,2) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
