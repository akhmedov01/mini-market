package memory

import (
	"context"
	"fmt"
	staff_service "staff/genproto"
	"staff/packages/helper"

	"github.com/google/uuid"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v4/pgxpool"
)

type staffRepo struct {
	db *pgxpool.Pool
}

func NewStaffRepo(db *pgxpool.Pool) *staffRepo {

	return &staffRepo{
		db: db,
	}

}

func (s *staffRepo) CreateStaff(ctx context.Context, req *staff_service.CreateStaff) (string, error) {

	id := uuid.NewString()

	//hashPas, err := helper.GeneratePasswordHash(req.Password)

	// if err != nil {
	// 	return "Error while generete password hash", err
	// }

	//req.Password = string(hashPas)

	query := `
	INSERT INTO 
		staffs(id,branch_id,tarif_id,staff_type,name,balance,birth_date,age,loging,password) 
	VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`

	_, err := s.db.Exec(ctx, query,
		id,
		req.BranchId,
		req.TarifId,
		req.TypeStaff,
		req.Name,
		req.Balance,
		req.BirthDate,
		req.Age,
		req.Login,
		req.Password,
	)
	if err != nil {
		fmt.Println("error:", err.Error())
		return "", err
	}

	return id, nil
}

func (s *staffRepo) UpdateStaff(ctx context.Context, req *staff_service.Staff) (string, error) {

	query := `
	UPDATE
		staffs
	SET
		branch_id=$2,tarif_id=$3,staff_type=$4,name=$5,balance=$6,birth_date=$7,age=$8,updated_at=NOW()
	WHERE
		id=$1`

	resp, err := s.db.Exec(ctx, query,
		req.Id,
		req.BranchId,
		req.TarifId,
		req.TypeStaff,
		req.Name,
		req.Balance,
		req.BirthDate,
		req.Age,
	)
	if err != nil {
		return "", err
	}
	if resp.RowsAffected() == 0 {
		return "", pgx.ErrNoRows
	}
	return "OK", nil

}

func (s *staffRepo) GetStaff(ctx context.Context, req *staff_service.IdRequest) (*staff_service.Staff, error) {

	query := `SELECT 
	id,
	branch_id,
	tarif_id,
	staff_type,
	name,
	balance,
	birth_date::text,
	age,
	loging,
	password,
	created_at::text,
	updated_at::text
	FROM staffs WHERE id=$1`

	resp := s.db.QueryRow(ctx, query, req.Id)

	var staff staff_service.Staff

	err := resp.Scan(
		&staff.Id,
		&staff.BranchId,
		&staff.TarifId,
		&staff.TypeStaff,
		&staff.Name,
		&staff.Balance,
		&staff.BirthDate,
		&staff.Age,
		&staff.Login,
		&staff.Password,
		&staff.CreatedAt,
		&staff.UpdatedAt,
	)

	if err != nil {
		fmt.Println("Error from Select")
		return &staff_service.Staff{}, err
	}

	return &staff, nil

}

