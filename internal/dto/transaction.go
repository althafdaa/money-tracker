package dto

type GetAllQueryParams struct {
	Page          int    `query:"page" validate:"required, numeric"`
	Limit         int    `query:"limit" validate:"required, numeric, min=1, max=20"`
	Search        string `query:"search"`
	CategoryID    []int  `query:"category_id"`
	SubcategoryID []int  `query:"subcategory_id"`
}

type GetAllValueRepository struct {
	Offset        int
	Limit         int
	Search        string
	CategoryID    []int
	SubcategoryID []int
}
