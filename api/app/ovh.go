package app

import (
	"net/http"

	"cloud-inventory/services"

	"github.com/labstack/echo/v4"
)

func GetAllServersHandler(c echo.Context) error {
	serverList, err := services.OvhListAllServers()
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, serverList)
}
