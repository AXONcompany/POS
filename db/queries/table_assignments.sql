-- name: AssignWaiterToTable :one
INSERT INTO table_assignments (table_id, user_id, venue_id)
VALUES ($1, $2, $3)
RETURNING *;

-- name: UnassignWaiterFromTable :exec
UPDATE table_assignments
SET unassigned_at = NOW()
WHERE table_id = $1 AND venue_id = $2 AND unassigned_at IS NULL;

-- name: GetActiveAssignment :one
SELECT * FROM table_assignments
WHERE table_id = $1 AND venue_id = $2 AND unassigned_at IS NULL
LIMIT 1;

-- name: ListAssignmentsByTable :many
SELECT ta.*, u.name AS waiter_name
FROM table_assignments ta
JOIN users u ON u.id = ta.user_id
WHERE ta.table_id = $1 AND ta.venue_id = $2
ORDER BY ta.assigned_at DESC;

-- name: ListActiveAssignmentsByVenue :many
SELECT ta.*, u.name AS waiter_name, t.table_number
FROM table_assignments ta
JOIN users u ON u.id = ta.user_id
JOIN tables t ON t.id_table = ta.table_id
WHERE ta.venue_id = $1 AND ta.unassigned_at IS NULL
ORDER BY t.table_number;
