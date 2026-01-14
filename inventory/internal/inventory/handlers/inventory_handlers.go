package handlers

import (
	"casino/inventory/internal/inventory/service"
	"net/http"

	"github.com/labstack/echo/v4"
)

type InventoryHandler struct {
	service *service.InventoryService
}

func NewInventoryCarHandler(service *service.InventoryService) *InventoryHandler {
	return &InventoryHandler{service: service}
}

func (a *InventoryHandler) SellUserInventory(c echo.Context) error {

	userId := c.Get("user_id").(int)

	total, err := a.service.SellUserInventory(userId)

	if err != nil {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "success",
			"data":    err,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
		"data":    total,
	})
}

func (a *InventoryHandler) SellUserCar(c echo.Context) error {

	req := struct {
		CarId int `json:"id"`
	}{}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid input",
		})
	}

	userId := c.Get("user_id").(int)

	total, err := a.service.SellUserCar(userId, req.CarId)

	if err != nil {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "success",
			"data":    err,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
		"data":    total,
	})
}
