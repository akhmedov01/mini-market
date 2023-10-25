package memory

import (
	"branch/config"
	"branch/storage"
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

type store struct {
	db       *pgxpool.Pool
	branches *branchRepo
	// staffs       *staffRepo
	// sales        *saleRepo
	// transactions *transactionRepo
	// staffTarifs  *staffTarifRepo
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

func (s *store) Branch() storage.BranchesI {
	if s.branches == nil {
		s.branches = NewBranchRepo(s.db)
	}
	return s.branches
}

// func (s *store) Staff() storage.StaffsI {
// 	if s.staffs == nil {
// 		s.staffs = NewStaffRepo(s.db)
// 	}
// 	return s.staffs
// }

// func (s *store) Sales() storage.SalesI {
// 	if s.sales == nil {
// 		s.sales = NewSaleRepo(s.db)
// 	}
// 	return s.sales
// }

// func (s *store) Transaction() storage.TransactionI {
// 	if s.transactions == nil {
// 		s.transactions = NewTransactionRepo(s.db)
// 	}
// 	return s.transactions
// }

// func (s *store) StaffTarif() storage.StaffTarifI {
// 	if s.staffTarifs == nil {
// 		s.staffTarifs = NewStaffTariffRepo(s.db)
// 	}
// 	return s.staffTarifs
// }
