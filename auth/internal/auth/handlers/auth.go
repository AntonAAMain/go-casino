package handlers

import (
	"casino/auth/internal/auth/service"
	"net/http"

	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	service *service.AuthService
}

func NewAuthHandler(service *service.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

func (a *AuthHandler) Register(c echo.Context) error {

	req := struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}{}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid input",
		})
	}

	token, err := a.service.Register(req.Name, req.Password)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})
}

func (a *AuthHandler) Login(c echo.Context) error {

	req := struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}{}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid input",
		})
	}

	token, err := a.service.Login(req.Name, req.Password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})
}
