CREATE TABLE IF NOT EXISTS table_assignments (
    id BIGSERIAL PRIMARY KEY,
    table_id BIGINT NOT NULL REFERENCES tables(id_table) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    venue_id INTEGER NOT NULL REFERENCES venues(id) ON DELETE CASCADE,
    assigned_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    unassigned_at TIMESTAMPTZ NULL
);

CREATE INDEX idx_table_assignments_table ON table_assignments(table_id);
CREATE INDEX idx_table_assignments_user ON table_assignments(user_id);
CREATE INDEX idx_table_assignments_venue ON table_assignments(venue_id);
