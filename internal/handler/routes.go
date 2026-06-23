package handler

import "github.com/labstack/echo/v4"

// RegisterRoutes wires all domain routes onto the public and admin groups.
func RegisterRoutes(pub, admin *echo.Group) {
	registerBootstrapRoutes(pub)
	registerPackageRoutes(pub, admin)
	registerServiceRoutes(pub, admin)
	registerScheduleRoutes(pub, admin)
	registerEnquiryRoutes(pub, admin)
	registerSettingsRoutes(pub, admin)
	registerGalleryRoutes(pub, admin)
}
