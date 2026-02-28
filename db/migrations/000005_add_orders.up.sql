CREATE TABLE IF NOT EXISTS order_statuses (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL,
    description TEXT
);

CREATE TABLE IF NOT EXISTS orders (
    id BIGSERIAL PRIMARY KEY,
    restaurant_id INTEGER NOT NULL REFERENCES restaurants(id) ON DELETE CASCADE,
    table_id BIGINT REFERENCES tables(id_table) ON DELETE SET NULL,
    user_id INTEGER NOT NULL REFERENCES users(id), -- mesero
    status_id INTEGER NOT NULL REFERENCES order_statuses(id),
    total_amount DECIMAL(10,2) NOT NULL DEFAULT 0.00,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ NULL
);

CREATE TABLE IF NOT EXISTS order_items (
    id BIGSERIAL PRIMARY KEY,
    order_id BIGINT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    product_id BIGINT NOT NULL REFERENCES products(id),
    quantity INTEGER NOT NULL,
    unit_price DECIMAL(10,2) NOT NULL,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Insert default statuses
INSERT INTO order_statuses (name, description) VALUES 
    ('PENDING', 'Pedido recién creado, pendiente de preparación'),
    ('PREPARING', 'Pedido en preparación en cocina/barra'),
    ('READY', 'Pedido listo para ser entregado a la mesa'),
    ('DELIVERED', 'Pedido entregado en la mesa'),
    ('PAID', 'Pedido pagado finalizado'),
    ('CANCELLED', 'Pedido cancelado');
