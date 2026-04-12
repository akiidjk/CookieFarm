package cmd

import (
	"context"
	"os"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/fang"
)

func ParseArgs(version string, theme fang.ColorScheme, useBanner *bool) {
	buildCmd(useBanner)

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
