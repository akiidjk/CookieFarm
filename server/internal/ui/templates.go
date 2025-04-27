package ui

import (
	"embed"
	"io/fs"
	"net/http"
	"path/filepath"

	"github.com/ByteTheCookies/cookieserver/internal/logger"
	"github.com/gofiber/template/html/v2"
)

//go:embed views/*
var embedViews embed.FS

func InitTemplateEngine(debug bool) *html.Engine {
	var engine *html.Engine
	if debug {
		path, err := filepath.Abs("internal/ui/views")
		if err != nil {
			logger.Log.Fatal().Err(err).Msg("Unable to resolve absolute template path")
		}
		logger.Log.Debug().Str("template_path", path).Msg("Using disk templates")
		engine = html.New(path, ".html")
	} else {
		subFS, err := fs.Sub(embedViews, "views")
		if err != nil {
			logger.Log.Fatal().Err(err).Msg("Failed to load embedded templates")
		}
		engine = html.NewFileSystem(http.FS(subFS), ".html")
	}

	engineFuncMap := map[string]interface{}{
		"mul": func(a, b int) int { return a * b },
		"add": func(a, b int) int { return a + b },
		"sub": func(a, b int) int { return a - b },
	}
	engine.AddFuncMap(engineFuncMap)

	return engine
}
