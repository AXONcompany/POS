INSERT INTO restaurants (id, name) VALUES (1, 'Test Restaurant') ON CONFLICT DO NOTHING;
INSERT INTO roles (id, name) VALUES (1, 'PROPIETARIO'), (2, 'CAJERO'), (3, 'MESERO') ON CONFLICT DO NOTHING;

INSERT INTO users (id, restaurant_id, role_id, name, email, password_hash) VALUES 
(1, 1, 1, 'Admin', 'admin@test.com', '$2a$10$VMOeIt2brgxez8dBUbQVZef7rDYal7Ktw.WexOh6bc270VMWvpfx6'),
(2, 1, 3, 'Mesero Test', 'mesero@test.com', '$2a$10$VMOeIt2brgxez8dBUbQVZef7rDYal7Ktw.WexOh6bc270VMWvpfx6')
ON CONFLICT (email) DO NOTHING;

-- table
INSERT INTO tables (id_table, table_number, capacity, status) VALUES (1, 1, 4, 'LIBRE') ON CONFLICT DO NOTHING;

-- Sincronizar secuencias para no tener errores de primary key
SELECT setval('tables_id_table_seq', (SELECT MAX(id_table) FROM tables));
SELECT setval('users_id_seq', (SELECT MAX(id) FROM users));
SELECT setval('restaurants_id_seq', (SELECT MAX(id) FROM restaurants));
SELECT setval('roles_id_seq', (SELECT MAX(id) FROM roles));
