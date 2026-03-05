INSERT INTO restaurants (id, name) VALUES (1, 'Test Restaurant') ON CONFLICT DO NOTHING;
INSERT INTO roles (id, name) VALUES (1, 'PROPIETARIO'), (2, 'CAJERO'), (3, 'MESERO') ON CONFLICT DO NOTHING;

INSERT INTO users (id, restaurant_id, role_id, name, email, password_hash) VALUES 
(1, 1, 1, 'Admin', 'admin@test.com', '$2a$10$7/Oa2jU.eT/09QOFsOhC4u8WvR.9L0oF1k1xJ68jEIf4iP3r0gD5K')
ON CONFLICT (email) DO NOTHING;

INSERT INTO users (id, restaurant_id, role_id, name, email, password_hash) VALUES 
(2, 1, 3, 'Mesero Test', 'mesero@test.com', '$2a$10$7/Oa2jU.eT/09QOFsOhC4u8WvR.9L0oF1k1xJ68jEIf4iP3r0gD5K')
ON CONFLICT (email) DO NOTHING;

-- table
INSERT INTO tables (id, restaurant_id, table_number, capacity, status) VALUES (1, 1, 1, 4, 'free') ON CONFLICT DO NOTHING;
