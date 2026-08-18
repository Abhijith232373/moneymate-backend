DELETE FROM auth.permissions WHERE name IN (
    'dashboard.read',
    'users.create', 'users.read', 'users.update', 'users.delete',
    'store.create', 'store.read', 'store.update', 'store.delete',
    'support.create', 'support.read', 'support.update', 'support.delete',
    'settings.read', 'settings.update',
    'audit.read'
);
