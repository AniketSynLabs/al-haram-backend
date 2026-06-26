package handler

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"mime"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"al-haram/config"

	"github.com/labstack/echo/v4"
)

const maxUploadSize = 20 << 20 // 20 MB

func registerUploadRoutes(admin *echo.Group, cfg *config.Config) {
	admin.POST("/upload", func(c echo.Context) error {
		return handleUpload(c, cfg)
	})
}

func handleUpload(c echo.Context, cfg *config.Config) error {
	if err := c.Request().ParseMultipartForm(maxUploadSize); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "file too large or bad request"})
	}

	file, header, err := c.Request().FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "missing file field"})
	}
	defer file.Close()

	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	ext := strings.ToLower(filepath.Ext(header.Filename))
	if ct := contentType; ct != "" {
		if exts, _ := mime.ExtensionsByType(ct); len(exts) > 0 {
			ext = exts[len(exts)-1]
		}
	}
	if ext == "" {
		ext = ".bin"
	}

	filename := fmt.Sprintf("%d-%s%s", time.Now().UnixMilli(), randStr(8), ext)

	data, err := io.ReadAll(file)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to read file"})
	}

	publicURL, err := uploadToSupabase(cfg, filename, contentType, data)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"url": publicURL})
}

func uploadToSupabase(cfg *config.Config, filename, contentType string, data []byte) (string, error) {
	uploadURL := fmt.Sprintf("%s/storage/v1/object/%s/%s",
		strings.TrimRight(cfg.SupabaseURL, "/"),
		cfg.SupabaseBucket,
		filename,
	)

	req, err := http.NewRequest(http.MethodPost, uploadURL, bytes.NewReader(data))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Authorization", "Bearer "+cfg.SupabaseKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("supabase upload failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("supabase upload error %d: %s", resp.StatusCode, string(body))
	}

	publicURL := fmt.Sprintf("%s/storage/v1/object/public/%s/%s",
		strings.TrimRight(cfg.SupabaseURL, "/"),
		cfg.SupabaseBucket,
		filename,
	)
	return publicURL, nil
}

func randStr(n int) string {
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
}
