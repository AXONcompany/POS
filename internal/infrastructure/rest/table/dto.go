package table

// CreateRequest maps to the frontend TableData creation payload.
type CreateRequest struct {
	Name       string  `json:"name"`
	Seats      int     `json:"seats" binding:"required,gt=0"`
	X          int     `json:"x"`
	Y          int     `json:"y"`
	Width      int     `json:"width"`
	Height     int     `json:"height"`
	Shape      string  `json:"shape"`
	Rotation   int     `json:"rotation"`
	Color      *string `json:"color"`
	Floor      int     `json:"floor"`
	IsMerged   bool    `json:"isMerged"`
	MergedFrom []int64 `json:"mergedFrom"`
}

// FullUpdateRequest replaces all mutable fields of a table (PUT /mesas/:id).
type FullUpdateRequest struct {
	Name             string  `json:"name"`
	Seats            int     `json:"seats"`
	Status           string  `json:"status"`
	Guests           int     `json:"guests"`
	AssignedWaiterId *int    `json:"assignedWaiterId"`
	CheckInTime      *string `json:"checkInTime"`
	X                int     `json:"x"`
	Y                int     `json:"y"`
	Width            int     `json:"width"`
	Height           int     `json:"height"`
	Shape            string  `json:"shape"`
	Rotation         int     `json:"rotation"`
	Color            *string `json:"color"`
	Floor            int     `json:"floor"`
	IsMerged         bool    `json:"isMerged"`
	MergedFrom       []int64 `json:"mergedFrom"`
}

// UpdateEstadoRequest is used by PATCH /mesas/:id/estado (status-only update).
type UpdateEstadoRequest struct {
	Estado string `json:"estado" binding:"required"`
}

// Response matches the frontend TableData interface (camelCase JSON keys).
type Response struct {
	ID               int64   `json:"id"`
	Name             string  `json:"name"`
	Seats            int     `json:"seats"`
	Status           string  `json:"status"`
	Guests           int     `json:"guests"`
	AssignedWaiterId *int    `json:"assignedWaiterId,omitempty"`
	CheckInTime      *string `json:"checkInTime,omitempty"`
	X                int     `json:"x"`
	Y                int     `json:"y"`
	Width            int     `json:"width"`
	Height           int     `json:"height"`
	Shape            string  `json:"shape"`
	Rotation         int     `json:"rotation"`
	Color            *string `json:"color,omitempty"`
	Floor            int     `json:"floor"`
	IsMerged         bool    `json:"isMerged"`
	MergedFrom       []int64 `json:"mergedFrom,omitempty"`
}
