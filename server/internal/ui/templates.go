package ui

import (
	"embed"
	"io/fs"
	"net/http"
	"path/filepath"

	"github.com/ByteTheCookies/backend/internal/logger"
)

//go:embed views/*
var EmbedViews embed.FS

func ViewsFS() http.FileSystem {
	sub, err := fs.Sub(EmbedViews, "views")
	if err != nil {
		panic("failed to create sub FS: " + err.Error())
	}
	return http.FS(sub)
}

func GetPathView() string {
	path, err := filepath.Abs("internal/ui/views")
	if err != nil {
		logger.Log.Fatal().Msgf("cannot resolve absolute path: %v", err)
	}
	return path
}
