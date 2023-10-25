package memory

import (
	"context"
	"fmt"
	"sale/packages/helper"

	sale_service "sale/genproto"

	"github.com/google/uuid"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v4/pgxpool"
)

type branchTransactionRepo struct {
	db *pgxpool.Pool
}

func NewBranchTransactionRepo(db *pgxpool.Pool) *branchTransactionRepo {

	return &branchTransactionRepo{
		db: db,
	}

}

func (b *branchTransactionRepo) Create(ctx context.Context, req *sale_service.CreateBranchTransaction) (string, error) {

	var id = uuid.NewString()

	query := `
	INSERT INTO
	branch_product_transactions (id,branch_id,staff_id,product_id,type,quantity,price)
	VALUES ($1,$2,$3,$4,$5,$6,$7)`

	_, err := b.db.Exec(ctx, query,
		id,
		req.BranchId,
		req.StaffId,
		req.ProductId,
		req.Type,
		req.Quantity,
		req.Price,
	)
	if err != nil {
		fmt.Println("error:", err.Error())
		return "", err
	}

	return id, nil
}

func (b *branchTransactionRepo) Update(ctx context.Context, req *sale_service.BranchTransaction) (string, error) {

	query := `
	UPDATE
		branch_product_transactions
	SET
		branch_id=$2,staff_id=$3,product_id=$4,type=$5,quantity=$6,price=$7,updated_at=NOW()
	WHERE
		id=$1`

	resp, err := b.db.Exec(ctx, query,
		req.Id,
		req.BranchId,
		req.StaffId,
		req.ProductId,
		req.Type,
		req.Quantity,
		req.Price,
	)

	if err != nil {
		return "", err
	}
	if resp.RowsAffected() == 0 {
		return "", pgx.ErrNoRows
	}
	return "OK", nil
}

func (b *branchTransactionRepo) Get(ctx context.Context, req *sale_service.IdRequest) (*sale_service.BranchTransaction, error) {

	query := `
	SELECT
	id,
	branch_id,
	staff_id,
	product_id,
	type,
	quantity,
	price,
	created_at::text,
	updated_at::text
	FROM branch_product_transactions
	WHERE id=$1`

	resp := b.db.QueryRow(ctx, query, req.Id)

	var branchTransaction sale_service.BranchTransaction

	err := resp.Scan(
		&branchTransaction.Id,
		&branchTransaction.BranchId,
		&branchTransaction.StaffId,
		&branchTransaction.ProductId,
		&branchTransaction.Type,
		&branchTransaction.Quantity,
		&branchTransaction.Price,
		&branchTransaction.CreatedAt,
		&branchTransaction.UpdatedAt,
	)

	if err != nil {
		return &sale_service.BranchTransaction{}, err
	}

	return &branchTransaction, nil
}

func (b *branchTransactionRepo) Delete(ctx context.Context, req *sale_service.IdRequest) (string, error) {

	query := `DELETE FROM branch_product_transactions WHERE id = $1`

	resp, err := b.db.Exec(ctx, query,
		req.Id,
	)

	if err != nil {
		return "", err
	}
	if resp.RowsAffected() == 0 {
		return "", pgx.ErrNoRows
	}

	return "Deleted suc", nil
}

func (b *branchTransactionRepo) GetAll(ctx context.Context, req *sale_service.GetAllBranchTransactionRequest) (resp *sale_service.GetAllBranchTransactionResponse, err error) {

	var (
		params  = make(map[string]interface{})
		filter  = " WHERE true"
		offsetQ = " OFFSET 0 "
		limit   = " LIMIT 10 "
		offset  = (req.GetPage() - 1) * req.GetLimit()
	)

	info := `
	SELECT
	id,
	branch_id,
	staff_id,
	product_id,
	type,
	quantity,
	price,
	created_at::text,
	updated_at::text
	FROM branch_product_transactions
	`

	var count int

	cQ := `
	SELECT
		COUNT(*)
	FROM 
		branch_product_transactions
	`

	// if req.ClientName != "" {
	// 	filter += ` AND client_name ILIKE '%' || @client_name || '%' `
	// 	params["client_name"] = req.ClientName
	// }

	if req.Limit > 0 {
		limit = fmt.Sprintf(" LIMIT %d", req.Limit)
	}
	if offset > 0 {
		offsetQ = fmt.Sprintf(" OFFSET %d", offset)
	}

	query := info + filter + limit + offsetQ
	countQuery := cQ + filter

	q, pArr := helper.ReplaceQueryParams(query, params)
	rows, err := b.db.Query(ctx, q, pArr...)
	if err != nil {
		return &sale_service.GetAllBranchTransactionResponse{}, err
	}

	err = b.db.QueryRow(context.Background(), countQuery).Scan(&count)

	if err != nil {
		return &sale_service.GetAllBranchTransactionResponse{}, err
	}

	defer rows.Close()

	result := []*sale_service.BranchTransaction{}

	for rows.Next() {

		var branchTransaction sale_service.BranchTransaction

		err := rows.Scan(
			&branchTransaction.Id,
			&branchTransaction.BranchId,
			&branchTransaction.StaffId,
			&branchTransaction.ProductId,
			&branchTransaction.Type,
			&branchTransaction.Quantity,
			&branchTransaction.Price,
			&branchTransaction.CreatedAt,
			&branchTransaction.UpdatedAt,
		)
		if err != nil {
			return &sale_service.GetAllBranchTransactionResponse{}, err
		}

		result = append(result, &branchTransaction)

	}

	return &sale_service.GetAllBranchTransactionResponse{BranchTransactions: result, Count: int64(count)}, nil

}
