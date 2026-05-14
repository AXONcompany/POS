CREATE TABLE IF NOT EXISTS order_divisions (
    id VARCHAR(50) PRIMARY KEY,
    order_id INTEGER NOT NULL REFERENCES orders(id),
    venue_id INTEGER NOT NULL REFERENCES venues(id),
    division_type VARCHAR(20) NOT NULL,
    amount NUMERIC(12, 2) NOT NULL,
    tax NUMERIC(12, 2) NOT NULL,
    total NUMERIC(12, 2) NOT NULL,
    is_paid BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

ALTER TABLE payments ADD CONSTRAINT fk_payments_division FOREIGN KEY (division_id) REFERENCES order_divisions(id);
