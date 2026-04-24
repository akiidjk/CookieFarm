package cmd

import (
	"context"
	"logger"
	"os"

	"client/tui"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/fang"
)

var (
	useBanner bool = true
	useTUI    bool = false
)

func startTui() {
	err := tui.StartTUI(logger.GetBanner("client"))
	if err != nil {
		logger.Log.Error().Err(err).Msg("Error starting TUI")
		if !logger.IsCompletionCommand() {
			logger.PrintBanner(useBanner, "client")
		}
	}
}

func ParseArgs(version string, theme fang.ColorScheme) {
	buildCmd(&useBanner, &useTUI)

	if useTUI {
		startTui()
	} else if !logger.IsCompletionCommand() {
		logger.PrintBanner(useBanner, "client")
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
