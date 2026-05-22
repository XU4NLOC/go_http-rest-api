package models

type Transaction struct {
	ID          int     `json:"id"`
	AccountID   int     `json:"account_id"`
	Amount      float64 `json:"amount"`
	Type        string  `json:"type"` // income | expense
	Description string  `json:"description"`
	CreatedAt   string  `json:"created_at"`
}

type Summary struct {
	AccountID     int     `json:"account_id"`
	TotalIncome   float64 `json:"total_income"`
	TotalExpenses float64 `json:"total_expenses"`
	Balance       float64 `json:"balance"`
}
