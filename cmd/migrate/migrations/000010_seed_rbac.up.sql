INSERT INTO auth.roles (name, description) VALUES
    ('admin', 'Full access to all resources'),
    ('user', 'Standard user access')
ON CONFLICT DO NOTHING;

INSERT INTO auth.permissions (name, resource, action, description) VALUES
    ('categories.create', 'categories', 'create', 'Create categories'),
    ('categories.read', 'categories', 'read', 'View categories'),
    ('categories.update', 'categories', 'update', 'Update categories'),
    ('categories.delete', 'categories', 'delete', 'Delete categories'),
    ('groups.create', 'category_groups', 'create', 'Create category groups'),
    ('groups.read', 'category_groups', 'read', 'View category groups'),
    ('groups.update', 'category_groups', 'update', 'Update category groups'),
    ('groups.delete', 'category_groups', 'delete', 'Delete category groups'),
    ('months.read', 'months', 'read', 'View months'),
    ('months.manage', 'months', 'manage', 'Create, close, and rollforward months'),
    ('months.reallocate', 'months', 'update', 'Reallocate funds between categories'),
    ('transactions.create', 'transactions', 'create', 'Create transactions'),
    ('transactions.read', 'transactions', 'read', 'View transactions'),
    ('transactions.update', 'transactions', 'update', 'Update transactions'),
    ('transactions.delete', 'transactions', 'delete', 'Delete transactions'),
    ('accounts.create', 'accounts', 'create', 'Create accounts'),
    ('accounts.read', 'accounts', 'read', 'View accounts'),
    ('accounts.update', 'accounts', 'update', 'Update accounts'),
    ('accounts.delete', 'accounts', 'delete', 'Delete accounts'),
    ('transfer.create', 'transfers', 'create', 'Create transfers between accounts')
ON CONFLICT DO NOTHING;

-- Assign all permissions to admin role
INSERT INTO auth.role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM auth.roles r
CROSS JOIN auth.permissions p
WHERE r.name = 'admin'
ON CONFLICT DO NOTHING;

-- Assign basic permissions to user role
INSERT INTO auth.role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM auth.roles r
CROSS JOIN auth.permissions p
WHERE r.name = 'user'
  AND p.name IN (
      'categories.create', 'categories.read', 'categories.update', 'categories.delete',
      'groups.create', 'groups.read', 'groups.update', 'groups.delete',
      'months.read',
      'transactions.create', 'transactions.read', 'transactions.update', 'transactions.delete',
      'accounts.create', 'accounts.read', 'accounts.update', 'accounts.delete',
      'transfer.create'
  )
ON CONFLICT DO NOTHING;
