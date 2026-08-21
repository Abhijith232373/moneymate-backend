CREATE TABLE audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    admin_id UUID NOT NULL,
    admin_name VARCHAR(255) NOT NULL,
    admin_role VARCHAR(100) NOT NULL,
    module VARCHAR(100) NOT NULL,
    action VARCHAR(255) NOT NULL,
    changes JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
