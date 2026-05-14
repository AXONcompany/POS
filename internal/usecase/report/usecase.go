package report

import (
	"context"
	"time"

	"github.com/AXONcompany/POS/internal/infrastructure/persistence/postgres"
)

type Repository interface {
	GetSalesReport(ctx context.Context, venueID int, startDate, endDate time.Time, reportType string) (float64, int, float64, []postgres.SalesReportRow, error)
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

func (uc *Usecase) GetTipsReport(ctx context.Context, venueID int, startDate, endDate time.Time) ([]postgres.TipsReportRow, error) {
	return uc.repo.GetTipsReport(ctx, venueID, startDate, endDate)
}
