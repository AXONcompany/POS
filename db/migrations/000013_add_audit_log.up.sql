CREATE TABLE IF NOT EXISTS audit_log (
    id         bigserial    PRIMARY KEY,
    created_at timestamptz  NOT NULL DEFAULT now(),
    entity_type varchar(64) NOT NULL,
    entity_id  bigint       NOT NULL,
    action     varchar(32)  NOT NULL,
    old_value  jsonb,
    new_value  jsonb,
    user_id    int          NOT NULL,
    venue_id   int          NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_audit_log_entity ON audit_log(entity_type, entity_id);
CREATE INDEX IF NOT EXISTS idx_audit_log_user   ON audit_log(user_id);
