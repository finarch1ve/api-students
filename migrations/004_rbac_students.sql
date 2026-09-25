CREATE TABLE IF NOT EXISTS roles (
    name VARCHAR(20) PRIMARY KEY,
    description VARCHAR(150) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
INSERT INTO roles (name, description) VALUES
    ('admin', 'Akses penuh terhadap seluruh data'),
    ('staff', 'Boleh melihat dan menambah data'),
    ('user', 'Hanya boleh mengelola datanya sendiri')
ON CONFLICT (name) DO NOTHING;

CREATE TABLE IF NOT EXISTS permissions (
    name VARCHAR(50) PRIMARY KEY,
    description VARCHAR(150) NOT NULL
);
INSERT INTO permissions (name, description) VALUES
    ('student:list', 'Melihat daftar seluruh mahasiswa'),
    ('student:read:any', 'Melihat data mahasiswa mana pun'),
    ('student:create', 'Menambahkan data mahasiswa baru'),
    ('student:update:any', 'Mengubah data mahasiswa mana pun'),
    ('student:delete', 'Menghapus data mahasiswa')
ON CONFLICT (name) DO NOTHING;

CREATE TABLE IF NOT EXISTS role_permissions (
    role_name VARCHAR(20) NOT NULL REFERENCES roles(name) ON DELETE CASCADE,
    permission_name VARCHAR(50) NOT NULL REFERENCES permissions(name) ON DELETE CASCADE,
    PRIMARY KEY (role_name, permission_name)
);
INSERT INTO role_permissions (role_name, permission_name) VALUES
    ('admin','student:list'),
    ('admin','student:read:any'),
    ('admin','student:create'),
    ('admin','student:update:any'),
    ('admin','student:delete'),
    ('staff','student:list'),
    ('staff','student:read:any'),
    ('staff','student:create')
ON CONFLICT DO NOTHING;

-- owner_id: siapa yang "punya" data mahasiswa ini
ALTER TABLE students ADD COLUMN IF NOT EXISTS owner_id INTEGER REFERENCES users(id);
UPDATE students SET owner_id = (SELECT id FROM users ORDER BY id LIMIT 1) WHERE owner_id IS NULL;
ALTER TABLE students ALTER COLUMN owner_id SET NOT NULL;