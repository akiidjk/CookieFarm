package main

import (
	"logger"
	"sharedconfig"

	"client/cmd"
	"client/config"
	"client/tui"
)

var (
	useBanner bool
	useTUI    bool
)

var VERSION = sharedconfig.GetVersion()

func startTui() {
	err := tui.StartTUI(logger.GetBanner("client"))
	if err != nil {
		logger.Log.Error().Err(err).Msg("Error starting TUI")
		if !logger.IsCompletionCommand() {
			logger.PrintBanner(useBanner, "client")
		}
	}
}

func main() {
	cm := config.GetInstance()
	cm.Read()

	cmd.ParseArgs(VERSION, logger.CookieCLIColorSchema, &useBanner)

	if useTUI {
		startTui()
	} else if !logger.IsCompletionCommand() {
		logger.PrintBanner(useBanner, "client")
	}
}
