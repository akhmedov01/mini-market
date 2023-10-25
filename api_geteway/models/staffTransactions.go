package models

type StaffTransaction struct {
	ID              string  `json:"id"`
	SaleID          string  `json:"sale_id"`
	StaffID         string  `json:"staff_id"`
	TransactionType string  `json:"transaction_type"`
	SourceType      string  `json:"source_type"`
	Amount          float64 `json:"amount"`
	Text            string  `json:"information_about"`
	Created_at      string  `json:"created_at"`
	Updated_at      string  `json:"updated_at"`
}

type CreateStaffTransaction struct {
	SaleID          string  `json:"sale_id"`
	StaffID         string  `json:"staff_id"`
	TransactionType string  `json:"transaction_type"`
	SourceType      string  `json:"source_type"`
	Amount          float64 `json:"amount"`
	Text            string  `json:"information_about"`
}

type GetAllStaffTransaction struct {
	StaffTransactions []StaffTransaction
	Count             int
}

type GetAllStaffTransactionRequest struct {
	Page            int
	Limit           int
	TransactionType string
}
