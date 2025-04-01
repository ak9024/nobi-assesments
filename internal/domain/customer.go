package domain

type Customer struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Balance  float64 `json:"balance"`
	Units    float64 `json:"units"`
	IsActive bool    `json:"is_active"`
}
