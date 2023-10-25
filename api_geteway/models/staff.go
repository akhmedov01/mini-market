package models

type Staff struct {
	ID            string  `json:"id"`
	BranchId      string  `json:"branch_id"`
	TarifId       string  `json:"tarif_id"`
	TypeStaffs    string  `json:"staff_type"`
	Name          string  `json:"name"`
	Balance       float64 `json:"balance"`
	Birthday_Date string  `json:"birth_date"`
	Age           int     `json:"age"`
	Loging        string  `json:"loging"`
	Password      string  `json:"password"`
	Created_at    string  `json:"created_at"`
	Updated_at    string  `json:"updated_at"`
}

type CreateStaff struct {
	BranchId      string  `json:"branch_id"`
	TarifId       string  `json:"tarif_id"`
	TypeStaffs    string  `json:"staff_type"`
	Name          string  `json:"name"`
	Balance       float64 `json:"balance"`
	Birthday_Date string  `json:"birth_date"`
	Age           int     `json:"age"`
	Loging        string  `json:"loging"`
	Password      string  `json:"password"`
}

type GetAllStaff struct {
	Staffs []Staff
	Count  int
}

type GetAllStaffRequest struct {
	Page  int
	Limit int
	Name  string
}

type UpdateBalanceRequest struct {
	TransactionType string
	SourceType      string
	ShopAssistent   StaffType
	Cashier         StaffType
	SaleId          string
	Text            string
}

type StaffType struct {
	StaffId string
	Amount  float64
}

type LoginReq struct {
	Loging   string `json:"loging"`
	Password string `json:"password"`
}
type LoginRes struct {
	Token string `json:"token"`
}

type RequestByUsername struct {
	Loging string
}

type ChangePassword struct {
	NewPassword string
	OldPassword string
}

type RequestByPassword struct {
	Id       string
	Password string
}
