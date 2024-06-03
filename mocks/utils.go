package mocks

import (
	"money-tracker/internal/domain"
	"testing"

	"github.com/stretchr/testify/mock"
)

type UtilsMock struct {
	mock.Mock
}

// GenerateRandomCode implements utils.Utils.
func (_m *UtilsMock) GenerateRandomCode(length int) string {
	panic("unimplemented")
}

// Slugify implements utils.Utils.
func (_m *UtilsMock) Slugify(text string) (string, error) {
	panic("unimplemented")
}

func (_m *UtilsMock) HashPassword(password string) (*string, *domain.Error) {
	args := _m.Called(password)
	if args.Get(0) == nil {
		return nil, args.Get(1).(*domain.Error)
	}
	return args.Get(0).(*string), nil
}

func NewUtilsMock(t *testing.T) *UtilsMock {
	mock := &UtilsMock{}
	mock.Mock.Test(t)
	return mock
}
