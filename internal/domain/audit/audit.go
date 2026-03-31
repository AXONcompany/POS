package audit

import "time"

type AuditEntry struct {
	ID         int64
	EntityType string
	EntityID   int64
	Action     string
	OldValue   interface{}
	NewValue   interface{}
	UserID     int
	VenueID    int
	CreatedAt  time.Time
}
