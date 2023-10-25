package models

type Sale struct {
	ID              string  `json:"id"`
	BranchId        string  `json:"branch_id"`
	ShopAssistentId string  `json:"shop_assistent_id"`
	CashierId       string  `json:"cashier_id"`
	PaymentType     string  `json:"payment_type"`
	Status          string  `json:"status"`
	ClientName      string  `json:"client_name"`
	Price           float64 `json:"price"`
	Created_at      string  `json:"created_at"`
	Updated_at      string  `json:"updated_at"`
}

type CreateSale struct {
	BranchId        string  `json:"branch_id"`
	ShopAssistentId string  `json:"shop_assistent_id"`
	CashierId       string  `json:"cashier_id"`
	PaymentType     string  `json:"payment_type"`
	Status          string  `json:"status"`
	ClientName      string  `json:"client_name"`
	Price           float64 `json:"price"`
}

type GetAllSale struct {
	Sales []Sale
	Count int
}

type GetAllSaleRequest struct {
	Page       int
	Limit      int
	ClientName string
}
