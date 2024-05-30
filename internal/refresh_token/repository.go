package refreshtoken

import (
	"money-tracker/internal/database/entity"
	"money-tracker/internal/domain"

	"gorm.io/gorm"
)

type RefreshTokenRepository interface {
	CreateOne(refresh entity.RefreshToken) (*entity.RefreshToken, *domain.Error)
	GetOneByRefreshToken(token string) (*entity.RefreshToken, *domain.Error)
	DeleteByRefreshToken(token string) *domain.Error
	DeleteByAccessToken(token string) *domain.Error
}
type refreshTokenRepository struct {
	db *gorm.DB
}

// GetOneByRefreshToken implements RefreshTokenRepository.
func (r *refreshTokenRepository) GetOneByRefreshToken(token string) (*entity.RefreshToken, *domain.Error) {
	var refresh entity.RefreshToken
	res := r.db.Raw("select * from refresh_token where refresh_token = ?", token).First(&refresh)

	if res.Error != nil {
		return nil, &domain.Error{
			Code: 500,
			Err:  res.Error,
		}
	}

	return &refresh, nil
}

// DeleteByAccessToken implements RefreshTokenRepository.
func (r *refreshTokenRepository) DeleteByAccessToken(token string) *domain.Error {
	res := r.db.Exec("delete from refresh_token where access_token = ?", token)

	if res.Error != nil {
		return &domain.Error{
			Code: 500,
			Err:  res.Error,
		}
	}

	return nil
}

// DeleteByToken implements RefreshTokenRepository.
func (r *refreshTokenRepository) DeleteByRefreshToken(token string) *domain.Error {
	res := r.db.Exec("delete from refresh_token where refresh_token = ?", token)

	if res.Error != nil {
		return &domain.Error{
			Code: 500,
			Err:  res.Error,
		}
	}

	return nil
}

// CreateOne implements RefreshTokenRepository.
func (r *refreshTokenRepository) CreateOne(refresh entity.RefreshToken) (*entity.RefreshToken, *domain.Error) {
	res := r.db.Raw("insert into refresh_token (access_token, refresh_token, user_id, expired_at) values (?, ?, ?, ?) returning *", refresh.AccessToken, refresh.RefreshToken, refresh.UserID, refresh.ExpiredAt).Scan(&refresh)

	if res.Error != nil {
		return nil, &domain.Error{
			Code: 500,
			Err:  res.Error,
		}
	}

	return &refresh, nil
}

func NewRefreshTokenRepository(db *gorm.DB) RefreshTokenRepository {
	return &refreshTokenRepository{db}
}
