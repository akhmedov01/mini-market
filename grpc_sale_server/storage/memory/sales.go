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

type saleRepo struct {
	db *pgxpool.Pool
}

func NewSaleRepo(db *pgxpool.Pool) *saleRepo {

	return &saleRepo{
		db: db,
	}

}

func (s *saleRepo) CreateSale(ctx context.Context, req *sale_service.CreateSale) (string, error) {

	var id = uuid.NewString()

	query := `
	INSERT INTO
	sales (id,branch_id,shop_assistent_id,cashier_id,payment_type,status,client_name,price)
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`

	if req.ShopAssistentId == "" {

		_, err := s.db.Exec(ctx, query,
			id,
			req.BranchId,
			nil,
			req.CashierId,
			req.PaymentType,
			req.Status,
			req.ClientName,
			req.Price,
		)
		if err != nil {
			fmt.Println("error:", err.Error())
			return "", err
		}

	}
	_, err := s.db.Exec(ctx, query,
		id,
		req.BranchId,
		req.ShopAssistentId,
		req.CashierId,
		req.PaymentType,
		req.Status,
		req.ClientName,
		req.Price,
	)
	if err != nil {
		fmt.Println("error:", err.Error())
		return "", err
	}

	return id, nil
}

func (s *saleRepo) UpdateSale(ctx context.Context, req *sale_service.Sale) (string, error) {

	query := `
	UPDATE
		sales
	SET
		branch_id=$2,shop_assistent_id=$3,cashier_id=$4,payment_type=$5,status=$6,client_name=$7,price=$8,updated_at=NOW()
	WHERE
		id=$1`

	resp, err := s.db.Exec(ctx, query,
		req.Id,
		req.BranchId,
		req.ShopAssistentId,
		req.CashierId,
		req.PaymentType,
		req.Status,
		req.ClientName,
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

func (s *saleRepo) GetSale(ctx context.Context, req *sale_service.IdRequest) (*sale_service.Sale, error) {

	query := `
	SELECT
	id,
	branch_id,
	shop_assistent_id,
	cashier_id,
	payment_type,
	status,
	client_name,
	price,
	created_at::text,
	updated_at::text
	FROM sales
	WHERE id=$1`

	resp := s.db.QueryRow(ctx, query, req.Id)

	var sale sale_service.Sale

	err := resp.Scan(
		&sale.Id,
		&sale.BranchId,
		&sale.ShopAssistentId,
		&sale.CashierId,
		&sale.PaymentType,
		&sale.Status,
		&sale.ClientName,
		&sale.Price,
		&sale.CreatedAt,
		&sale.UpdatedAt,
	)

	if err != nil {
		return &sale_service.Sale{}, err
	}

	return &sale, nil
}

func (s *saleRepo) DeleteSale(ctx context.Context, req *sale_service.IdRequest) (string, error) {

	query := `DELETE FROM sales WHERE id = $1`

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

func (s *saleRepo) GetAllSale(ctx context.Context, req *sale_service.GetAllSaleRequest) (resp *sale_service.GetAllSaleResponse, err error) {

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
	shop_assistent_id,
	cashier_id,
	payment_type,
	status,
	client_name,
	price,
	created_at::text,
	updated_at::text
	FROM sales
	`

	var count int

	cQ := `
	SELECT
		COUNT(*)
	FROM 
		sale
	`

	if req.ClientName != "" {
		filter += ` AND client_name ILIKE '%' || @client_name || '%' `
		params["client_name"] = req.ClientName
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
	rows, err := s.db.Query(ctx, q, pArr...)
	if err != nil {
		return &sale_service.GetAllSaleResponse{}, err
	}

	err = s.db.QueryRow(context.Background(), countQuery).Scan(&count)

	if err != nil {
		return &sale_service.GetAllSaleResponse{}, err
	}

	defer rows.Close()

	result := []*sale_service.Sale{}

	for rows.Next() {

		var sale sale_service.Sale

		err := rows.Scan(
			&sale.Id,
			&sale.BranchId,
			&sale.ShopAssistentId,
			&sale.CashierId,
			&sale.PaymentType,
			&sale.Status,
			&sale.ClientName,
			&sale.Price,
			&sale.CreatedAt,
			&sale.UpdatedAt,
		)
		if err != nil {
			return &sale_service.GetAllSaleResponse{}, err
		}

		result = append(result, &sale)

	}

	return &sale_service.GetAllSaleResponse{Sales: result, Count: int64(count)}, nil

}

/* func (s *saleRepo) CancelSale(id string) (models.Sale, error) {

	sales, err := s.read()

	if err != nil {
		return models.Sale{}, err
	}

	for i, v := range sales {
		if v.ID == id {

			sales[i].Status = 2

			err = s.write(sales)

			if err != nil {
				return models.Sale{}, err
			}

			return v, nil

		}
	}

	return models.Sale{}, err

}

func (s *saleRepo) BranchTotal() (map[string]models.BranchTotalSumAndCount, error) {

	sales, err := s.read()

	if err != nil {
		return nil, err
	}

	resultMap := make(map[string]models.BranchTotalSumAndCount)

	var sum float64
	var count int

	for _, v := range sales {

		if _, ok := resultMap[v.BranchId]; ok {

			sum = resultMap[v.BranchId].Sum
			count = resultMap[v.BranchId].Count

		} else {

			sum, count = 0.0, 0

		}

		if v.Status == 1 {

			resultMap[v.BranchId] = models.BranchTotalSumAndCount{

				Sum:   sum,
				Count: count,
			}

		}

	}

	return resultMap, nil

}

func (s *saleRepo) GetSalesInDay() (map[string]map[string]float64, error) {

	sales, err := s.read()

	if err != nil {
		return nil, err
	}

	salesInDayMap := make(map[string]map[string]float64)

	for _, v := range sales {

		salesInDayMap[v.BranchId][v.Created_at] += v.Price

	}

	return salesInDayMap, nil

}

func (u *saleRepo) read() ([]models.Sale, error) {

	var (
		sales []models.Sale
	)

	data, err := os.ReadFile(u.fileName)
	if err != nil {
		log.Printf("Error while Read data: %+v", err)
		return nil, err
	}

	err = json.Unmarshal(data, &sales)
	if err != nil {
		log.Printf("Error while Unmarshal data: %+v", err)
		return nil, err
	}

	return sales, nil
}

func (u *saleRepo) write(sales []models.Sale) error {

	body, err := json.Marshal(sales)
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
