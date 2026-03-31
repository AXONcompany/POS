package postgres

import (
	"context"
	"fmt"
	"math"

	"github.com/jackc/pgx/v5"

	domainOrder "github.com/AXONcompany/POS/internal/domain/order"
)

type OrderRepository struct {
	db *DB
}

func NewOrderRepository(db *DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) Create(ctx context.Context, o *domainOrder.Order) (*domainOrder.Order, error) {
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	// Insert order
	query := `
		INSERT INTO orders (venue_id, table_id, user_id, pos_terminal_id, status_id, total_amount)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`
	err = tx.QueryRow(ctx, query, o.VenueID, o.TableID, o.UserID, o.POSTerminalID, o.StatusID, o.TotalAmount).Scan(&o.ID, &o.CreatedAt, &o.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("insert order: %w", err)
	}

	// Insert items
	if len(o.Items) > 0 {
		var b pgx.Batch
		itemQuery := `
			INSERT INTO order_items (order_id, product_id, quantity, unit_price, notes)
			VALUES ($1, $2, $3, $4, $5)
		`
		for i := range o.Items {
			o.Items[i].OrderID = o.ID
			b.Queue(itemQuery, o.ID, o.Items[i].ProductID, o.Items[i].Quantity, o.Items[i].UnitPrice, o.Items[i].Notes)
		}

		br := tx.SendBatch(ctx, &b)
		_, err := br.Exec()
		if err != nil {
			br.Close()
			return nil, fmt.Errorf("insert items batch: %w", err)
		}
		if err := br.Close(); err != nil {
			return nil, fmt.Errorf("close batch: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit tx: %w", err)
	}

	return o, nil
}

func (r *OrderRepository) GetByID(ctx context.Context, id int64, venueID int) (*domainOrder.Order, error) {
	query := `
		SELECT o.id, o.venue_id, o.table_id, o.user_id, o.pos_terminal_id, o.status_id, o.total_amount, o.created_at, o.updated_at, os.name as status
		FROM orders o
		JOIN order_statuses os ON o.status_id = os.id
		WHERE o.id = $1 AND o.venue_id = $2 AND o.deleted_at IS NULL
	`
	var o domainOrder.Order
	err := r.db.Pool.QueryRow(ctx, query, id, venueID).Scan(
		&o.ID, &o.VenueID, &o.TableID, &o.UserID, &o.POSTerminalID, &o.StatusID, &o.TotalAmount, &o.CreatedAt, &o.UpdatedAt, &o.Status,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("order not found")
		}
		return nil, fmt.Errorf("get order by id: %w", err)
	}

	// Get items
	itemsQuery := `
		SELECT id, order_id, product_id, quantity, unit_price, notes, created_at
		FROM order_items
		WHERE order_id = $1
	`
	rows, err := r.db.Pool.Query(ctx, itemsQuery, id)
	if err != nil {
		return nil, fmt.Errorf("query order items: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var item domainOrder.OrderItem
		if err := rows.Scan(&item.ID, &item.OrderID, &item.ProductID, &item.Quantity, &item.UnitPrice, &item.Notes, &item.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan order item: %w", err)
		}
		o.Items = append(o.Items, item)
	}

	return &o, nil
}

func (r *OrderRepository) GetStatusByID(ctx context.Context, id int64, venueID int) (int, error) {
	var statusID int
	err := r.db.Pool.QueryRow(ctx,
		`SELECT status_id FROM orders WHERE id = $1 AND venue_id = $2 AND deleted_at IS NULL`,
		id, venueID,
	).Scan(&statusID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return 0, fmt.Errorf("order not found")
		}
		return 0, fmt.Errorf("get order status: %w", err)
	}
	return statusID, nil
}

func (r *OrderRepository) UpdateStatus(ctx context.Context, id int64, venueID int, statusID int) error {
	query := `
		UPDATE orders 
		SET status_id = $1, updated_at = NOW()
		WHERE id = $2 AND venue_id = $3 AND deleted_at IS NULL
	`
	cmd, err := r.db.Pool.Exec(ctx, query, statusID, id, venueID)
	if err != nil {
		return fmt.Errorf("update status: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("order not found or not updated")
	}

	return nil
}

func (r *OrderRepository) GetOrderItem(ctx context.Context, itemID, orderID int64) (*domainOrder.OrderItem, error) {
	var item domainOrder.OrderItem
	err := r.db.Pool.QueryRow(ctx,
		`SELECT id, order_id, product_id, quantity, unit_price, notes, created_at, cancelled_at
		 FROM order_items WHERE id = $1 AND order_id = $2`,
		itemID, orderID,
	).Scan(&item.ID, &item.OrderID, &item.ProductID, &item.Quantity, &item.UnitPrice, &item.Notes, &item.CreatedAt, &item.CancelledAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("order item not found")
		}
		return nil, fmt.Errorf("get order item: %w", err)
	}
	return &item, nil
}

