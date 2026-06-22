package handler

import (
	"embed"
	"io"
	"io/fs"
	"net/http"

	"github.com/labstack/echo/v4"
)

var distFS fs.FS

// RegisterStaticFS must be called before ServeStatic is used.
func RegisterStaticFS(embedded embed.FS) {
	sub, err := fs.Sub(embedded, "dist")
	if err != nil {
		panic("web/dist not embedded: " + err.Error())
	}
	distFS = sub
}

// ServeStatic serves files from the embedded dist/ folder.
// Unknown paths fall back to index.html so React Router can handle them.
func ServeStatic(c echo.Context) error {
	urlPath := c.Request().URL.Path

	// Check if the file actually exists in the embedded FS
	rel := urlPath
	if len(rel) > 0 && rel[0] == '/' {
		rel = rel[1:]
	}
	if rel == "" {
		rel = "index.html"
	}

	if f, err := distFS.Open(rel); err == nil {
		f.Close()
		// Serve the real asset as-is
		http.FileServer(http.FS(distFS)).ServeHTTP(c.Response().Writer, c.Request())
		return nil
	}

	// File not found → SPA fallback: serve index.html
	index, err := distFS.Open("index.html")
	if err != nil {
		return echo.ErrNotFound
	}
	defer index.Close()

	stat, err := index.Stat()
	if err != nil {
		return echo.ErrInternalServerError
	}

	http.ServeContent(c.Response().Writer, c.Request(), "index.html", stat.ModTime(), index.(io.ReadSeeker))
	return nil
}
