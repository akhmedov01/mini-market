package memory

import (
	"context"
	"fmt"
	"gocrud/packages/helper"
	staff_service "staff/genproto"

	"github.com/google/uuid"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v4/pgxpool"
)

type staffTarifRepo struct {
	db *pgxpool.Pool
}

func NewStaffTariffRepo(db *pgxpool.Pool) *staffTarifRepo {

	return &staffTarifRepo{
		db: db,
	}

}

func (s *staffTarifRepo) CreateStaffTarif(ctx context.Context, req *staff_service.CreateTarif) (string, error) {

	id := uuid.NewString()

	query := `
	INSERT INTO 
		staff_tarifs(id,name,type,amount_for_cash,amount_for_card,founded_at) 
	VALUES($1,$2,$3,$4,$5,$6)`

	_, err := s.db.Exec(ctx, query,
		id,
		req.Name,
		req.Type,
		req.AmountForCash,
		req.AmountForCard,
		req.FoundedAt,
	)
	if err != nil {
		fmt.Println("error:", err.Error())
		return "", err
	}

	return id, nil
}

func (s *staffTarifRepo) UpdateStaffTarif(ctx context.Context, req *staff_service.Tarif) (string, error) {

	query := `
	UPDATE
		staff_tarifs
	SET
		name=$2,type=$3,amount_for_cash=$4,amount_for_card=$5,founded_at=$6,updated_at=NOW()
	WHERE
		id=$1`

	resp, err := s.db.Exec(ctx, query,
		req.Id,
		req.Name,
		req.Type,
		req.AmountForCash,
		req.AmountForCard,
		req.FoundedAt,
	)
	if err != nil {
		return "", err
	}
	if resp.RowsAffected() == 0 {
		return "", pgx.ErrNoRows
	}
	return "OK", nil
}

func (s *staffTarifRepo) GetStaffTarif(ctx context.Context, req *staff_service.IdRequest) (*staff_service.Tarif, error) {

	query := `SELECT 
	id,
	name,
	type,
	amount_for_cash,
	amount_for_card,
	founded_at::text,
	created_at::text,
	updated_at::text
	FROM staff_tarifs WHERE id=$1`

	resp := s.db.QueryRow(ctx, query, req.Id)

	var tarif staff_service.Tarif

	err := resp.Scan(
		&tarif.Id,
		&tarif.Name,
		&tarif.Type,
		&tarif.AmountForCash,
		&tarif.AmountForCard,
		&tarif.FoundedAt,
		&tarif.CreatedAt,
		&tarif.UpdatedAt,
	)

	if err != nil {
		fmt.Println("Error from Select")
		return &staff_service.Tarif{}, err
	}

	return &tarif, nil
}

func (s *staffTarifRepo) DeleteStaffTarif(ctx context.Context, req *staff_service.IdRequest) (string, error) {

	query := `DELETE FROM staff_tarifs WHERE id = $1`

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

func (s *staffTarifRepo) GetAllStaffTarif(ctx context.Context, req *staff_service.GetAllTarifRequest) (*staff_service.GetAllTarifResponse, error) {

	var (
		params  = make(map[string]interface{})
		filter  = "WHERE true"
		offsetQ = " OFFSET 0 "
		limit   = " LIMIT 10 "
		offset  = (req.Page - 1) * req.Limit
	)

	info := `
	SELECT
	id,
	name,
	type,
	amount_for_cash,
	amount_for_card,
	founded_at::text,
	created_at::text,
	updated_at::text
	FROM staff_tarifs
	`

	var count int

	cQ := `
	SELECT
		COUNT(*)
	FROM 
		staff_tarifs
	`
	if req.Name != "" {
		filter += ` AND name ILIKE '%' || @name || '%' `
		params["name"] = req.Name
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
		return &staff_service.GetAllTarifResponse{}, err
	}

	err = s.db.QueryRow(context.Background(), countQuery).Scan(&count)

	if err != nil {
		return &staff_service.GetAllTarifResponse{}, err
	}

	defer rows.Close()

	result := []*staff_service.Tarif{}

	for rows.Next() {

		var tarif staff_service.Tarif

		err := rows.Scan(
			&tarif.Id,
			&tarif.Name,
			&tarif.Type,
			&tarif.AmountForCash,
			&tarif.AmountForCard,
			&tarif.FoundedAt,
			&tarif.CreatedAt,
			&tarif.UpdatedAt,
		)
		if err != nil {
			return &staff_service.GetAllTarifResponse{}, err
		}

		result = append(result, &tarif)

	}

	return &staff_service.GetAllTarifResponse{Tarifs: result, Count: int64(count)}, nil
}

/* func (u *staffTarifRepo) read() ([]models.StaffTarif, error) {

	var (
		tarifs []models.StaffTarif
	)

	data, err := os.ReadFile(u.fileName)
	if err != nil {
		log.Printf("Error while Read data: %+v", err)
		return nil, err
	}

	err = json.Unmarshal(data, &tarifs)
	if err != nil {
		log.Printf("Error while Unmarshal data: %+v", err)
		return nil, err
	}

	return tarifs, nil
}

func (u *staffTarifRepo) write(tarifs []models.StaffTarif) error {

	body, err := json.Marshal(tarifs)
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