func (r *OrderRepository) CancelItemWithInventoryRestore(ctx context.Context, itemID, orderID int64, venueID int, restorations []domainOrder.StockDeduction) error {
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	// Soft delete del item — si ya estaba cancelado, 0 rows afectadas
	cmd, err := tx.Exec(ctx,
		`UPDATE order_items SET cancelled_at = NOW() WHERE id = $1 AND cancelled_at IS NULL`,
		itemID,
	)
	if err != nil {
		return fmt.Errorf("cancel order item: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return domainOrder.ErrItemAlreadyCancelled
	}

	// Restaurar stock por ingrediente
	for _, r := range restorations {
		qty := int64(math.Round(r.Quantity))
		_, err := tx.Exec(ctx,
			`UPDATE ingredients SET stock = stock + $1 WHERE id = $2 AND deleted_at IS NULL`,
			qty, r.IngredientID,
		)
		if err != nil {
			return fmt.Errorf("restore stock ingredient %d: %w", r.IngredientID, err)
		}
	}

	// Restar del total de la orden (obtenemos unit_price * quantity del item cancelado)
	_, err = tx.Exec(ctx,
		`UPDATE orders
		 SET total_amount = total_amount - (
		     SELECT unit_price * quantity FROM order_items WHERE id = $1
		 ), updated_at = NOW()
		 WHERE id = $2 AND venue_id = $3 AND deleted_at IS NULL`,
		itemID, orderID, venueID,
	)
	if err != nil {
		return fmt.Errorf("update order total: %w", err)
	}

	return tx.Commit(ctx)
}

func (r *OrderRepository) AddItemsWithInventory(ctx context.Context, orderID int64, venueID int, items []domainOrder.OrderItem, deductions []domainOrder.StockDeduction) error {
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	// INSERT batch en order_items
	if len(items) > 0 {
		var b pgx.Batch
		itemQuery := `INSERT INTO order_items (order_id, product_id, quantity, unit_price, notes) VALUES ($1, $2, $3, $4, $5)`
		for _, item := range items {
			b.Queue(itemQuery, orderID, item.ProductID, item.Quantity, item.UnitPrice, item.Notes)
		}
		br := tx.SendBatch(ctx, &b)
		_, err := br.Exec()
		if err != nil {
			br.Close()
			return fmt.Errorf("insert items batch: %w", err)
		}
		if err := br.Close(); err != nil {
			return fmt.Errorf("close batch: %w", err)
		}
	}

	// UPDATE stock por ingrediente con guard atomico
	stockQuery := `UPDATE ingredients SET stock = stock - $1 WHERE id = $2 AND deleted_at IS NULL AND stock >= $1`
	for _, d := range deductions {
		qty := int64(math.Round(d.Quantity))
		cmd, err := tx.Exec(ctx, stockQuery, qty, d.IngredientID)
		if err != nil {
			return fmt.Errorf("update stock ingredient %d: %w", d.IngredientID, err)
		}
		if cmd.RowsAffected() == 0 {
			return domainOrder.ErrInsufficientStock
		}
	}

	// UPDATE total_amount de la orden
	var addedTotal float64
	for _, item := range items {
		addedTotal += item.UnitPrice * float64(item.Quantity)
	}
	_, err = tx.Exec(ctx,
		`UPDATE orders SET total_amount = total_amount + $1, updated_at = NOW() WHERE id = $2 AND venue_id = $3 AND deleted_at IS NULL`,
		addedTotal, orderID, venueID,
	)
	if err != nil {
		return fmt.Errorf("update order total: %w", err)
	}

	return tx.Commit(ctx)
}

func (r *OrderRepository) ListByTable(ctx context.Context, tableID int64, venueID int) ([]domainOrder.Order, error) {
	query := `
		SELECT o.id, o.venue_id, o.table_id, o.user_id, o.pos_terminal_id, o.status_id, o.total_amount, o.created_at, o.updated_at, os.name as status
		FROM orders o
		JOIN order_statuses os ON o.status_id = os.id
		WHERE o.table_id = $1 AND o.venue_id = $2 AND o.deleted_at IS NULL
		ORDER BY o.created_at DESC
	`
	rows, err := r.db.Pool.Query(ctx, query, tableID, venueID)
	if err != nil {
		return nil, fmt.Errorf("list orders by table query: %w", err)
	}
	defer rows.Close()

	var orders []domainOrder.Order
	for rows.Next() {
		var o domainOrder.Order
		if err := rows.Scan(&o.ID, &o.VenueID, &o.TableID, &o.UserID, &o.POSTerminalID, &o.StatusID, &o.TotalAmount, &o.CreatedAt, &o.UpdatedAt, &o.Status); err != nil {
			return nil, fmt.Errorf("scan order: %w", err)
		}
		orders = append(orders, o)
	}

	return orders, nil
}
