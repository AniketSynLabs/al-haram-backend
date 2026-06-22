package handler

import (
	"net/http"

	"al-haram/internal/model"
	"al-haram/internal/service"

	"github.com/labstack/echo/v4"
)

func registerScheduleRoutes(pub, admin *echo.Group) {
	pub.GET("/schedule", getSchedule)
	admin.GET("/schedule", getSchedule)
	admin.POST("/schedule", createScheduleEntry)
	admin.PUT("/schedule/:id", updateScheduleEntry)
	admin.DELETE("/schedule/:id", deleteScheduleEntry)
}

func getSchedule(c echo.Context) error {
	entries, err := service.ListSchedule()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err500(err))
	}
	return c.JSON(http.StatusOK, entries)
}

func createScheduleEntry(c echo.Context) error {
	var s model.ScheduleEntry
	if err := c.Bind(&s); err != nil {
		return c.JSON(http.StatusBadRequest, err400(err))
	}
	created, err := service.CreateScheduleEntry(s)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err500(err))
	}
	return c.JSON(http.StatusCreated, created)
}

func updateScheduleEntry(c echo.Context) error {
	var s model.ScheduleEntry
	if err := c.Bind(&s); err != nil {
		return c.JSON(http.StatusBadRequest, err400(err))
	}
	if err := service.UpdateScheduleEntry(c.Param("id"), s); err != nil {
		return c.JSON(http.StatusInternalServerError, err500(err))
	}
	return c.JSON(http.StatusOK, ok())
}

func deleteScheduleEntry(c echo.Context) error {
	if err := service.DeleteScheduleEntry(c.Param("id")); err != nil {
		return c.JSON(http.StatusInternalServerError, err500(err))
	}
	return c.JSON(http.StatusOK, ok())
}
