package category

type CategoryService interface{}
type categoryService struct {
	categoryRepo CategoryRepository
}

func NewCategoryService(categoryRepo CategoryRepository) CategoryService {
	return &categoryService{categoryRepo}
}
