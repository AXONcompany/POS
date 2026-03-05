package report

import (
	"context"
	"time"

	"github.com/AXONcompany/POS/internal/domain/ingredient"
	"github.com/AXONcompany/POS/internal/infrastructure/persistence/postgres"
)

type Repository interface {
	GetSalesReport(ctx context.Context, venueID int, startDate, endDate time.Time, reportType string) (float64, int, float64, []postgres.SalesReportRow, error)
	GetInventoryReport(ctx context.Context, venueID int) ([]*ingredient.Ingredient, float64, error)
	GetTipsReport(ctx context.Context, venueID int, startDate, endDate time.Time) ([]postgres.TipsReportRow, error)
}

type Usecase struct {
	repo Repository
}

func NewUsecase(repo Repository) *Usecase {
	return &Usecase{repo: repo}
}

func (uc *Usecase) GetSalesReport(ctx context.Context, venueID int, startDate, endDate time.Time, reportType string) (map[string]interface{}, error) {
	totalSales, totalOrders, avgTicket, details, err := uc.repo.GetSalesReport(ctx, venueID, startDate, endDate, reportType)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"periodo": map[string]string{
			"inicio": startDate.Format("2006-01-02"),
			"fin":    endDate.Format("2006-01-02"),
		},
		"total_ventas":    totalSales,
		"total_ordenes":   totalOrders,
		"ticket_promedio": avgTicket,
		"detalle":         details,
	}, nil
}

func (uc *Usecase) GetInventoryReport(ctx context.Context, venueID int) (map[string]interface{}, error) {
	items, totalValue, err := uc.repo.GetInventoryReport(ctx, venueID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"ingredientes_bajo_stock": items,
		"valor_total_inventario":  totalValue,
	}, nil
}

func (uc *Usecase) GetTipsReport(ctx context.Context, venueID int, startDate, endDate time.Time) ([]postgres.TipsReportRow, error) {
	return uc.repo.GetTipsReport(ctx, venueID, startDate, endDate)
}
