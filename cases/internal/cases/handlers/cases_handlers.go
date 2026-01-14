package handlers

import (
	"casino/cases/internal/cases/service"
	"net/http"

	"github.com/labstack/echo/v4"
)

type CasesHandler struct {
	service *service.CasesService
}

func NewBoxCarHandler(service *service.CasesService) *CasesHandler {
	return &CasesHandler{service: service}
}

func (a *CasesHandler) OpenBox(c echo.Context) error {
	req := struct {
		BoxId int `json:"box_id"`
	}{}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid input",
		})
	}

	userID := c.Get("user_id").(int)

	boxCar, err := a.service.OpenCase(userID, req.BoxId)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
		"data":    boxCar,
	})
}
