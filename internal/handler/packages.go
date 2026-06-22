package handler

import (
	"net/http"

	"al-haram/internal/model"
	"al-haram/internal/service"

	"github.com/labstack/echo/v4"
)

func registerPackageRoutes(pub, admin *echo.Group) {
	pub.GET("/packages", getPackages)
	admin.PUT("/packages/:id", updatePackage)
}

func getPackages(c echo.Context) error {
	pkgs, err := service.ListPackages()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err500(err))
	}
	return c.JSON(http.StatusOK, pkgs)
}

func updatePackage(c echo.Context) error {
	var p model.Package
	if err := c.Bind(&p); err != nil {
		return c.JSON(http.StatusBadRequest, err400(err))
	}
	if err := service.UpdatePackage(c.Param("id"), p); err != nil {
		return c.JSON(http.StatusInternalServerError, err500(err))
	}
	return c.JSON(http.StatusOK, ok())
}