func (s *staffRepo) DeleteStaff(ctx context.Context, req *staff_service.IdRequest) (string, error) {

	query := `DELETE FROM staffs WHERE id = $1`

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

func (s *staffRepo) GetAllStaff(ctx context.Context, req *staff_service.GetAllStaffRequest) (resp *staff_service.GetAllStaffResponse, err error) {

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
	tarif_id,
	staff_type,
	name,
	balance,
	birth_date::text,
	age,
	loging,
	password,
	created_at::text,
	updated_at::text
	FROM staffs
	`

	var count int

	cQ := `
	SELECT
		COUNT(*)
	FROM 
		staff_tarifs
	`

	if req.GetName() != "" {
		filter += ` AND name ILIKE '%' || @name || '%' `
		params["name"] = req.GetName()
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
		return &staff_service.GetAllStaffResponse{}, err
	}

	err = s.db.QueryRow(context.Background(), countQuery).Scan(&count)

	if err != nil {
		return &staff_service.GetAllStaffResponse{}, err
	}

	defer rows.Close()

	result := []*staff_service.Staff{}

	for rows.Next() {

		var staff staff_service.Staff

		err := rows.Scan(
			&staff.Id,
			&staff.BranchId,
			&staff.TarifId,
			&staff.TypeStaff,
			&staff.Name,
			&staff.Balance,
			&staff.BirthDate,
			&staff.Age,
			&staff.Login,
			&staff.Password,
			&staff.CreatedAt,
			&staff.UpdatedAt,
		)
		if err != nil {
			return &staff_service.GetAllStaffResponse{}, err
		}

		result = append(result, &staff)

	}

	return &staff_service.GetAllStaffResponse{Staffs: result, Count: int64(count)}, nil

}

func (s *staffRepo) UpdateBalance(ctx context.Context, req *staff_service.UpdateBalanceRequest) (string, error) {

	tr, err := s.db.Begin(ctx)

	defer func() {
		if err != nil {
			tr.Rollback(ctx)
		} else {
			tr.Commit(ctx)
		}
	}()

	updateQuery1 := `
	UPDATE staffs
	SET balance=balance+$2, updated_at=NOW()
	WHERE id=$1`

	if req.TransactionType == "Withdraw" {
		req.Cashier.Amount = -(req.Cashier.Amount)
		req.ShopAssistent.Amount = -(req.ShopAssistent.Amount)
	}

	_, err = tr.Exec(ctx, updateQuery1,
		req.Cashier.StaffId,
		req.Cashier.Amount,
	)

	insertQuery1 := `
	INSERT INTO 
		staff_transactions(id,sale_id,staff_id,transaction_type,source_type,amount,information_about) 
	VALUES($1,$2,$3,$4,$5,$6,$7)`

	_, err = tr.Exec(ctx, insertQuery1,
		uuid.NewString(),
		req.SaleId,
		req.Cashier.StaffId,
		req.TransactionType,
		req.SourceType,
		req.Cashier.Amount,
		req.Text,
	)
	if err != nil {
		fmt.Println("error:", err.Error())
		return "", err
	}

	if req.ShopAssistent.StaffId != "" {
		updateQuery2 := `
		UPDATE staffs
		SET balance=balance+$2,updated_at=NOW()
		WHERE id=$1`

		_, err = tr.Exec(ctx, updateQuery2,
			req.ShopAssistent.StaffId,
			req.ShopAssistent.Amount,
		)

		insertQuery2 := `
		INSERT INTO 
			staff_transactions(id,sale_id,staff_id,transaction_type,source_type,amount,information_about) 
		VALUES($1,$2,$3,$4,$5,$6,$7)`

		_, err = tr.Exec(ctx, insertQuery2,
			uuid.NewString(),
			req.SaleId,
			req.ShopAssistent.StaffId,
			req.TransactionType,
			req.SourceType,
			req.ShopAssistent.Amount,
			req.Text,
		)
		if err != nil {
			fmt.Println("error:", err.Error())
			return "", err
		}
	}

	return "", nil
}

func (s *staffRepo) GetByUsername(ctx context.Context, req *staff_service.RequestByUsername) (*staff_service.Staff, error) {

	selectQ := `SELECT 
	id,
	branch_id,
	tarif_id,
	staff_type,
	name,
	balance,
	birth_date::text,
	age,
	loging,
	password,
	created_at::text,
	updated_at::text
	FROM staffs WHERE loging=$1`

	staff := &staff_service.Staff{}
	err := s.db.QueryRow(ctx, selectQ, req.Login).Scan(
		&staff.Id,
		&staff.BranchId,
		&staff.TarifId,
		&staff.TypeStaff,
		&staff.Name,
		&staff.Balance,
		&staff.BirthDate,
		&staff.Age,
		&staff.Login,
		&staff.Password,
		&staff.CreatedAt,
		&staff.UpdatedAt,
	)
	if err != nil {
		return staff, err
	}
	return staff, nil
}

func (s *staffRepo) ChangePassword(ctx context.Context, req *staff_service.RequestByPassword) (string, error) {

	query := `
	UPDATE
		staffs
	SET
		password=$2,updated_at=NOW()
	WHERE
		id=$1`

	resp, err := s.db.Exec(ctx, query,
		req.Id,
		req.Password,
	)
	if err != nil {
		return "", err
	}
	if resp.RowsAffected() == 0 {
		return "", pgx.ErrNoRows
	}
	return "OK", nil

}

/* func (s *staffRepo) ChangeBalance(id string, amount float64) (string, error) {

	staffs, err := s.read()

	if err != nil {
		return "error", err
	}

	for i, v := range staffs {
		if v.ID == id {

			staffs[i].Balance += amount

		}
	}

	err = s.write(staffs)

	if err != nil {
		return "error", err
	}

	return "Balance suc changed", nil

}

func (s *staffRepo) GetMapOfStaffs() (map[string]models.Staff, error) {

	staffs, err := s.read()

	if err != nil {
		return nil, err
	}

	staffsMap := make(map[string]models.Staff)

	for _, v := range staffs {

		staffsMap[v.ID] = v

	}

	return staffsMap, nil

}

func (u *staffRepo) read() ([]models.Staff, error) {

	var (
		staffs []models.Staff
	)

	data, err := os.ReadFile(u.fileName)
	if err != nil {
		log.Printf("Error while Read data: %+v", err)
		return nil, err
	}

	err = json.Unmarshal(data, &staffs)
	if err != nil {
		log.Printf("Error while Unmarshal data: %+v", err)
		return nil, err
	}

	return staffs, nil
}

func (u *staffRepo) write(staffs []models.Staff) error {

	body, err := json.Marshal(staffs)
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
