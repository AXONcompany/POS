-- Seed data para la estructura Owner > Venue > POS Terminal
-- Password hash corresponde a: "password123"

-- Roles
INSERT INTO roles (id, name) VALUES (1, 'PROPIETARIO'), (2, 'CAJERO'), (3, 'MESERO') ON CONFLICT DO NOTHING;

-- Owner (propietario)
INSERT INTO owners (id, name, email, password_hash) VALUES
(1, 'Propietario Demo', 'owner@test.com', '$2a$10$wN2G4860D4/E4o.VWeHqAeVhI.PttO.I1vA3u9z/3Kk.m7P.Z8T42')
ON CONFLICT DO NOTHING;

-- Venue (sede)
INSERT INTO venues (id, owner_id, name, address, phone) VALUES
(1, 1, 'Sede Principal', 'Calle 1 #10-20, Centro', '3001234567')
ON CONFLICT DO NOTHING;

-- POS Terminal
INSERT INTO pos_terminals (id, venue_id, terminal_name) VALUES
(1, 1, 'Caja 1')
ON CONFLICT DO NOTHING;

-- Usuarios de la sede
INSERT INTO users (id, venue_id, role_id, name, email, password_hash) VALUES
(1, 1, 1, 'Admin', 'admin@test.com', '$2a$10$wN2G4860D4/E4o.VWeHqAeVhI.PttO.I1vA3u9z/3Kk.m7P.Z8T42'),
(2, 1, 2, 'Cajero Test', 'cajero@test.com', '$2a$10$wN2G4860D4/E4o.VWeHqAeVhI.PttO.I1vA3u9z/3Kk.m7P.Z8T42'),
(3, 1, 3, 'Mesero Test', 'mesero@test.com', '$2a$10$wN2G4860D4/E4o.VWeHqAeVhI.PttO.I1vA3u9z/3Kk.m7P.Z8T42')
ON CONFLICT (email) DO NOTHING;

-- Mesa de ejemplo
INSERT INTO tables (id, venue_id, table_number, capacity, status) VALUES
(1, 1, 1, 4, 'free')
ON CONFLICT DO NOTHING;
