package storage

import (
	"context"
	sale_service "sale/genproto"
)

type StoregeI interface {
	Sales() SalesI
	Transaction() TransactionI
	SaleProduct() SaleProductI
	BranchTransaction() BranchTransactionI
}

type SaleProductI interface {
	Create(context.Context, *sale_service.CreateSaleProduct) (string, error)
	Update(context.Context, *sale_service.SaleProduct) (string, error)
	Get(context.Context, *sale_service.IdRequest) (*sale_service.SaleProduct, error)
	GetAll(context.Context, *sale_service.GetAllSaleProductRequest) (*sale_service.GetAllSaleProductResponse, error)
	Delete(context.Context, *sale_service.IdRequest) (string, error)
}

type BranchTransactionI interface {
	Create(context.Context, *sale_service.CreateBranchTransaction) (string, error)
	Update(context.Context, *sale_service.BranchTransaction) (string, error)
	Get(context.Context, *sale_service.IdRequest) (*sale_service.BranchTransaction, error)
	GetAll(context.Context, *sale_service.GetAllBranchTransactionRequest) (*sale_service.GetAllBranchTransactionResponse, error)
	Delete(context.Context, *sale_service.IdRequest) (string, error)
}

type SalesI interface {
	CreateSale(context.Context, *sale_service.CreateSale) (string, error)
	UpdateSale(context.Context, *sale_service.Sale) (string, error)
	GetSale(context.Context, *sale_service.IdRequest) (*sale_service.Sale, error)
	GetAllSale(context.Context, *sale_service.GetAllSaleRequest) (*sale_service.GetAllSaleResponse, error)
	DeleteSale(context.Context, *sale_service.IdRequest) (string, error)
	/* CancelSale(id string) (models.Sale, error)
	BranchTotal() (map[string]models.BranchTotalSumAndCount, error)
	GetSalesInDay() (map[string]map[string]float64, error) */
}

type TransactionI interface {
	CreateStaffTransaction(context.Context, *sale_service.CreateTransaction) (string, error)
	UpdateStaffTransaction(context.Context, *sale_service.Transaction) (string, error)
	GetStaffTransaction(context.Context, *sale_service.IdRequest) (*sale_service.Transaction, error)
	GetAllStaffTransaction(context.Context, *sale_service.GetAllTransactionRequest) (*sale_service.GetAllTransactionResponse, error)
	DeleteStaffTransaction(context.Context, *sale_service.IdRequest) (string, error)
	//FindErnedSum(dateFrom, dateTo string) (map[string]float64, error)
}
