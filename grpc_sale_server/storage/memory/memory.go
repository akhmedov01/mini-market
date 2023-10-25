package memory

import (
	"context"
	"fmt"
	"sale/config"
	"sale/storage"

	"github.com/jackc/pgx/v4/pgxpool"
)

type store struct {
	db                *pgxpool.Pool
	sales             *saleRepo
	transaction       *transactionRepo
	saleProduct       *saleProductRepo
	branchTransaction *branchTransactionRepo
}

func NewStorage(ctx context.Context, cfg config.Config) (storage.StoregeI, error) {
	config, err := pgxpool.ParseConfig(
		fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
			cfg.PostgresUser,
			cfg.PostgresPassword,
			cfg.PostgresHost,
			cfg.PostgresPort,
			cfg.PostgresDatabase,
		),
	)

	if err != nil {
		fmt.Println("ParseConfig:", err.Error())
		return nil, err
	}

	config.MaxConns = cfg.PostgresMaxConnections
	pool, err := pgxpool.ConnectConfig(ctx, config)
	if err != nil {
		fmt.Println("ConnectConfig:", err.Error())
		return nil, err
	}
	return &store{
		db: pool,
	}, nil
}

func (s *store) Sales() storage.SalesI {
	if s.sales == nil {
		s.sales = NewSaleRepo(s.db)
	}
	return s.sales
}

func (s *store) Transaction() storage.TransactionI {
	if s.transaction == nil {
		s.transaction = NewTransactionRepo(s.db)
	}
	return s.transaction
}

func (s *store) SaleProduct() storage.SaleProductI {
	if s.saleProduct == nil {
		s.saleProduct = NewSaleProductRepo(s.db)
	}
	return s.saleProduct
}

func (s *store) BranchTransaction() storage.BranchTransactionI {
	if s.branchTransaction == nil {
		s.branchTransaction = NewBranchTransactionRepo(s.db)
	}
	return s.branchTransaction
}
