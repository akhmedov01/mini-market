package models

type Branch struct {
	ID         string      `json:"id"`
	Name       string      `json:"name"`
	Address    string      `json:"address"`
	Created_At interface{} `json:"created_at"`
	Updated_at string      `json:"updated_at"`
}

type CreateBranch struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

type IdRequest struct {
	ID string
}

type GetAllBranch struct {
	Branches []Branch
	Count    int
}

type GetAllBranchRequest struct {
	Page  int
	Limit int
	Name  string
}

type BranchTotalSumAndCount struct {
	Sum   float64
	Count int
}
