package db

type Document struct {
	ID       int       `json:"id"`
	Document string    `json:"document"`
	Vector   []float64 `json:"vector"`
}
