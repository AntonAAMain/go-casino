package auth

import (
	"errors"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type Claims struct {
	UserID uint   `json:"user_id"`
	Role   string `jsom:"role"`
	jwt.RegisteredClaims
}

var JwtSecret = []byte("your_secret_key")

func Authorize(token string, c echo.Context) (int, string, error) {

	user_id, role, err := ParseToken(token, string(JwtSecret))

	if err != nil {
		return 0, "", err
	}

	return int(user_id), role, nil
}

func ParseToken(tokenStr, secret string) (uint, string, error) {
	token, err := jwt.ParseWithClaims(
		tokenStr,
		&Claims{},
		func(t *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		},
	)
	if err != nil {
		return 0, "", err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return 0, "", errors.New("invalid token")
	}

	return claims.UserID, claims.Role, nil
}

func Protected(isAdmin bool) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(401, map[string]string{"error": "missing token"})
			}

			token := strings.TrimPrefix(authHeader, "Bearer ")

			userID, role, err := Authorize(token, c)
			if err != nil {
				return c.JSON(401, map[string]string{"error": "unauthorized"})
			}
			
			if isAdmin && role == "USER" {
				return c.JSON(401, map[string]string{"error": "not enough rights"})

			}

			c.Set("user_id", userID)
			c.Set("role", role)

			return next(c)
		}
	}
}
