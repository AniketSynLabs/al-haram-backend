package handler

import (
	"net/http"

	"al-haram/internal/model"
	"al-haram/internal/service"

	"github.com/labstack/echo/v4"
)

func registerServiceRoutes(pub, admin *echo.Group) {
	pub.GET("/services", getServices)
	admin.PUT("/services/:id", updateService)
}

func getServices(c echo.Context) error {
	svcs, err := service.ListServices()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err500(err))
	}
	return c.JSON(http.StatusOK, svcs)
}

func updateService(c echo.Context) error {
	var s model.Service
	if err := c.Bind(&s); err != nil {
		return c.JSON(http.StatusBadRequest, err400(err))
	}
	if err := service.UpdateService(c.Param("id"), s); err != nil {
		return c.JSON(http.StatusInternalServerError, err500(err))
	}
	return c.JSON(http.StatusOK, ok())
}
