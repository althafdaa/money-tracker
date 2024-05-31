package entity

type CategoryType string

const (
	CategoryIncome  CategoryType = "income"
	CategoryExpense CategoryType = "expense"
)

type Category struct {
	ID   int          `json:"id"`
	Name string       `json:"name"`
	Slug string       `json:"slug"`
	Type CategoryType `json:"type"`
}

type PgCategory struct {
	ID   int
	Name string
	Slug string
	Type CategoryType
}
