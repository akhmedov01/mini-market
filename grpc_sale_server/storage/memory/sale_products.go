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

type saleProductRepo struct {
	db *pgxpool.Pool
}

func NewSaleProductRepo(db *pgxpool.Pool) *saleProductRepo {

	return &saleProductRepo{
		db: db,
	}

}

func (s *saleProductRepo) Create(ctx context.Context, req *sale_service.CreateSaleProduct) (string, error) {

	var id = uuid.NewString()

	query := `
	INSERT INTO
	sale_products (id,sale_id,product_id,quantity,price)
	VALUES ($1,$2,$3,$4,$5)`

	_, err := s.db.Exec(ctx, query,
		id,
		req.SaleId,
		req.ProductId,
		req.Quantity,
		req.Price,
	)
	if err != nil {
		fmt.Println("error:", err.Error())
		return "", err
	}

	return id, nil
}

func (s *saleProductRepo) Update(ctx context.Context, req *sale_service.SaleProduct) (string, error) {

	query := `
	UPDATE
		sale_products
	SET
		sale_id=$2,product_id=$3,quantity=$4,price=$5,updated_at=NOW()
	WHERE
		id=$1`

	resp, err := s.db.Exec(ctx, query,
		req.Id,
		req.SaleId,
		req.ProductId,
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

func (s *saleProductRepo) Get(ctx context.Context, req *sale_service.IdRequest) (*sale_service.SaleProduct, error) {

	query := `
	SELECT
	id,
	sale_id,
	product_id,
	quantity,
	price,
	created_at::text,
	updated_at::text
	FROM sale_products
	WHERE id=$1`

	resp := s.db.QueryRow(ctx, query, req.Id)

	var saleProducts sale_service.SaleProduct

	err := resp.Scan(
		&saleProducts.Id,
		&saleProducts.SaleId,
		&saleProducts.ProductId,
		&saleProducts.Quantity,
		&saleProducts.Price,
		&saleProducts.CreatedAt,
		&saleProducts.UpdatedAt,
	)

	if err != nil {
		return &sale_service.SaleProduct{}, err
	}

	return &saleProducts, nil
}

func (s *saleProductRepo) Delete(ctx context.Context, req *sale_service.IdRequest) (string, error) {

	query := `DELETE FROM sale_products WHERE id = $1`

	resp, err := s.db.Exec(ctx, query,
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

func (s *saleProductRepo) GetAll(ctx context.Context, req *sale_service.GetAllSaleProductRequest) (resp *sale_service.GetAllSaleProductResponse, err error) {

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
	product_id,
	quantity,
	price,
	created_at::text,
	updated_at::text
	FROM sale_products
	`

	var count int

	cQ := `
	SELECT
		COUNT(*)
	FROM 
		sale_products
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
	rows, err := s.db.Query(ctx, q, pArr...)
	if err != nil {
		return &sale_service.GetAllSaleProductResponse{}, err
	}

	err = s.db.QueryRow(context.Background(), countQuery).Scan(&count)

	if err != nil {
		return &sale_service.GetAllSaleProductResponse{}, err
	}

	defer rows.Close()

	result := []*sale_service.SaleProduct{}

	for rows.Next() {

		var saleProducts sale_service.SaleProduct

		err := rows.Scan(
			&saleProducts.Id,
			&saleProducts.SaleId,
			&saleProducts.ProductId,
			&saleProducts.Quantity,
			&saleProducts.Price,
			&saleProducts.CreatedAt,
			&saleProducts.UpdatedAt,
		)
		if err != nil {
			return &sale_service.GetAllSaleProductResponse{}, err
		}

		result = append(result, &saleProducts)

	}

	return &sale_service.GetAllSaleProductResponse{SaleProducts: result, Count: int64(count)}, nil

}
