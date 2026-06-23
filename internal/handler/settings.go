package handler

import (
	"net/http"

	"al-haram/internal/model"
	"al-haram/internal/service"

	"github.com/labstack/echo/v4"
)

func registerSettingsRoutes(pub, admin *echo.Group) {
	pub.GET("/settings", getSettings)
	pub.GET("/bank-details", getBankDetails)
	pub.GET("/policies", getPolicies)

	admin.GET("/settings", getSettings)
	admin.PUT("/settings", updateSettings)
	admin.GET("/bank-details", getBankDetails)
	admin.PUT("/bank-details", updateBankDetails)
	admin.GET("/policies", getPolicies)
	admin.PUT("/policies/:id", updatePolicy)
}

func getSettings(c echo.Context) error {
	settings, err := service.GetSettings()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err500(err))
	}
	return c.JSON(http.StatusOK, settings)
}

func updateSettings(c echo.Context) error {
	var body map[string]string
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, err400(err))
	}
	if err := service.UpdateSettings(body); err != nil {
		return c.JSON(http.StatusInternalServerError, err500(err))
	}
	service.InvalidateBootstrapCache()
	return c.JSON(http.StatusOK, ok())
}

func getBankDetails(c echo.Context) error {
	b, err := service.GetBankDetails()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err500(err))
	}
	return c.JSON(http.StatusOK, b)
}

func updateBankDetails(c echo.Context) error {
	var b model.BankDetails
	if err := c.Bind(&b); err != nil {
		return c.JSON(http.StatusBadRequest, err400(err))
	}
	if err := service.UpdateBankDetails(b); err != nil {
		return c.JSON(http.StatusInternalServerError, err500(err))
	}
	service.InvalidateBootstrapCache()
	return c.JSON(http.StatusOK, ok())
}

func getPolicies(c echo.Context) error {
	policies, err := service.ListPolicies()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err500(err))
	}
	return c.JSON(http.StatusOK, policies)
}

func updatePolicy(c echo.Context) error {
	var body struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, err400(err))
	}
	if err := service.UpdatePolicy(c.Param("id"), body.Title, body.Content); err != nil {
		return c.JSON(http.StatusInternalServerError, err500(err))
	}
	service.InvalidateBootstrapCache()
	return c.JSON(http.StatusOK, ok())
}
