package memory

import (
	"context"
	"fmt"
	sale_service "sale/genproto"
	"sale/packages/helper"

	"github.com/google/uuid"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v4/pgxpool"
)

type transactionRepo struct {
	db *pgxpool.Pool
}

func NewTransactionRepo(db *pgxpool.Pool) *transactionRepo {

	return &transactionRepo{
		db: db,
	}

}

func (t *transactionRepo) CreateStaffTransaction(ctx context.Context, req *sale_service.CreateTransaction) (string, error) {

	id := uuid.NewString()

	query := `
	INSERT INTO
		staff_transactions(id,sale_id,staff_id,transaction_type,source_type,amount,information_about)
	VALUES($1,$2,$3,$4,$5,$6,$7)`

	_, err := t.db.Exec(ctx, query,
		id,
		req.SaleId,
		req.StaffId,
		req.TransactionType,
		req.SourceType,
		req.Amount,
		req.Text,
	)
	if err != nil {
		fmt.Println("error:", err.Error())
		return "", err
	}

	return id, nil

}

func (t *transactionRepo) UpdateStaffTransaction(ctx context.Context, req *sale_service.Transaction) (string, error) {

	query := `
	UPDATE
		staff_transactions
	SET
		sale_id=$2,staff_id=$3,transaction_type=$4,source_type=$5,amount=$6,information_about=$7,updated_at=NOW()
	WHERE
		id=$1`

	resp, err := t.db.Exec(ctx, query,
		req.Id,
		req.SaleId,
		req.StaffId,
		req.TransactionType,
		req.SourceType,
		req.Amount,
		req.Text,
	)
	if err != nil {
		return "", err
	}
	if resp.RowsAffected() == 0 {
		return "", pgx.ErrNoRows
	}
	return "OK", nil
}

func (t *transactionRepo) GetStaffTransaction(ctx context.Context, req *sale_service.IdRequest) (*sale_service.Transaction, error) {

	query := `SELECT
	id,
	sale_id,
	staff_id,
	transaction_type,
	source_type,
	amount,
	information_about,
	created_at::text,
	updated_at::text
	FROM staff_transactions WHERE id=$1`

	resp := t.db.QueryRow(ctx, query, req.Id)

	var transaction sale_service.Transaction

	err := resp.Scan(
		&transaction.Id,
		&transaction.SaleId,
		&transaction.StaffId,
		&transaction.TransactionType,
		&transaction.SourceType,
		&transaction.Amount,
		&transaction.Text,
		&transaction.CreatedAt,
		&transaction.UpdatedAt,
	)

	if err != nil {
		fmt.Println("Error from Select")
		return &sale_service.Transaction{}, err
	}

	return &transaction, nil
}

func (t *transactionRepo) DeleteStaffTransaction(ctx context.Context, req *sale_service.IdRequest) (string, error) {

	query := `DELETE FROM staff_transactions WHERE id = $1`

	resp, err := t.db.Exec(ctx, query,
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

func (t *transactionRepo) GetAllStaffTransaction(ctx context.Context, req *sale_service.GetAllTransactionRequest) (resp *sale_service.GetAllTransactionResponse, err error) {

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
	sale_id,
	staff_id,
	transaction_type,
	source_type,
	amount,
	information_about,
	created_at::text,
	updated_at::text
	FROM staff_transactions
	`
	var count int

	cQ := `
	SELECT
		COUNT(*)
	FROM 
		transactions
	`

	if req.TransactionType != "" {
		filter += ` AND transaction_type ILIKE '%' || @transaction_type || '%' `
		params["transaction_type"] = req.TransactionType
	}
	if req.Limit > 0 {
		limit = fmt.Sprintf(" LIMIT %d", req.Limit)
	}
	if offset > 0 {
		offsetQ = fmt.Sprintf(" OFFSET %d", offset)
	}

	query := info + filter + limit + offsetQ
	countQuery := cQ + filter

	q, pArr := helper.ReplaceQueryParams(query, params)
	rows, err := t.db.Query(ctx, q, pArr...)
	if err != nil {
		return &sale_service.GetAllTransactionResponse{}, err
	}

	err = t.db.QueryRow(context.Background(), countQuery).Scan(&count)

	if err != nil {
		return &sale_service.GetAllTransactionResponse{}, err
	}

	defer rows.Close()

	result := []*sale_service.Transaction{}

	for rows.Next() {

		var transaction sale_service.Transaction

		err := rows.Scan(
			&transaction.Id,
			&transaction.SaleId,
			&transaction.StaffId,
			&transaction.TransactionType,
			&transaction.SourceType,
			&transaction.Amount,
			&transaction.Text,
			&transaction.CreatedAt,
			&transaction.UpdatedAt,
		)
		if err != nil {
			return &sale_service.GetAllTransactionResponse{}, err
		}

		result = append(result, &transaction)

	}

	return &sale_service.GetAllTransactionResponse{Transactions: result, Count: int64(count)}, nil
}

/* func (t *transactionRepo) FindErnedSum(dateFrom, dateTo string) (map[string]float64, error) {

	counterMap := make(map[string]float64)

	transactions, err := t.read()

	if err != nil {
		return nil, err
	}

	parsingFromDate, err := time.Parse("2006-01-02", dateFrom)

	if err != nil {
		return nil, err
	}

	parsingToDate, err := time.Parse("2006-01-02", dateTo)

	if err != nil {
		return nil, err
	}

	for _, v := range transactions {

		parsingCreatedAt, err := time.Parse("2006-01-02", v.Created_at)

		if err != nil {
			return nil, err
		}

		if parsingToDate.After(parsingCreatedAt) && parsingFromDate.Before(parsingCreatedAt) {

			counterMap[v.StaffID] += v.Amount

		}

	}

	return counterMap, nil

}

func (u *transactionRepo) read() ([]models.StaffTransaction, error) {

	var (
		transactions []models.StaffTransaction
	)

	data, err := os.ReadFile(u.fileName)
	if err != nil {
		log.Printf("Error while Read data: %+v", err)
		return nil, err
	}

	err = json.Unmarshal(data, &transactions)
	if err != nil {
		log.Printf("Error while Unmarshal data: %+v", err)
		return nil, err
	}

	return transactions, nil
}

func (u *transactionRepo) write(transactions []models.StaffTransaction) error {

	body, err := json.Marshal(transactions)
	if err != nil {
		return err
	}

	err = os.WriteFile(u.fileName, body, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}
*/
