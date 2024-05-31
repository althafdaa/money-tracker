package entity

type Subcategory struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Slug       string `json:"slug"`
	CategoryID int    `json:"category_id"`
	UserID     int    `json:"user_id"`
}

type PgSubcategory struct {
	ID       int
	Name     string
	Slug     string
	Category PgCategory
	User     PgUser
}
