CREATE TABLE tenants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    saas saas_type NOT NULL,
    name VARCHAR(100) NOT NULL,
    plan_id UUID NOT NULL REFERENCES plans(id),
    email_user_key VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(saas, name),
    UNIQUE(email_user_key)
);
