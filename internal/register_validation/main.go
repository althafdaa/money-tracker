package registervalidator

import (
	"money-tracker/internal/database/entity"

	"github.com/go-playground/validator/v10"
)

type ValidationRegister struct {
	validator *validator.Validate
}

func (v *ValidationRegister) RegisterCategoryType() error {
	err := v.validator.RegisterValidation("category_type", func(fl validator.FieldLevel) bool {
		return fl.Field().String() == string(entity.CategoryExpense) || fl.Field().String() == string(entity.CategoryIncome)
	})

	if err != nil {
		return err
	}

	return nil
}

func NewValidationRegiser(
	validator *validator.Validate,
) *ValidationRegister {
	return &ValidationRegister{
		validator,
	}
}
