ALTER TABLE order_items ADD COLUMN IF NOT EXISTS cancelled_at timestamptz NULL;

CREATE INDEX IF NOT EXISTS idx_order_items_cancelled_at ON order_items(cancelled_at);
