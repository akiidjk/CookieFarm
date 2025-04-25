package ui

import (
	"embed"
	"io/fs"
	"net/http"
	"path/filepath"

	"github.com/ByteTheCookies/backend/internal/logger"
	"github.com/gofiber/template/html/v2"
)

//go:embed views/*
var EmbedViews embed.FS

func InitTemplateEngine(debug bool) *html.Engine {
	var engine *html.Engine

	if debug {
		path, err := filepath.Abs("internal/ui/views")
		if err != nil {
			logger.Log.Fatal().Msgf("cannot resolve absolute path: %v", err)
		}
		logger.Log.Debug().Msgf("using template path: %s", path)
		engine = html.New(path, ".html")
	} else {
		sub, err := fs.Sub(EmbedViews, "views")
		if err != nil {
			panic("failed to create sub FS: " + err.Error())
		}
		engine = html.NewFileSystem(http.FS(sub), ".html")
	}

	engine.AddFunc("mul", func(a, b int) int {
		return a * b
	})
	engine.AddFunc("add1", func(a int) int {
		return a + 1
	})

	engine.AddFunc("rangeN", func(n int) []int {
		out := make([]int, n)
		for i := range n {
			out[i] = i
		}
		return out
	})

	return engine
}
