package handler

import (
	"net/http"

	"al-haram/internal/service"

	"github.com/labstack/echo/v4"
)

func registerBootstrapRoutes(pub *echo.Group) {
	pub.GET("/bootstrap", getBootstrap)
}

func getBootstrap(c echo.Context) error {
	data, err := service.GetBootstrap()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err500(err))
	}
	return c.JSON(http.StatusOK, data)
}
