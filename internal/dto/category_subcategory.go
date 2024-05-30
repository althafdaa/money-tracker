package dto

type CategoryBody struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type SubcategoryBody struct {
	Name       string `json:"name"`
	Slug       string `json:"slug"`
	CategoryID int    `json:"category_id"`
	UserID     int    `json:"user_id"`
}
