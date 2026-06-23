package handler

import (
	"net/http"
	"time"

	"al-haram/internal/model"
	"al-haram/internal/service"

	"github.com/labstack/echo/v4"
)

func registerPackageRoutes(pub, admin *echo.Group) {
	pub.GET("/packages", getPackages)
	admin.POST("/packages", createPackage)
	admin.PUT("/packages/:id", updatePackage)
	admin.DELETE("/packages/:id", deletePackage)
}

func getPackages(c echo.Context) error {
	pkgs, err := service.ListPackages()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err500(err))
	}
	return c.JSON(http.StatusOK, pkgs)
}

func createPackage(c echo.Context) error {
	var p model.Package
	if err := c.Bind(&p); err != nil {
		return c.JSON(http.StatusBadRequest, err400(err))
	}
	if p.ID == "" {
		p.ID = "pkg_" + time.Now().Format("20060102150405")
	}
	created, err := service.CreatePackage(p)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err500(err))
	}
	service.InvalidateBootstrapCache()
	return c.JSON(http.StatusCreated, created)
}

func updatePackage(c echo.Context) error {
	var p model.Package
	if err := c.Bind(&p); err != nil {
		return c.JSON(http.StatusBadRequest, err400(err))
	}
	if err := service.UpdatePackage(c.Param("id"), p); err != nil {
		return c.JSON(http.StatusInternalServerError, err500(err))
	}
	service.InvalidateBootstrapCache()
	return c.JSON(http.StatusOK, ok())
}

func deletePackage(c echo.Context) error {
	if err := service.DeletePackage(c.Param("id")); err != nil {
		return c.JSON(http.StatusInternalServerError, err500(err))
	}
	service.InvalidateBootstrapCache()
	return c.JSON(http.StatusOK, ok())
}
