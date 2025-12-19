package frontend

import (
	"embed"
	"io/fs"
	"net/http"
	"strings"

	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
)

//go:embed dist/*
var distFS embed.FS

// Handler returns a fasthttp handler that serves the embedded frontend files
func Handler() fasthttp.RequestHandler {
	// Get the dist subdirectory
	distSubFS, err := fs.Sub(distFS, "dist")
	if err != nil {
		panic("failed to get dist subdirectory: " + err.Error())
	}

	// Create file server
	fileServer := http.FileServer(http.FS(distSubFS))

	// Wrap with SPA fallback
	spaHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		// Try to serve the file
		if path != "/" && !strings.HasPrefix(path, "/api") {
			// Check if file exists
			if f, err := distSubFS.Open(strings.TrimPrefix(path, "/")); err == nil {
				f.Close()
				fileServer.ServeHTTP(w, r)
				return
			}
		}

		// For root or non-existent files (SPA routes), serve index.html
		if path == "/" || (!strings.HasPrefix(path, "/api") && !strings.Contains(path, ".")) {
			r.URL.Path = "/"
			fileServer.ServeHTTP(w, r)
			return
		}

		// Serve the actual file
		fileServer.ServeHTTP(w, r)
	})

	// Convert to fasthttp handler
	return fasthttpadaptor.NewFastHTTPHandler(spaHandler)
}

// IsEmbedded returns true if the frontend dist folder is embedded
func IsEmbedded() bool {
	entries, err := distFS.ReadDir("dist")
	if err != nil {
		return false
	}
	return len(entries) > 0
}
