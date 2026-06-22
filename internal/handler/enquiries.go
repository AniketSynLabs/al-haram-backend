package handler

import (
	"net/http"

	"al-haram/internal/model"
	"al-haram/internal/service"

	"github.com/labstack/echo/v4"
)

func registerEnquiryRoutes(pub, admin *echo.Group) {
	pub.POST("/enquiries", createEnquiry)
	admin.GET("/enquiries", getEnquiries)
	admin.PUT("/enquiries/:id/status", updateEnquiryStatus)
}

func createEnquiry(c echo.Context) error {
	var e model.Enquiry
	if err := c.Bind(&e); err != nil {
		return c.JSON(http.StatusBadRequest, err400(err))
	}
	created, err := service.CreateEnquiry(e)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err500(err))
	}
	return c.JSON(http.StatusCreated, created)
}

func getEnquiries(c echo.Context) error {
	enquiries, err := service.ListEnquiries(c.QueryParam("status"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err500(err))
	}
	return c.JSON(http.StatusOK, enquiries)
}

func updateEnquiryStatus(c echo.Context) error {
	var body struct {
		Status string `json:"status"`
	}
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, err400(err))
	}
	if err := service.UpdateEnquiryStatus(c.Param("id"), body.Status); err != nil {
		return c.JSON(http.StatusInternalServerError, err500(err))
	}
	return c.JSON(http.StatusOK, ok())
}
