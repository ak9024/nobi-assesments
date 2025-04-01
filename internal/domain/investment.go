package domain

type Investment struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	TotalUnits   float64 `json:"total_units"`
	TotalBalance float64 `json:"total_balance"`
	NAB          float64 `json:"nab"`
}

type CustomerInvestment struct {
	ID           string  `json:"id"`
	CustomerID   string  `json:"customer_id"`
	InvestmentID string  `json:"investment_id"`
	Units        float64 `json:"units"`
}

type CustomerPortfolio struct {
	Customer   string     `json:"customer_id"`
	Investment Investment `json:"investment"`
	Portfolio  struct {
		ID      string  `json:"id"`
		Units   float64 `json:"units"`
		Balance float64 `json:"balance"`
	} `json:"portfolio"`
}
