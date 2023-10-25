package models

type CreateStaffTarif struct {
	Name          string  `json:"name"`
	Type          string  `json:"type"`
	AmountForCash float64 `json:"amount_for_cash"`
	AmountForCard float64 `json:"amount_for_card"`
	FoundedAt     string  `json:"founded_at"`
}

type StaffTarif struct {
	Id            string  `json:"id"`
	Name          string  `json:"name"`
	Type          string  `json:"type"`
	AmountForCash float64 `json:"amount_for_cash"`
	AmountForCard float64 `json:"amount_for_card"`
	FoundedAt     string  `json:"founded_at"`
	CreatedAt     string  `json:"created_at"`
	Updated_at    string  `json:"updated_at"`
}

type GetAllStaffTarifRequest struct {
	Page  int
	Limit int
	Name  string
}

type GetAllStaffTarif struct {
	StaffTarifs []StaffTarif
	Count       int
}
