package refreshtoken

import (
	"money-tracker/internal/database/entity"
	"money-tracker/internal/domain"
)

type RefreshTokenService interface {
	GenerateRefreshToken(*entity.RefreshToken) (*entity.RefreshToken, *domain.Error)
}
type refreshTokenService struct {
	refreshTokenRepository RefreshTokenRepository
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
