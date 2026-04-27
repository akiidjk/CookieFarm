package cmd

import (
	"context"
	"logger"
	"os"

	"client/tui"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/fang"
)

var useTUI bool = false

func startTui() {
	err := tui.StartTUI(logger.GetBanner("client"))
	if err != nil {
		logger.Log.Error().Err(err).Msg("Error starting TUI")
	}
}

func ParseArgs(version string, theme fang.ColorScheme) {
	buildCmd(&useTUI)

	if useTUI {
		startTui()
	}

	err := fang.Execute(
		context.TODO(),
		rootCmd,
		fang.WithVersion(version),
		fang.WithColorSchemeFunc(func(ld lipgloss.LightDarkFunc) fang.ColorScheme {
			return theme
		}),
	)
	if err != nil {
		os.Exit(1)
	}
}
