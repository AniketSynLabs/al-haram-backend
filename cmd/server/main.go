package main

import (
	"log"
	"net/http"
	"os"

	"al-haram/config"
	"al-haram/internal/db"
	"al-haram/internal/handler"
	mw "al-haram/internal/middleware"
	"al-haram/web"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	_ = godotenv.Load()

	cfg := config.Load()

	if err := db.Connect(cfg.DatabaseURL); err != nil {
		log.Fatalf("cannot connect to database: %v", err)
	}
	// if err := db.Migrate(); err != nil {
	// 	log.Fatalf("migration failed: %v", err)
	// }

	e := echo.New()
	e.HideBanner = true

	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	pub := e.Group("/api")
	admin := e.Group("/api/admin", mw.BasicAuth(cfg.AdminUsername, cfg.AdminPassword))

	handler.RegisterRoutes(pub, admin, cfg)

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	// Serve uploaded files from disk
	uploadDir := os.Getenv("UPLOAD_DIR")
	if uploadDir == "" {
		uploadDir = "/app/uploads"
	}
	if err := os.MkdirAll(uploadDir, 0o755); err != nil {
		log.Printf("warning: cannot create upload dir: %v", err)
	}
	e.Static("/uploads", uploadDir)

	handler.RegisterStaticFS(web.Dist)
	e.GET("/*", handler.ServeStatic)

	log.Printf("🚀 Server running on :%s", cfg.Port)
	log.Fatal(e.Start(":" + cfg.Port))
}
