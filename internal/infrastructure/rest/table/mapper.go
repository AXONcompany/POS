package table

import (
	"github.com/AXONcompany/POS/internal/domain/table"
)

// ToDomain convierte el JSON de entrada a la Entidad de Dominio
func ToDomain(req CreateRequest) *table.Table {
	return &table.Table{
		Number:   req.Number,
		Capacity: req.Capacity,
		Status:   req.Status,
	}
}

// ToUpdateDomain convierte el JSON de actualización a la estructura de updates del dominio
func ToUpdateDomain(req UpdateRequest) *table.TableUpdates {
	return &table.TableUpdates{
		Number:      req.Number,
		Capacity:    req.Capacity,
		Status:      req.Status,
		ArrivalTime: req.ArrivalTime,
	}
}

// ToResponse convierte la Entidad de Dominio a JSON de respuesta
func ToResponse(t *table.Table) Response {
	return Response{
		ID:          t.ID,
		Number:      t.Number,
		Capacity:    t.Capacity,
		Status:      t.Status,
		ArrivalTime: t.ArrivalTime,
		CreatedAt:   t.CreatedAt,
	}
}

// ToResponseList convierte una lista completa
func ToResponseList(tables []table.Table) []Response {
	responses := make([]Response, len(tables))
	for i, t := range tables {
		responses[i] = ToResponse(&t)
	}
	return responses
}
