package auth

import (
	"money-tracker/internal/config"
	"money-tracker/internal/database/entity"
	"money-tracker/internal/domain"
	"money-tracker/internal/dto"
	refreshtoken "money-tracker/internal/modules/refresh_token"
	"money-tracker/internal/modules/user"
	"money-tracker/internal/utils"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type AuthService interface {
	generateNewToken(user *entity.User) (*dto.TokenResponse, *domain.Error)
	generateAndUpdateNewToken(user *entity.User, refreshTokenID int) (*dto.TokenResponse, *domain.Error)
	tokenGenerator(user *entity.User) (*dto.NewTokenDto, *domain.Error)
	LoginWithGoogle(code string) (*dto.TokenResponse, *domain.Error)
	RefreshToken(refreshToken string) (*dto.TokenResponse, *domain.Error)
	GetSelf(user *dto.ATClaims) (*dto.SelfResponse, *domain.Error)
	GenerateGoogleLoginUrl() string
}

type authService struct {
	refreshTokenService refreshtoken.RefreshTokenService
	config              config.Config
	userService         user.UserService
	utils               utils.Utils
}

// GetSelf implements AuthService.
func (a *authService) GetSelf(user *dto.ATClaims) (*dto.SelfResponse, *domain.Error) {
	data, err := a.userService.GetOneUserFromID(user.UserID)
	if err != nil {
		return nil, err
	}

	return &dto.SelfResponse{
		User: *data,
		Token: dto.ExpirationResponse{
			AccessTokenExpiresIn:  int(user.ExpiresAt.Time.Unix()),
			RefreshTokenExpiresIn: int(user.ExpiresAt.Time.Unix()),
		},
	}, nil
}

// GenerateGoogleLoginUrl implements AuthService.
func (a *authService) GenerateGoogleLoginUrl() string {
	return a.config.AuthCodeURL()
}

// RefreshToken implements AuthService.
func (a *authService) RefreshToken(refreshToken string) (*dto.TokenResponse, *domain.Error) {
	refresh, err := a.refreshTokenService.CheckRefreshTokenValidity(refreshToken)
	if err != nil {
		return nil, err
	}

	user, err := a.userService.GetOneUserFromID(int(refresh.UserID))

	if err != nil {
		return nil, err
	}

	newToken, err := a.generateAndUpdateNewToken(user, refresh.ID)

	if err != nil {
		return nil, err
	}

	return newToken, nil
}

// LoginWithGoogle implements AuthService.
func (a *authService) LoginWithGoogle(code string) (*dto.TokenResponse, *domain.Error) {
	googleUserData, userErr := a.config.ParseAccessTokenToUserData(code)

	if userErr != nil {
		return nil, userErr
	}
	checkedUser, existErr := a.userService.CheckEmail(googleUserData.Email)

	if existErr != nil {
		return nil, existErr
	}

	if checkedUser != nil {
		token, err := a.generateNewToken(checkedUser)
		if err != nil {
			return nil, err
		}
		return token, nil
	} else {
		newUser, newUserErr := a.userService.CreateUserFromGoogle(googleUserData)
		if newUserErr != nil {
			return nil, newUserErr
		}

		token, newTokenErr := a.generateNewToken(newUser)

		if newTokenErr != nil {
			return nil, newTokenErr
		}

		return token, nil
	}
}

// generateAndUpdateNewToken implements AuthService.
func (a *authService) generateAndUpdateNewToken(user *entity.User, refreshTokenID int) (*dto.TokenResponse, *domain.Error) {
	token, err := a.tokenGenerator(user)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	data, dataErr := a.refreshTokenService.UpdateRefreshTokenByRefreshTokenID(refreshTokenID, &entity.RefreshToken{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		ExpiredAt:    &token.RefreshTokenExpiresIn,
		UpdatedAt:    &now,
		UserID:       user.ID,
	})

	if dataErr != nil {
		return nil, dataErr
	}

	return &dto.TokenResponse{
		AccessToken:  data.AccessToken,
		TokenType:    "Bearer",
		RefreshToken: data.RefreshToken,
		ExpirationResponse: dto.ExpirationResponse{
			AccessTokenExpiresIn:  int(data.ExpiredAt.Unix()),
			RefreshTokenExpiresIn: int(data.ExpiredAt.Unix()),
		},
	}, nil

}

// tokenGenerator implements AuthService.
func (a *authService) tokenGenerator(user *entity.User) (*dto.NewTokenDto, *domain.Error) {
	secret := os.Getenv("JWT_SECRET")
	expiredTime := time.Now().Add(time.Minute * 60)
	claims := &dto.ATClaims{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiredTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, signError := token.SignedString([]byte(secret))

	if signError != nil {
		return nil, &domain.Error{
			Err:  signError,
			Code: fiber.StatusInternalServerError,
		}

	}

	expiredRefreshTime := time.Now().Add(time.Hour * 24 * 7)
	refreshToken := a.utils.GenerateRandomCode(32)

	return &dto.NewTokenDto{
		AccessToken:           accessToken,
		TokenType:             "Bearer",
		RefreshToken:          refreshToken,
		AccessTokenExpiresIn:  expiredTime,
		RefreshTokenExpiresIn: expiredRefreshTime,
	}, nil
}

// GenerateToken implements AuthService.
func (a *authService) generateNewToken(user *entity.User) (*dto.TokenResponse, *domain.Error) {
	newToken, tokenERr := a.tokenGenerator(user)

	if tokenERr != nil {
		return nil, tokenERr
	}

	_, err := a.refreshTokenService.GenerateRefreshToken(&entity.RefreshToken{
		UserID:       user.ID,
		AccessToken:  newToken.AccessToken,
		RefreshToken: newToken.RefreshToken,
		ExpiredAt:    &newToken.RefreshTokenExpiresIn,
	})

	if err != nil {
		return nil, err
	}

	return &dto.TokenResponse{
		AccessToken:  newToken.AccessToken,
		TokenType:    newToken.TokenType,
		RefreshToken: newToken.RefreshToken,
		ExpirationResponse: dto.ExpirationResponse{
			AccessTokenExpiresIn:  int(newToken.AccessTokenExpiresIn.Unix()),
			RefreshTokenExpiresIn: int(newToken.RefreshTokenExpiresIn.Unix()),
		},
	}, nil

}

func NewAuthService(
	refreshTokenService refreshtoken.RefreshTokenService,
	config config.Config,
	userService user.UserService,
	utils utils.Utils,
) AuthService {
	return &authService{
		refreshTokenService,
		config,
		userService,
		utils,
	}
}
