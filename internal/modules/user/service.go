package user

import (
	"errors"
	"money-tracker/internal/database/entity"
	"money-tracker/internal/domain"
	"money-tracker/internal/dto"
	"money-tracker/internal/utils"

	"gorm.io/gorm"
)

type UserService interface {
	CheckEmail(email string) (*entity.User, *domain.Error)
	GetOneUserFromID(id int) (*entity.User, *domain.Error)
	CreateUserFromGoogle(body *dto.GoogleUserData) (*entity.User, *domain.Error)
}

type userService struct {
	userRepository UserRepository
	utils          utils.Utils
}

// GetOneUserFromID implements UserService.
func (u *userService) GetOneUserFromID(id int) (*entity.User, *domain.Error) {
	res, err := u.userRepository.GetOneUserByID(id)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// CreateUser implements UserService.
func (u *userService) CreateUserFromGoogle(body *dto.GoogleUserData) (*entity.User, *domain.Error) {
	hash, err := u.utils.HashPassword("")

	if err != nil {
		return nil, err
	}

	user, err := u.userRepository.CreateOne(entity.User{
		Name:              body.Name,
		Email:             body.Email,
		ProfilePictureUrl: body.Picture,
		Hash:              *hash,
	})

	if err != nil {
		return nil, err
	}

	return user, nil
}

// CheckEmail implements UserService.
func (u *userService) CheckEmail(email string) (*entity.User, *domain.Error) {
	user, err := u.userRepository.GetOneUserByEmail(email)
	if err != nil {
		if errors.Is(err.Err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

func NewUserService(userRepository UserRepository, utils utils.Utils) UserService {
	return &userService{
		userRepository,
		utils,
	}
}
