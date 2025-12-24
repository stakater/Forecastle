package web

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed build/*
var frontendFS embed.FS

// GetFrontendFS returns the embedded frontend filesystem
func GetFrontendFS() http.FileSystem {
	// Strip the "build" prefix so files are served from root
	sub, err := fs.Sub(frontendFS, "build")
	if err != nil {
		panic(err)
	}
	return http.FS(sub)
}
