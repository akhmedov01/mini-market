package storage

import (
	"context"
	staff_service "staff/genproto"

	"time"
)

type StoregeI interface {
	Staff() StaffsI
	// Sales() SalesI
	// Transaction() TransactionI
	StaffTarif() StaffTarifI
}

type CacheI interface {
	Cache() RedisI
}

type RedisI interface {
	Create(context.Context, string, interface{}, time.Duration) error
	Get(context.Context, string, interface{}) (bool, error)
	Delete(context.Context, string) error
}

type StaffsI interface {
	CreateStaff(context.Context, *staff_service.CreateStaff) (string, error)
	UpdateStaff(context.Context, *staff_service.Staff) (string, error)
	GetStaff(context.Context, *staff_service.IdRequest) (*staff_service.Staff, error)
	GetAllStaff(context.Context, *staff_service.GetAllStaffRequest) (*staff_service.GetAllStaffResponse, error)
	DeleteStaff(context.Context, *staff_service.IdRequest) (string, error)
	UpdateBalance(context.Context, *staff_service.UpdateBalanceRequest) (string, error)
	GetByUsername(context.Context, *staff_service.RequestByUsername) (*staff_service.Staff, error)
	ChangePassword(context.Context, *staff_service.RequestByPassword) (string, error)
	/* ChangeBalance(id string, amount float64) (string, error)
	GetMapOfStaffs() (map[string]models.Staff, error) */
}

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

type StaffTarifI interface {
	CreateStaffTarif(context.Context, *staff_service.CreateTarif) (string, error)
	UpdateStaffTarif(context.Context, *staff_service.Tarif) (string, error)
	GetStaffTarif(context.Context, *staff_service.IdRequest) (*staff_service.Tarif, error)
	GetAllStaffTarif(context.Context, *staff_service.GetAllTarifRequest) (*staff_service.GetAllTarifResponse, error)
	DeleteStaffTarif(context.Context, *staff_service.IdRequest) (string, error)
}
