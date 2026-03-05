CREATE TABLE IF NOT EXISTS payments (
    id SERIAL PRIMARY KEY,
    order_id INTEGER NOT NULL REFERENCES orders(id),
    division_id VARCHAR(50),
    payment_method VARCHAR(20) NOT NULL CHECK (payment_method IN ('efectivo', 'tarjeta', 'multiple')),
    amount NUMERIC(12, 2) NOT NULL,
    tip NUMERIC(12, 2) NOT NULL DEFAULT 0,
    total NUMERIC(12, 2) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pendiente' CHECK (status IN ('pendiente', 'aprobado', 'rechazado')),
    reference VARCHAR(100),
    restaurant_id INTEGER NOT NULL REFERENCES restaurants(id),
    user_id INTEGER NOT NULL REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_payments_order_id ON payments(order_id);
CREATE INDEX idx_payments_restaurant_id ON payments(restaurant_id);
