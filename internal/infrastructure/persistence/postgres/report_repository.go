package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/AXONcompany/POS/internal/domain/ingredient"
)

// ReportRepository ejecuta queries de reportes directamente.
type ReportRepository struct {
	db *DB
}

func NewReportRepository(db *DB) *ReportRepository {
	return &ReportRepository{db: db}
}

// SalesReportRow contiene una fila del reporte de ventas.
type SalesReportRow struct {
	Period      string  `json:"periodo"`
	TotalSales  float64 `json:"total_ventas"`
	TotalOrders int     `json:"total_ordenes"`
}

func (r *ReportRepository) GetSalesReport(ctx context.Context, venueID int, startDate, endDate time.Time, reportType string) (float64, int, float64, []SalesReportRow, error) {
	summaryQuery := `
		SELECT COALESCE(SUM(p.total), 0), COUNT(DISTINCT p.order_id), COALESCE(AVG(p.total), 0)
		FROM payments p
		WHERE p.venue_id = $1
		  AND p.status = 'aprobado'
		  AND p.created_at >= $2
		  AND p.created_at < $3`

	var totalSales float64
	var totalOrders int
	var avgTicket float64
	err := r.db.Pool.QueryRow(ctx, summaryQuery, venueID, startDate, endDate.Add(24*time.Hour)).Scan(&totalSales, &totalOrders, &avgTicket)
	if err != nil {
		return 0, 0, 0, nil, fmt.Errorf("sales report summary: %w", err)
	}

	detailQuery := `
		SELECT DATE(p.created_at) as period, SUM(p.total), COUNT(*)
		FROM payments p
		WHERE p.venue_id = $1
		  AND p.status = 'aprobado'
		  AND p.created_at >= $2
		  AND p.created_at < $3
		GROUP BY DATE(p.created_at)
		ORDER BY period`

	rows, err := r.db.Pool.Query(ctx, detailQuery, venueID, startDate, endDate.Add(24*time.Hour))
	if err != nil {
		return 0, 0, 0, nil, fmt.Errorf("sales report detail: %w", err)
	}
	defer rows.Close()

	details := make([]SalesReportRow, 0)
	for rows.Next() {
		var row SalesReportRow
		var date time.Time
		if err := rows.Scan(&date, &row.TotalSales, &row.TotalOrders); err != nil {
			return 0, 0, 0, nil, fmt.Errorf("scan row: %w", err)
		}
		row.Period = date.Format("2006-01-02")
		details = append(details, row)
	}

	return totalSales, totalOrders, avgTicket, details, nil
}

func (r *ReportRepository) GetInventoryReport(ctx context.Context, venueID int) ([]*ingredient.Ingredient, float64, error) {
	query := `
		SELECT id, venue_id, ingredient_name, unit_of_measure, stock, ingredient_type, created_at
		FROM ingredients
		WHERE venue_id = $1 AND stock < 10 AND deleted_at IS NULL
		ORDER BY stock ASC`

	rows, err := r.db.Pool.Query(ctx, query, venueID)
	if err != nil {
		return nil, 0, fmt.Errorf("inventory report: %w", err)
	}
	defer rows.Close()

	items := make([]*ingredient.Ingredient, 0)
	for rows.Next() {
		i := &ingredient.Ingredient{}
		if err := rows.Scan(&i.ID, &i.VenueID, &i.Name, &i.UnitOfMeasure, &i.Stock, &i.IngredientType, &i.CreatedAt); err != nil {
			return nil, 0, fmt.Errorf("scan ingredient: %w", err)
		}
		items = append(items, i)
	}

	var totalValue float64
	err = r.db.Pool.QueryRow(ctx, `SELECT COALESCE(SUM(stock * 1000), 0) FROM ingredients WHERE venue_id = $1 AND deleted_at IS NULL`, venueID).Scan(&totalValue)
	if err != nil {
		return nil, 0, fmt.Errorf("total value: %w", err)
	}

	return items, totalValue, nil
}

// TipsReportRow contiene el resumen de propinas por mesero.
type TipsReportRow struct {
	UserID      int     `json:"mesero_id"`
	UserName    string  `json:"mesero_nombre"`
	TotalTips   float64 `json:"total_propinas"`
	OrdersCount int     `json:"numero_ordenes"`
}

func (r *ReportRepository) GetTipsReport(ctx context.Context, venueID int, startDate, endDate time.Time) ([]TipsReportRow, error) {
	query := `
		SELECT p.user_id, u.name, COALESCE(SUM(p.tip), 0), COUNT(*)
		FROM payments p
		JOIN users u ON u.id = p.user_id
		WHERE p.venue_id = $1
		  AND p.status = 'aprobado'
		  AND p.tip > 0
		  AND p.created_at >= $2
		  AND p.created_at < $3
		GROUP BY p.user_id, u.name
		ORDER BY SUM(p.tip) DESC`

	rows, err := r.db.Pool.Query(ctx, query, venueID, startDate, endDate.Add(24*time.Hour))
	if err != nil {
		return nil, fmt.Errorf("tips report: %w", err)
	}
	defer rows.Close()

	result := make([]TipsReportRow, 0)
	for rows.Next() {
		var row TipsReportRow
		if err := rows.Scan(&row.UserID, &row.UserName, &row.TotalTips, &row.OrdersCount); err != nil {
			return nil, fmt.Errorf("scan tip row: %w", err)
		}
		result = append(result, row)
	}

	return result, nil
}
