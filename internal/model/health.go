package model

// Health is the answer to GET /api/health. It is a struct rather than a map so
// the field order in the JSON is fixed.
type Health struct {
	Status string `json:"status"`
}
