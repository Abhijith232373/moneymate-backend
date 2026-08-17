INSERT INTO auth.permissions (id, name, description) VALUES
    (gen_random_uuid(), 'dashboard.read', 'View analytics and system overview'),
    (gen_random_uuid(), 'users.create', 'Create user accounts'),
    (gen_random_uuid(), 'users.read', 'Read user accounts'),
    (gen_random_uuid(), 'users.update', 'Update user accounts'),
    (gen_random_uuid(), 'users.delete', 'Delete user accounts'),
    (gen_random_uuid(), 'store.create', 'Create merchant stores'),
    (gen_random_uuid(), 'store.read', 'Read merchant stores'),
    (gen_random_uuid(), 'store.update', 'Update merchant stores'),
    (gen_random_uuid(), 'store.delete', 'Delete merchant stores'),
    (gen_random_uuid(), 'support.create', 'Create support tickets'),
    (gen_random_uuid(), 'support.read', 'Read support tickets'),
    (gen_random_uuid(), 'support.update', 'Update support tickets'),
    (gen_random_uuid(), 'support.delete', 'Delete support tickets'),
    (gen_random_uuid(), 'settings.read', 'Read system settings'),
    (gen_random_uuid(), 'settings.update', 'Update system settings'),
    (gen_random_uuid(), 'audit.read', 'Read audit logs')
ON CONFLICT (name) DO NOTHING;
