package middleware

import (
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func BasicAuth(username, password string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			auth := c.Request().Header.Get("Authorization")
			if !strings.HasPrefix(auth, "Basic ") {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
			}
			decoded, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(auth, "Basic "))
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
			}
			parts := strings.SplitN(string(decoded), ":", 2)
			if len(parts) != 2 || parts[0] != username || parts[1] != password {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid credentials"})
			}
			return next(c)
		}
	}
}
