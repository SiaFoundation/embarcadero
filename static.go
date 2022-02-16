package main

import (
	"embed"
	"io/fs"
	"net/http"
	"os"
	"path"
)

var (
	//go:embed app/dist/apps/embarcadero
	Assets   embed.FS
	buildDir = "app/dist/apps/embarcadero"
)

// Static file server for UI

type fsFunc func(name string) (fs.File, error)

func (f fsFunc) Open(name string) (fs.File, error) {
	return f(name)
}

func buildUIHandler() http.Handler {
	defaultPath := path.Join(buildDir, "index.html")

	handler := fsFunc(func(name string) (fs.File, error) {
		assetPath := path.Join(buildDir, name)
		f, err := Assets.Open(assetPath)

		if os.IsNotExist(err) {
			// Fallback to index.html
			return Assets.Open(defaultPath)
		}

		return f, err
	})

	return http.FileServer(http.FS(handler))
}
