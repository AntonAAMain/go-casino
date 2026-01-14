package handlers

// import (
// 	"strings"

// 	"github.com/labstack/echo/v4"
// )

// func (h *AuthHandler) Protected(next echo.HandlerFunc) echo.HandlerFunc {
// 	return func(c echo.Context) error {
// 		authHeader := c.Request().Header.Get("Authorization")
// 		if authHeader == "" {
// 			return c.JSON(401, map[string]string{"error": "missing token"})
// 		}

// 		token := strings.TrimPrefix(authHeader, "Bearer ")

// 		user_id, role, err := h.service.Authorize(token, c)
// 		if err != nil {
// 			return c.JSON(401, map[string]string{"error": "unauthorized"})
// 		}

// 		c.Set("user_id", user_id)
// 		c.Set("role", role)

// 		return next(c)
// 	}
// }
