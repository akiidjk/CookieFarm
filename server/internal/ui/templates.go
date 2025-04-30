package ui

import (
	"path/filepath"

	"github.com/ByteTheCookies/cookieserver/internal/logger"
	"github.com/gofiber/template/html/v2"
)

func InitTemplateEngine(debug bool) *html.Engine {
	path, err := filepath.Abs("internal/ui/views")
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Unable to resolve absolute template path")
	}
	logger.Log.Debug().Str("template_path", path).Msg("Using disk templates")
	engine := html.New(path, ".html")

	engineFuncMap := map[string]interface{}{
		"mul":    func(a, b int) int { return a * b },
		"add":    func(a, b int) int { return a + b },
		"sub":    func(a, b int) int { return a - b },
		"subu64": func(a, b uint64) uint64 { return a - b },
	}
	engine.AddFuncMap(engineFuncMap)

	return engine
}
