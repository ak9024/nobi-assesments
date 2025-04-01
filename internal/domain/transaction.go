package domain

import "time"

type Transaction struct {
	ID              string    `json:"id"`
	CustomerID      string    `json:"customer_id"`
	InvestmentID    string    `json:"investment_id"`
	Type            string    `json:"type"` // DEPOSIT or WITHDRAW
	Amount          float64   `json:"amount"`
	Units           float64   `json:"units"`
	NAB             float64   `json:"nab"`
	TransactionDate time.Time `json:"transaction_date"`
}

type DepositRequest struct {
	CustomerID   string  `json:"customer_id"`
	InvestmentID string  `json:"investment_id"`
	Amount       float64 `json:"amount"`
}

type WithdrawRequest struct {
	CustomerID   string  `json:"customer_id"`
	InvestmentID string  `json:"investment_id"`
	Amount       float64 `json:"amount"`
}

type TransactionResponse struct {
	TransactionID  string  `json:"transaction_id"`
	Message        string  `json:"message"`
	Amount         float64 `json:"amount"`
	Units          float64 `json:"units_added,omitempty"`
	UnitsReduced   float64 `json:"units_reduced,omitempty"`
	NAB            float64 `json:"nab"`
	TotalUnits     float64 `json:"total_units,omitempty"`
	RemainingUnits float64 `json:"remaining_units,omitempty"`
	CurrentBalance float64 `json:"current_balance"`
}
