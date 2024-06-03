package mocks

import (
	"money-tracker/internal/database/entity"
	"money-tracker/internal/domain"

	"github.com/stretchr/testify/mock"
)

type UserRepositoryMock struct {
	mock.Mock
}

// CreateOne implements user.UserRepository.
func (_m *UserRepositoryMock) CreateOne(user entity.User) (*entity.User, *domain.Error) {
	panic("unimplemented")
}

func (_m *UserRepositoryMock) GetOneUserByID(id int) (*entity.User, *domain.Error) {
	args := _m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Get(1).(*domain.Error)
	}
	return args.Get(0).(*entity.User), nil
}

func (_m *UserRepositoryMock) GetOneUserByEmail(email string) (*entity.User, *domain.Error) {
	args := _m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Get(1).(*domain.Error)
	}
	return args.Get(0).(*entity.User), nil
}

func (_m *UserRepositoryMock) CreateUser(user entity.User) (*entity.User, *domain.Error) {
	args := _m.Called(user)
	return args.Get(0).(*entity.User), args.Get(1).(*domain.Error)
}

func NewUserRepositoryMock(t interface {
	mock.TestingT
}) *UserRepositoryMock {
	mock := &UserRepositoryMock{}
	mock.Mock.Test(t)
	return mock
}
