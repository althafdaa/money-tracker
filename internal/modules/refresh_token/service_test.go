package refreshtoken

import (
	"money-tracker/internal/database/entity"
	"money-tracker/internal/domain"
	"money-tracker/mocks"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type refreshTokenInstance struct {
	refreshTokenRepo *mocks.RefreshTokenRepositoryMock
}

func createRefreshTokenInstance(t *testing.T) *refreshTokenInstance {
	refreshTokenRepo := mocks.NewRefreshTokenRepositoryMock(t)
	return &refreshTokenInstance{
		refreshTokenRepo,
	}
}

func TestGenerateRefreshToken(t *testing.T) {
	t.Run("Test GenerateRefreshToken", func(t *testing.T) {
		instance := createRefreshTokenInstance(t)
		expiredAt := time.Now().Add(time.Hour * 24)
		now := time.Now()
		body := entity.RefreshToken{
			AccessToken:  "random",
			UserID:       1,
			RefreshToken: "random",
			ExpiredAt:    &expiredAt,
		}

		expected := &entity.RefreshToken{
			AccessToken:  "random",
			UserID:       1,
			RefreshToken: "random",
			ExpiredAt:    &expiredAt,
			CreatedAt:    &now,
			UpdatedAt:    &now,
		}

		instance.refreshTokenRepo.On("CreateOne", body).Return(expected, &domain.Error{})
		res, err := instance.refreshTokenRepo.CreateOne(body)

		assert.NotNil(t, err)
		assert.Equal(t, expected, res)
	})

	t.Run("failed to generate refresh token", func(t *testing.T) {
		instance := createRefreshTokenInstance(t)
		expiredAt := time.Now().Add(time.Hour * 24)
		body := entity.RefreshToken{
			AccessToken:  "random",
			UserID:       1,
			RefreshToken: "random",
			ExpiredAt:    &expiredAt,
		}

		instance.refreshTokenRepo.On("CreateOne", body).Return(nil, &domain.Error{})
		_, err := instance.refreshTokenRepo.CreateOne(body)

		assert.NotNil(t, err)

	})
}
