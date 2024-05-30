package middleware

import (
	"errors"
	"money-tracker/internal/dto"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type AuthMiddleware struct{}

func (a *AuthMiddleware) Init(c *fiber.Ctx) error {
	tokenString := c.Get("Authorization")
	secret := os.Getenv("JWT_SECRET")
	if tokenString == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "UNAUTHORIZED",
			"code":  fiber.StatusUnauthorized,
		})
	}

	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	claims := &dto.ATClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrSignatureInvalid) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "TOKEN_INVALID",
				"code":  fiber.StatusUnauthorized,
			})
		}
		if errors.Is(err, jwt.ErrTokenExpired) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "TOKEN_EXPIRED",
				"code":  fiber.StatusUnauthorized,
			})
		}

		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
			"code":  fiber.StatusUnauthorized,
		})
	}

	if !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "TOKEN_INVALID",
			"code":  fiber.StatusUnauthorized,
		})
	}

	c.Locals("user", claims)

	return c.Next()
}

func NewAuthMiddleware() *AuthMiddleware {
	return &AuthMiddleware{}
}
