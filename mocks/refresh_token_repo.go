package mocks

import (
	"money-tracker/internal/database/entity"
	"money-tracker/internal/domain"

	"github.com/stretchr/testify/mock"
)

type RefreshTokenRepositoryMock struct {
	mock.Mock
}

func (_m *RefreshTokenRepositoryMock) CreateOne(refresh entity.RefreshToken) (*entity.RefreshToken, *domain.Error) {
	ret := _m.Called(refresh)

	var r0 *entity.RefreshToken
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*entity.RefreshToken)
	}

	var r1 *domain.Error
	if ret.Get(1) != nil {
		r1 = ret.Get(1).(*domain.Error)
	}

	return r0, r1
}

func NewRefreshTokenRepositoryMock(t interface {
	mock.TestingT
}) *RefreshTokenRepositoryMock {
	mock := &RefreshTokenRepositoryMock{}
	mock.Mock.Test(t)
	return mock
}
