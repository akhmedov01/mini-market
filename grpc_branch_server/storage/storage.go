package storage

import (
	sale_service "branch/genproto"
	"context"

	"time"
)

type StoregeI interface {
	Branch() BranchesI
	// Staff() StaffsI
	// Sales() SalesI
	// Transaction() TransactionI
	// StaffTarif() StaffTarifI
}

type CacheI interface {
	Cache() RedisI
}

type RedisI interface {
	Create(context.Context, string, interface{}, time.Duration) error
	Get(context.Context, string, interface{}) (bool, error)
	Delete(context.Context, string) error
}

type BranchesI interface {
	Create(context.Context, *sale_service.CreateBranch) (string, error)
	Update(context.Context, *sale_service.Branch) (string, error)
	GetBranch(context.Context, *sale_service.IdReqRes) (*sale_service.Branch, error)
	GetAllBranch(context.Context, *sale_service.GetAllBranchRequest) (*sale_service.GetAllBranchResponse, error)
	DeleteBranch(context.Context, *sale_service.IdReqRes) (string, error)
}

// type StaffsI interface {
// 	CreateStaff(context.Context, models.CreateStaff) (string, error)
// 	UpdateStaff(context.Context, string, models.CreateStaff) (string, error)
// 	GetStaff(context.Context, models.IdRequest) (*models.Staff, error)
// 	GetAllStaff(context.Context, models.GetAllStaffRequest) (models.GetAllStaff, error)
// 	DeleteStaff(context.Context, models.IdRequest) (string, error)
// 	UpdateBalance(context.Context, models.UpdateBalanceRequest) (string, error)
// 	GetByUsername(context.Context, models.RequestByUsername) (models.Staff, error)
// 	ChangePassword(context.Context, models.RequestByPassword) (string, error)
// 	/* ChangeBalance(id string, amount float64) (string, error)
// 	GetMapOfStaffs() (map[string]models.Staff, error) */
// }

// type SalesI interface {
// 	CreateSale(context.Context, models.CreateSale) (string, error)
// 	UpdateSale(context.Context, string, models.CreateSale) (string, error)
// 	GetSale(context.Context, models.IdRequest) (*models.Sale, error)
// 	GetAllSale(context.Context, models.GetAllSaleRequest) (models.GetAllSale, error)
// 	DeleteSale(context.Context, models.IdRequest) (string, error)
// 	/* CancelSale(id string) (models.Sale, error)
// 	BranchTotal() (map[string]models.BranchTotalSumAndCount, error)
// 	GetSalesInDay() (map[string]map[string]float64, error) */
// }

// type TransactionI interface {
// 	CreateStaffTransaction(context.Context, models.CreateStaffTransaction) (string, error)
// 	UpdateStaffTransaction(context.Context, string, models.CreateStaffTransaction) (string, error)
// 	GetStaffTransaction(context.Context, models.IdRequest) (*models.StaffTransaction, error)
// 	GetAllStaffTransaction(context.Context, models.GetAllStaffTransactionRequest) (models.GetAllStaffTransaction, error)
// 	DeleteStaffTransaction(context.Context, models.IdRequest) (string, error)
// 	//FindErnedSum(dateFrom, dateTo string) (map[string]float64, error)
// }

// type StaffTarifI interface {
// 	CreateStaffTarif(context.Context, models.CreateStaffTarif) (string, error)
// 	UpdateStaffTarif(context.Context, string, models.CreateStaffTarif) (string, error)
// 	GetStaffTarif(context.Context, models.IdRequest) (*models.StaffTarif, error)
// 	GetAllStaffTarif(context.Context, models.GetAllStaffTarifRequest) (models.GetAllStaffTarif, error)
// 	DeleteStaffTarif(context.Context, models.IdRequest) (string, error)
// }
