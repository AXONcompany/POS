package postgres

import (
	"context"
	"encoding/json"
	"fmt"

	domainAudit "github.com/AXONcompany/POS/internal/domain/audit"
)

type AuditRepository struct {
	db *DB
}

func NewAuditRepository(db *DB) *AuditRepository {
	return &AuditRepository{db: db}
}

func (r *AuditRepository) SaveAudit(ctx context.Context, entry *domainAudit.AuditEntry) error {
	oldVal, err := json.Marshal(entry.OldValue)
	if err != nil {
		return fmt.Errorf("marshal old_value: %w", err)
	}
	newVal, err := json.Marshal(entry.NewValue)
	if err != nil {
		return fmt.Errorf("marshal new_value: %w", err)
	}

	_, err = r.db.Pool.Exec(ctx,
		`INSERT INTO audit_log (entity_type, entity_id, action, old_value, new_value, user_id, venue_id)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		entry.EntityType, entry.EntityID, entry.Action, oldVal, newVal, entry.UserID, entry.VenueID,
	)
	if err != nil {
		return fmt.Errorf("save audit: %w", err)
	}
	return nil
}
