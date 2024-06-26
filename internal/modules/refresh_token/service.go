package refreshtoken

import (
	"errors"
	"money-tracker/internal/database/entity"
	"money-tracker/internal/domain"
	"time"

	"gorm.io/gorm"
)

type RefreshTokenService interface {
	GenerateRefreshToken(*entity.RefreshToken) (*entity.RefreshToken, *domain.Error)
	UpdateRefreshTokenByRefreshTokenID(refreshTokenID int, refresh *entity.RefreshToken) (*entity.RefreshToken, *domain.Error)
	CheckRefreshTokenValidity(string) (*entity.RefreshToken, *domain.Error)
	LogoutRefreshToken(string) *domain.Error
}
type refreshTokenService struct {
	refreshTokenRepository RefreshTokenRepository
}

// UpdateRefreshTokenByRefreshTokenID implements RefreshTokenService.
func (r *refreshTokenService) UpdateRefreshTokenByRefreshTokenID(refreshTokenID int, refresh *entity.RefreshToken) (*entity.RefreshToken, *domain.Error) {
	data, err := r.refreshTokenRepository.UpdateRefreshTokenByID(refreshTokenID, refresh)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// LogoutRefreshToken implements RefreshTokenService.
func (r *refreshTokenService) LogoutRefreshToken(accessToken string) *domain.Error {
	err := r.refreshTokenRepository.DeleteByAccessToken(accessToken)
	if err != nil {
		return err
	}
	return nil
}

// GetOneByRefreshToken implements RefreshTokenService.
func (r *refreshTokenService) CheckRefreshTokenValidity(token string) (*entity.RefreshToken, *domain.Error) {

	refresh, err := r.refreshTokenRepository.GetOneByRefreshToken(token)

	if err != nil {
		if errors.Is(err.Err, gorm.ErrRecordNotFound) {
			return nil, &domain.Error{
				Code: 401,
				Err:  errors.New("REFRESH_TOKEN_NOT_FOUND"),
			}
		}
		return nil, err
	}

	err = r.refreshTokenRepository.DeleteByRefreshToken(token)
	if err != nil {
		return nil, err
	}

	isExpired := time.Now().After(*refresh.ExpiredAt)

	if isExpired {
		return nil, &domain.Error{
			Code: 401,
			Err:  errors.New("REFRESH_TOKEN_EXPIRED"),
		}
	}

	return refresh, nil

}

// GenerateRefreshToken implements RefreshTokenService.
func (r *refreshTokenService) GenerateRefreshToken(refresh *entity.RefreshToken) (*entity.RefreshToken, *domain.Error) {
	token, err := r.refreshTokenRepository.CreateOne(*refresh)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func NewRefreshTokenService(refreshTokenRepository RefreshTokenRepository) RefreshTokenService {
	return &refreshTokenService{refreshTokenRepository}
}
