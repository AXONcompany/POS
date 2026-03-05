-- ============================================================
-- Migration 000011 DOWN: Rollback restructure
-- ============================================================
-- WARNING: This rollback recreates the old schema but data
-- migration back to restaurant_id is approximate.

BEGIN;

-- 1. Recreate restaurants table
CREATE TABLE IF NOT EXISTS restaurants (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    address TEXT,
    phone VARCHAR(50),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Migrate venues back to restaurants
INSERT INTO restaurants (id, name, address, phone, is_active, created_at, updated_at)
SELECT id, name, COALESCE(address, ''), COALESCE(phone, ''), is_active, created_at, updated_at
FROM venues;

-- 2. Recreate waitress and table_waitress
CREATE TABLE IF NOT EXISTS waitress (
    id_user bigint PRIMARY KEY,
    FOREIGN KEY (id_user) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS table_waitress (
    id bigserial PRIMARY KEY,
    created_at timestamptz NOT NULL DEFAULT NOW(),
    updated_at timestamptz NOT NULL DEFAULT NOW(),
    deleted_at timestamptz NULL,
    table_id bigint NOT NULL,
    waitress_id bigint NOT NULL,
    CONSTRAINT fk_table FOREIGN KEY (table_id) REFERENCES tables(id_table) ON DELETE CASCADE,
    CONSTRAINT fk_waitress FOREIGN KEY (waitress_id) REFERENCES waitress(id_user) ON DELETE CASCADE
);

-- 3. Restore restaurant_id on users
ALTER TABLE users ADD COLUMN restaurant_id INTEGER;
UPDATE users SET restaurant_id = venue_id;
ALTER TABLE users ALTER COLUMN restaurant_id SET NOT NULL;
ALTER TABLE users ADD CONSTRAINT users_restaurant_id_fkey FOREIGN KEY (restaurant_id) REFERENCES restaurants(id) ON DELETE CASCADE;
ALTER TABLE users DROP CONSTRAINT IF EXISTS fk_users_venue;
ALTER TABLE users DROP COLUMN venue_id;
DROP INDEX IF EXISTS idx_users_venue_id;
CREATE INDEX idx_users_restaurant_id ON users(restaurant_id);

-- 4. Restore restaurant_id on orders
ALTER TABLE orders ADD COLUMN restaurant_id INTEGER;
UPDATE orders SET restaurant_id = venue_id;
ALTER TABLE orders ALTER COLUMN restaurant_id SET NOT NULL;
ALTER TABLE orders ADD CONSTRAINT orders_restaurant_id_fkey FOREIGN KEY (restaurant_id) REFERENCES restaurants(id) ON DELETE CASCADE;
ALTER TABLE orders DROP CONSTRAINT IF EXISTS fk_orders_venue;
ALTER TABLE orders DROP CONSTRAINT IF EXISTS fk_orders_pos;
ALTER TABLE orders DROP COLUMN venue_id;
ALTER TABLE orders DROP COLUMN pos_terminal_id;
DROP INDEX IF EXISTS idx_orders_venue_id;

-- 5. Restore restaurant_id on payments
ALTER TABLE payments ADD COLUMN restaurant_id INTEGER;
UPDATE payments SET restaurant_id = venue_id;
ALTER TABLE payments ALTER COLUMN restaurant_id SET NOT NULL;
ALTER TABLE payments ADD CONSTRAINT payments_restaurant_id_fkey FOREIGN KEY (restaurant_id) REFERENCES restaurants(id);
ALTER TABLE payments DROP CONSTRAINT IF EXISTS fk_payments_venue;
ALTER TABLE payments DROP CONSTRAINT IF EXISTS fk_payments_pos;
ALTER TABLE payments DROP COLUMN venue_id;
ALTER TABLE payments DROP COLUMN pos_terminal_id;
DROP INDEX IF EXISTS idx_payments_venue_id;
CREATE INDEX idx_payments_restaurant_id ON payments(restaurant_id);

-- 6. Remove venue_id from data tables
ALTER TABLE ingredients DROP CONSTRAINT IF EXISTS fk_ingredients_venue;
ALTER TABLE ingredients DROP COLUMN IF EXISTS venue_id;
DROP INDEX IF EXISTS idx_ingredients_venue_id;

ALTER TABLE products DROP CONSTRAINT IF EXISTS fk_products_venue;
ALTER TABLE products DROP COLUMN IF EXISTS venue_id;
DROP INDEX IF EXISTS idx_products_venue_id;

ALTER TABLE categories DROP CONSTRAINT IF EXISTS fk_categories_venue;
ALTER TABLE categories DROP COLUMN IF EXISTS venue_id;
DROP INDEX IF EXISTS idx_categories_venue_id;

-- Restore global unique on table_number
ALTER TABLE tables DROP CONSTRAINT IF EXISTS uq_tables_venue_number;
ALTER TABLE tables DROP CONSTRAINT IF EXISTS fk_tables_venue;
ALTER TABLE tables DROP COLUMN IF EXISTS venue_id;
DROP INDEX IF EXISTS idx_tables_venue_id;
ALTER TABLE tables ADD CONSTRAINT tables_table_number_key UNIQUE (table_number);

-- 7. Drop new tables
DROP TABLE IF EXISTS pos_terminals CASCADE;
DROP TABLE IF EXISTS venues CASCADE;
DROP TABLE IF EXISTS owners CASCADE;

COMMIT;
