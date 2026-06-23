package handler

import (
	"net/http"
	"time"

	"al-haram/internal/model"
	"al-haram/internal/service"

	"github.com/labstack/echo/v4"
)

func registerServiceRoutes(pub, admin *echo.Group) {
	pub.GET("/services", getServices)
	admin.POST("/services", createService)
	admin.PUT("/services/:id", updateService)
	admin.DELETE("/services/:id", deleteService)
}

func getServices(c echo.Context) error {
	svcs, err := service.ListServices()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err500(err))
	}
	return c.JSON(http.StatusOK, svcs)
}

func createService(c echo.Context) error {
	var s model.Service
	if err := c.Bind(&s); err != nil {
		return c.JSON(http.StatusBadRequest, err400(err))
	}
	if s.ID == "" {
		s.ID = "svc_" + time.Now().Format("20060102150405")
	}
	created, err := service.CreateService(s)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err500(err))
	}
	service.InvalidateBootstrapCache()
	return c.JSON(http.StatusCreated, created)
}

func updateService(c echo.Context) error {
	var s model.Service
	if err := c.Bind(&s); err != nil {
		return c.JSON(http.StatusBadRequest, err400(err))
	}
	if err := service.UpdateService(c.Param("id"), s); err != nil {
		return c.JSON(http.StatusInternalServerError, err500(err))
	}
	service.InvalidateBootstrapCache()
	return c.JSON(http.StatusOK, ok())
}

func deleteService(c echo.Context) error {
	if err := service.DeleteService(c.Param("id")); err != nil {
		return c.JSON(http.StatusInternalServerError, err500(err))
	}
	service.InvalidateBootstrapCache()
	return c.JSON(http.StatusOK, ok())
}
