-- ============================================================
-- Migration 000011: Restructure Owner > Venue > POS Terminal
-- ============================================================
-- This is an atomic migration that:
-- 1. Creates owners, venues, pos_terminals tables
-- 2. Migrates data from restaurants to venues
-- 3. Adds venue_id to ingredients, products, categories, tables
-- 4. Updates users, orders, payments (restaurant_id -> venue_id)
-- 5. Drops waitress, table_waitress, restaurants

BEGIN;

-- =========================================
-- 1. Create owners table
-- =========================================
CREATE TABLE owners (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_owners_email ON owners(email);

-- =========================================
-- 2. Create venues table
-- =========================================
CREATE TABLE venues (
    id SERIAL PRIMARY KEY,
    owner_id INTEGER NOT NULL REFERENCES owners(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    address TEXT,
    phone VARCHAR(50),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_venues_owner_id ON venues(owner_id);

-- =========================================
-- 3. Create pos_terminals table
-- =========================================
CREATE TABLE pos_terminals (
    id SERIAL PRIMARY KEY,
    venue_id INTEGER NOT NULL REFERENCES venues(id) ON DELETE CASCADE,
    terminal_name VARCHAR(100) NOT NULL,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_pos_terminals_venue_id ON pos_terminals(venue_id);

-- =========================================
-- 4. Migrate existing data: restaurants -> owners + venues
-- =========================================

-- Create a default owner from the first PROPIETARIO user (if exists)
INSERT INTO owners (name, email, password_hash)
SELECT u.name, u.email, u.password_hash
FROM users u
WHERE u.role_id = 1
ORDER BY u.id
LIMIT 1
ON CONFLICT (email) DO NOTHING;

-- If no PROPIETARIO user exists, create a placeholder owner
INSERT INTO owners (name, email, password_hash)
SELECT 'Default Owner', 'owner@default.com', '$2a$10$placeholder'
WHERE NOT EXISTS (SELECT 1 FROM owners);

-- Create venues from existing restaurants
INSERT INTO venues (owner_id, name, address, phone, is_active, created_at, updated_at)
SELECT
    (SELECT id FROM owners ORDER BY id LIMIT 1),
    r.name,
    COALESCE(r.address, ''),
    COALESCE(r.phone, ''),
    r.is_active,
    r.created_at,
    r.updated_at
FROM restaurants r;

-- If no restaurants exist, create a default venue
INSERT INTO venues (owner_id, name)
SELECT (SELECT id FROM owners ORDER BY id LIMIT 1), 'Default Venue'
WHERE NOT EXISTS (SELECT 1 FROM venues);

-- Create a default POS terminal for each venue
INSERT INTO pos_terminals (venue_id, terminal_name)
SELECT v.id, 'Terminal Principal'
FROM venues v;

-- =========================================
-- 5. Update users: restaurant_id -> venue_id
-- =========================================
ALTER TABLE users ADD COLUMN venue_id INTEGER;

-- Map users to venues based on their restaurant_id
UPDATE users u SET venue_id = v.id
FROM (
    SELECT v.id, r_id
    FROM venues v
    CROSS JOIN LATERAL (
        SELECT r.id AS r_id FROM restaurants r
        ORDER BY r.id
        LIMIT 1
    ) sub
) v
WHERE u.restaurant_id = v.r_id;

-- For any users that didn't get mapped, assign to first venue
UPDATE users SET venue_id = (SELECT id FROM venues ORDER BY id LIMIT 1)
WHERE venue_id IS NULL;

ALTER TABLE users ALTER COLUMN venue_id SET NOT NULL;
ALTER TABLE users ADD CONSTRAINT fk_users_venue FOREIGN KEY (venue_id) REFERENCES venues(id) ON DELETE CASCADE;

-- Drop old restaurant FK and column
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_restaurant_id_fkey;
ALTER TABLE users DROP COLUMN restaurant_id;

DROP INDEX IF EXISTS idx_users_restaurant_id;
CREATE INDEX idx_users_venue_id ON users(venue_id);

-- =========================================
-- 6. Add venue_id to ingredients
-- =========================================
ALTER TABLE ingredients ADD COLUMN venue_id INTEGER;
UPDATE ingredients SET venue_id = (SELECT id FROM venues ORDER BY id LIMIT 1);
ALTER TABLE ingredients ALTER COLUMN venue_id SET NOT NULL;
ALTER TABLE ingredients ADD CONSTRAINT fk_ingredients_venue FOREIGN KEY (venue_id) REFERENCES venues(id) ON DELETE CASCADE;
CREATE INDEX idx_ingredients_venue_id ON ingredients(venue_id);

-- =========================================
-- 7. Add venue_id to products
-- =========================================
ALTER TABLE products ADD COLUMN venue_id INTEGER;
UPDATE products SET venue_id = (SELECT id FROM venues ORDER BY id LIMIT 1);
ALTER TABLE products ALTER COLUMN venue_id SET NOT NULL;
ALTER TABLE products ADD CONSTRAINT fk_products_venue FOREIGN KEY (venue_id) REFERENCES venues(id) ON DELETE CASCADE;
CREATE INDEX idx_products_venue_id ON products(venue_id);

-- =========================================
-- 8. Add venue_id to categories
-- =========================================
ALTER TABLE categories ADD COLUMN venue_id INTEGER;
UPDATE categories SET venue_id = (SELECT id FROM venues ORDER BY id LIMIT 1);
ALTER TABLE categories ALTER COLUMN venue_id SET NOT NULL;
ALTER TABLE categories ADD CONSTRAINT fk_categories_venue FOREIGN KEY (venue_id) REFERENCES venues(id) ON DELETE CASCADE;
CREATE INDEX idx_categories_venue_id ON categories(venue_id);

-- =========================================
-- 9. Add venue_id to tables + fix unique constraint
-- =========================================
ALTER TABLE tables ADD COLUMN venue_id INTEGER;
UPDATE tables SET venue_id = (SELECT id FROM venues ORDER BY id LIMIT 1);
ALTER TABLE tables ALTER COLUMN venue_id SET NOT NULL;
ALTER TABLE tables ADD CONSTRAINT fk_tables_venue FOREIGN KEY (venue_id) REFERENCES venues(id) ON DELETE CASCADE;
CREATE INDEX idx_tables_venue_id ON tables(venue_id);

-- Change table_number from global unique to per-venue unique
ALTER TABLE tables DROP CONSTRAINT IF EXISTS tables_table_number_key;
ALTER TABLE tables ADD CONSTRAINT uq_tables_venue_number UNIQUE (venue_id, table_number);

-- =========================================
-- 10. Update orders: restaurant_id -> venue_id + pos_terminal_id
-- =========================================
ALTER TABLE orders ADD COLUMN venue_id INTEGER;
ALTER TABLE orders ADD COLUMN pos_terminal_id INTEGER;

UPDATE orders SET venue_id = (SELECT id FROM venues ORDER BY id LIMIT 1);
ALTER TABLE orders ALTER COLUMN venue_id SET NOT NULL;

ALTER TABLE orders ADD CONSTRAINT fk_orders_venue FOREIGN KEY (venue_id) REFERENCES venues(id);
ALTER TABLE orders ADD CONSTRAINT fk_orders_pos FOREIGN KEY (pos_terminal_id) REFERENCES pos_terminals(id);
CREATE INDEX idx_orders_venue_id ON orders(venue_id);

ALTER TABLE orders DROP CONSTRAINT IF EXISTS orders_restaurant_id_fkey;
ALTER TABLE orders DROP COLUMN restaurant_id;

-- =========================================
-- 11. Update payments: restaurant_id -> venue_id + pos_terminal_id
-- =========================================
ALTER TABLE payments ADD COLUMN venue_id INTEGER;
ALTER TABLE payments ADD COLUMN pos_terminal_id INTEGER;

UPDATE payments SET venue_id = (SELECT id FROM venues ORDER BY id LIMIT 1);
ALTER TABLE payments ALTER COLUMN venue_id SET NOT NULL;

ALTER TABLE payments ADD CONSTRAINT fk_payments_venue FOREIGN KEY (venue_id) REFERENCES venues(id);
ALTER TABLE payments ADD CONSTRAINT fk_payments_pos FOREIGN KEY (pos_terminal_id) REFERENCES pos_terminals(id);

ALTER TABLE payments DROP CONSTRAINT IF EXISTS payments_restaurant_id_fkey;
DROP INDEX IF EXISTS idx_payments_restaurant_id;
ALTER TABLE payments DROP COLUMN restaurant_id;

CREATE INDEX idx_payments_venue_id ON payments(venue_id);

-- =========================================
-- 12. Drop obsolete tables
-- =========================================
DROP TABLE IF EXISTS table_waitress CASCADE;
DROP TABLE IF EXISTS waitress CASCADE;
DROP TABLE IF EXISTS restaurants CASCADE;

-- =========================================
-- 13. Sync sequences
-- =========================================
SELECT setval('owners_id_seq', COALESCE((SELECT MAX(id) FROM owners), 1));
SELECT setval('venues_id_seq', COALESCE((SELECT MAX(id) FROM venues), 1));
SELECT setval('pos_terminals_id_seq', COALESCE((SELECT MAX(id) FROM pos_terminals), 1));

COMMIT;
