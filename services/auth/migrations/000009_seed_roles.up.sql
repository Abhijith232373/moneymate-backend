INSERT INTO auth.roles (id, name, description) VALUES
    ('a0000000-0000-0000-0000-000000000001', 'user', 'Standard end user account'),
    ('a0000000-0000-0000-0000-000000000002', 'merchant', 'Merchant/business account'),
    ('a0000000-0000-0000-0000-000000000003', 'admin', 'Internal administrator account')
ON CONFLICT (name) DO NOTHING;