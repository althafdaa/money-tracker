package user

import (
	"errors"
	"money-tracker/internal/database/entity"
	"money-tracker/internal/domain"
	"money-tracker/internal/dto"
	"money-tracker/mocks"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type userInstance struct {
	userService UserService
	userRepo    *mocks.UserRepositoryMock
	utils       *mocks.UtilsMock
}

func createUserInstance(t *testing.T) userInstance {
	userRepo := mocks.NewUserRepositoryMock(t)
	utils := mocks.NewUtilsMock(t)
	userService := NewUserService(userRepo, utils)
	return userInstance{
		userRepo:    userRepo,
		userService: userService,
		utils:       utils,
	}
}

func TestGetOneUserFromID(t *testing.T) {
	t.Run("success get one user from id", func(t *testing.T) {
		instance := createUserInstance(t)

		instance.userRepo.On("GetOneUserByID", 1).Return(&entity.User{
			ID: 1,
		}, nil)

		data, err := instance.userService.GetOneUserFromID(1)

		assert.Nil(t, err)
		assert.Equal(t, 1, data.ID)
	})

	t.Run("failed get one user from id", func(t *testing.T) {
		instance := createUserInstance(t)

		instance.userRepo.On("GetOneUserByID", 1).Return(nil, &domain.Error{
			Code: 500,
			Err:  errors.New("error"),
		})

		_, err := instance.userService.GetOneUserFromID(1)

		assert.NotNil(t, err)
		assert.Equal(t, 500, err.Code)
		assert.Equal(t, errors.New("error"), err.Err)
	})

	t.Run("not found get one user from id", func(t *testing.T) {
		instance := createUserInstance(t)

		instance.userRepo.On("GetOneUserByID", 1).Return(nil, &domain.Error{
			Code: 500,
			Err:  gorm.ErrRecordNotFound,
		})

		_, err := instance.userService.GetOneUserFromID(1)

		assert.NotNil(t, err)
		assert.Equal(t, 500, err.Code)
		assert.Equal(t, gorm.ErrRecordNotFound, err.Err)

	})
}

func TestCheckEmail(t *testing.T) {
	t.Run("success check email", func(t *testing.T) {
		instance := createUserInstance(t)

		instance.userRepo.On("GetOneUserByEmail", "test@gmail.com").Return(&entity.User{
			Email: "test@gmail.com",
		}, nil)

		res, err := instance.userService.CheckEmail("test@gmail.com")
		assert.Nil(t, err)
		assert.Equal(t, "test@gmail.com", res.Email)
	})

	t.Run("email not found", func(t *testing.T) {
		instance := createUserInstance(t)

		instance.userRepo.On("GetOneUserByEmail", "test@mail.com").Return(nil, &domain.Error{
			Code: 500,
			Err:  gorm.ErrRecordNotFound,
		})

		data, err := instance.userService.CheckEmail("test@mail.com")

		assert.Nil(t, data)
		assert.Nil(t, err)
	})

	t.Run("database call error", func(t *testing.T) {
		instance := createUserInstance(t)

		instance.userRepo.On("GetOneUserByEmail", "test@mail.com").Return(nil, &domain.Error{
			Code: 500,
			Err:  errors.New("error"),
		})

		_, err := instance.userService.CheckEmail("test@mail.com")

		assert.NotNil(t, err)
		assert.Equal(t, 500, err.Code)
	})
}

func TestCreateUserFromGoogle(t *testing.T) {
	t.Run("success create user from google", func(t *testing.T) {
		instance := createUserInstance(t)
		hash := "hash"

		now := time.Now()

		repoBody := entity.User{
			Name:              "test",
			Email:             "test@gmail.com",
			ProfilePictureUrl: "test.jpg",
			Hash:              hash,
		}

		expected := &entity.User{
			Name:              "test",
			Email:             "test@gmail.com",
			ProfilePictureUrl: "test.jpg",
			Hash:              hash,
			CreatedAt:         &now,
			UpdatedAt:         &now,
		}

		instance.utils.On("HashPassword", "").Return(&hash, nil)
		instance.userRepo.On("CreateOne", repoBody).Return(expected, nil)

		data, err := instance.userService.CreateUserFromGoogle(&dto.GoogleUserData{
			Name:    "test",
			Email:   "test@gmail.com",
			Picture: "test.jpg",
		})

		assert.Nil(t, err)
		assert.Equal(t, expected, data)
	})

	t.Run("hashing failed", func(t *testing.T) {
		instance := createUserInstance(t)

		instance.utils.On("HashPassword", "").Return(nil, &domain.Error{
			Code: 500,
			Err:  errors.New("error"),
		})

		_, err := instance.userService.CreateUserFromGoogle(&dto.GoogleUserData{})

		assert.NotNil(t, err)
		assert.Equal(t, 500, err.Code)
	})

	t.Run("create user failed", func(t *testing.T) {
		instance := createUserInstance(t)
		hash := "hash"

		repoBody := entity.User{
			Name:              "test",
			Email:             "test@gmail.com",
			ProfilePictureUrl: "test.jpg",
			Hash:              hash,
		}

		instance.utils.On("HashPassword", "").Return(&hash, nil)

		instance.utils.On("HashPassword", "").Return(&hash, nil)
		instance.userRepo.On("CreateOne", repoBody).Return(nil, &domain.Error{
			Code: 500,
			Err:  errors.New("error"),
		})

		_, err := instance.userService.CreateUserFromGoogle(&dto.GoogleUserData{
			Name:    "test",
			Email:   "test@gmail.com",
			Picture: "test.jpg",
		})

		assert.NotNil(t, err)
		assert.Equal(t, 500, err.Code)
	})
}
