package subcategory

type SubcategoryService interface{}
type subcategoryService struct {
	subcategoryRepository SubcategoryRepository
}

func NewSubcategoryService(subcategoryRepository subcategoryRepository) SubcategoryService {
	return &subcategoryService{
		subcategoryRepository,
	}
}
