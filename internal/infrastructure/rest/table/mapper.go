package table

import (
	"time"

	domainTable "github.com/AXONcompany/POS/internal/domain/table"
)

// ToDomain converts a CreateRequest to a domain Table.
func ToDomain(req CreateRequest) *domainTable.Table {
	width := req.Width
	if width == 0 {
		width = 110
	}
	height := req.Height
	if height == 0 {
		height = 110
	}
	shape := req.Shape
	if shape == "" {
		shape = "square"
	}
	floor := req.Floor
	if floor == 0 {
		floor = 1
	}
	name := req.Name
	if name == "" {
		name = "Mesa"
	}

	return &domainTable.Table{
		Name:       name,
		Capacity:   req.Seats,
		Status:     "free",
		X:          req.X,
		Y:          req.Y,
		Width:      width,
		Height:     height,
		Shape:      shape,
		Rotation:   req.Rotation,
		Color:      req.Color,
		Floor:      floor,
		IsMerged:   req.IsMerged,
		MergedFrom: req.MergedFrom,
	}
}

// ToFullUpdateDomain converts a FullUpdateRequest to a domain Table for full replacement.
func ToFullUpdateDomain(req FullUpdateRequest) domainTable.Table {
	t := domainTable.Table{
		Name:             req.Name,
		Capacity:         req.Seats,
		Status:           req.Status,
		Guests:           req.Guests,
		X:                req.X,
		Y:                req.Y,
		Width:            req.Width,
		Height:           req.Height,
		Shape:            req.Shape,
		Rotation:         req.Rotation,
		Color:            req.Color,
		Floor:            req.Floor,
		IsMerged:         req.IsMerged,
		MergedFrom:       req.MergedFrom,
		AssignedWaiterID: req.AssignedWaiterId,
	}
	if req.CheckInTime != nil {
		if parsed, err := time.Parse(time.RFC3339, *req.CheckInTime); err == nil {
			t.ArrivalTime = &parsed
		}
	}
	return t
}

// ToResponse converts a domain Table to the REST response (camelCase keys match frontend).
func ToResponse(t *domainTable.Table) Response {
	r := Response{
		ID:         t.ID,
		Name:       t.Name,
		Seats:      t.Capacity,
		Status:     t.Status,
		Guests:     t.Guests,
		X:          t.X,
		Y:          t.Y,
		Width:      t.Width,
		Height:     t.Height,
		Shape:      t.Shape,
		Rotation:   t.Rotation,
		Color:      t.Color,
		Floor:      t.Floor,
		IsMerged:   t.IsMerged,
		MergedFrom: t.MergedFrom,
	}
	if t.AssignedWaiterID != nil {
		r.AssignedWaiterId = t.AssignedWaiterID
	}
	if t.ArrivalTime != nil {
		s := t.ArrivalTime.Format(time.RFC3339)
		r.CheckInTime = &s
	}
	return r
}

// ToResponseList converts a slice of domain Tables to REST responses.
func ToResponseList(tables []domainTable.Table) []Response {
	responses := make([]Response, len(tables))
	for i, t := range tables {
		responses[i] = ToResponse(&t)
	}
	return responses
}
