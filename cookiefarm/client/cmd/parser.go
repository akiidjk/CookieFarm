package cmd

import (
	"context"
	"os"

	"github.com/charmbracelet/fang"
)

func ParseArgs(version string, theme fang.ColorScheme, useBanner *bool) {
	buildCmd(useBanner)

	err := fang.Execute(
		context.TODO(),
		rootCmd,
		fang.WithVersion(version),
		fang.WithTheme(theme),
	)
	if err != nil {
		os.Exit(1)
	}
}
