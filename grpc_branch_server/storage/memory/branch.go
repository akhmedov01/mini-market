package memory

import (
	branch_service "branch/genproto"
	"branch/packages/helper"
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v4/pgxpool"
)

type branchRepo struct {
	db *pgxpool.Pool
}

func NewBranchRepo(db *pgxpool.Pool) *branchRepo {

	return &branchRepo{
		db: db,
	}

}

func (b *branchRepo) Create(ctx context.Context, req *branch_service.CreateBranch) (string, error) {

	id := uuid.NewString()

	query := `
	INSERT INTO 
		branches(id,name,address) 
	VALUES($1,$2,$3)`

	_, err := b.db.Exec(ctx, query,
		id,
		req.Name,
		req.Address,
	)
	if err != nil {
		fmt.Println("error:", err.Error())
		return "", err
	}

	return id, nil

}

func (b *branchRepo) Update(ctx context.Context, req *branch_service.Branch) (string, error) {

	query := `
	UPDATE branches
	SET name=$2,address=$3,updated_at=NOW()
	WHERE id=$1`

	resp, err := b.db.Exec(ctx, query,
		req.Id,
		req.Name,
		req.Address,
	)
	if err != nil {
		return "", err
	}
	if resp.RowsAffected() == 0 {
		return "", pgx.ErrNoRows
	}
	return "OK", nil
}

func (b *branchRepo) GetBranch(ctx context.Context, req *branch_service.IdReqRes) (*branch_service.Branch, error) {

	query := `
	SELECT
	id,
	name,
	address,
	created_at::text,
	updated_at::text
	FROM branches WHERE id = $1`

	resp := b.db.QueryRow(ctx, query, req.Id)

	var branch branch_service.Branch

	err := resp.Scan(
		&branch.Id,
		&branch.Name,
		&branch.Address,
		&branch.CreatedAt,
		&branch.UpdatedAt,
	)

	if err != nil {
		return &branch_service.Branch{}, err
	}

	return &branch, nil
}

func (b *branchRepo) GetAllBranch(ctx context.Context, req *branch_service.GetAllBranchRequest) (*branch_service.GetAllBranchResponse, error) {

	var (
		params  = make(map[string]interface{})
		filter  = " WHERE true"
		offsetQ = " OFFSET 0 "
		limit   = " LIMIT 10 "
		offset  = (req.GetPage() - 1) * req.GetLimit()
	)
	s := `
	SELECT
	id,
	name,
	address,
	created_at::text,
	updated_at::text
	FROM branches
	`
	var count int

	cQ := `
	SELECT
		COUNT(*)
	FROM
		branches
	`

	if req.GetName() != "" {
		filter += ` AND name ILIKE '%' || @name || '%' `
		params["name"] = req.GetName()
	}
	if req.Limit > 0 {
		limit = fmt.Sprintf(" LIMIT %d", req.GetLimit())
	}
	if offset > 0 {
		offsetQ = fmt.Sprintf(" OFFSET %d", offset)
	}

	query := s + filter + limit + offsetQ
	countQuery := cQ + filter

	q, pArr := helper.ReplaceQueryParams(query, params)
	rows, err := b.db.Query(ctx, q, pArr...)
	if err != nil {
		return &branch_service.GetAllBranchResponse{}, err
	}

	err = b.db.QueryRow(context.Background(), countQuery).Scan(&count)

	if err != nil {
		return &branch_service.GetAllBranchResponse{}, err
	}

	defer rows.Close()

	result := []*branch_service.Branch{}

	for rows.Next() {

		var branch branch_service.Branch

		err := rows.Scan(
			&branch.Id,
			&branch.Name,
			&branch.Address,
			&branch.CreatedAt,
			&branch.UpdatedAt,
		)
		if err != nil {
			return &branch_service.GetAllBranchResponse{}, err
		}

		result = append(result, &branch)

	}

	return &branch_service.GetAllBranchResponse{Branches: result, Count: int64(count)}, nil

}

func (b *branchRepo) DeleteBranch(ctx context.Context, req *branch_service.IdReqRes) (string, error) {

	query := `DELETE FROM branches WHERE id = $1`

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
