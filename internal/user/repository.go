package user

import (
	"errors"
	"money-tracker/internal/database/entity"
	"money-tracker/internal/domain"

	"gorm.io/gorm"
)

type UserRepository interface {
	GetOneUserByEmail(email string) (*entity.User, *domain.Error)
	GetOneUserByID(id int) (*entity.User, *domain.Error)
	CreateOne(user entity.User) (*entity.User, *domain.Error)
}
type userRepository struct {
	db *gorm.DB
}

// GetOneUserByID implements UserRepository.
func (u *userRepository) GetOneUserByID(id int) (*entity.User, *domain.Error) {
	var user *entity.User
	res := u.db.Raw("select * from user_data where email = ? limit 1", id).First(&user)

	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, &domain.Error{
			Code: 500,
			Err:  res.Error,
		}
	}

	return user, nil
}

// CreateOne implements UserRepository.
func (u *userRepository) CreateOne(user entity.User) (*entity.User, *domain.Error) {
	res := u.db.Raw("insert into user_data (email, name, profile_picture_url, hash) values (?, ?, ?, ?) returning *", user.Email, user.Name, user.ProfilePictureUrl, user.Hash).Scan(&user)

	if res.Error != nil {
		return nil, &domain.Error{
			Code: 500,
			Err:  res.Error,
		}
	}

	return &user, nil
}

// GetUserByEmail implements UserRepository.
func (u *userRepository) GetOneUserByEmail(email string) (*entity.User, *domain.Error) {
	var user *entity.User
	res := u.db.Raw("select * from user_data where email = ? limit 1", email).First(&user)

	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, &domain.Error{
			Code: 500,
			Err:  res.Error,
		}
	}

	return user, nil
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}
