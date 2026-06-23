package handler

import (
	"net/http"

	"al-haram/internal/model"
	"al-haram/internal/service"

	"github.com/labstack/echo/v4"
)

func registerGalleryRoutes(pub, admin *echo.Group) {
	pub.GET("/gallery", listGallery)
	admin.GET("/gallery", listGallery)
	admin.POST("/gallery", createGalleryItem)
	admin.PUT("/gallery/:id", updateGalleryItem)
	admin.DELETE("/gallery/:id", deleteGalleryItem)
}

func listGallery(c echo.Context) error {
	items, err := service.ListGallery()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err500(err))
	}
	return c.JSON(http.StatusOK, items)
}

func createGalleryItem(c echo.Context) error {
	var g model.GalleryItem
	if err := c.Bind(&g); err != nil {
		return c.JSON(http.StatusBadRequest, err400(err))
	}
	created, err := service.CreateGalleryItem(g)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err500(err))
	}
	service.InvalidateBootstrapCache()
	return c.JSON(http.StatusCreated, created)
}

func updateGalleryItem(c echo.Context) error {
	var g model.GalleryItem
	if err := c.Bind(&g); err != nil {
		return c.JSON(http.StatusBadRequest, err400(err))
	}
	if err := service.UpdateGalleryItem(c.Param("id"), g); err != nil {
		return c.JSON(http.StatusInternalServerError, err500(err))
	}
	service.InvalidateBootstrapCache()
	return c.JSON(http.StatusOK, ok())
}

func deleteGalleryItem(c echo.Context) error {
	if err := service.DeleteGalleryItem(c.Param("id")); err != nil {
		return c.JSON(http.StatusInternalServerError, err500(err))
	}
	service.InvalidateBootstrapCache()
	return c.JSON(http.StatusOK, ok())
}
