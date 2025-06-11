package main

import (
	_ "embed"
	"fmt"
	"os"
	"strings"

	"github.com/ByteTheCookies/cookieclient/cmd"
	"github.com/ByteTheCookies/cookieclient/internal/config"
	"github.com/ByteTheCookies/cookieclient/internal/logger"
	"github.com/ByteTheCookies/cookieclient/internal/tui"
)

func isCompletionCommand() bool {
	for _, arg := range os.Args {
		if strings.Contains(arg, "__complete") || strings.Contains(arg, "completion") {
			return true
		}
	}
	return false
}

//go:embed banner.txt
var banner string

func main() {
	cm := config.GetConfigManager()
	debug := false
	cm.SetUseTUI(true)
	logger.SetTUI(true)

	for _, arg := range os.Args {
		switch arg {
		case "-D", "--debug":
			debug = true
		case "--no-tui", "-N", "-h", "--help":
			cm.SetUseTUI(false)
			logger.SetTUI(false)
		case "--no-banner", "-B":
			cm.SetUseBanner(false)

		case "-v", "--version":
			cm.SetUseTUI(false)
			logger.SetTUI(false)
			cm.SetUseBanner(false)
		}
	}

	if cm.GetUseTUI() && os.Getenv("COOKIECLIENT_NO_TUI") == "" {
		if err := tui.StartTUI(banner, debug); err != nil {
			fmt.Printf("Error starting TUI: %v\nFalling back to CLI mode\n", err)
			if cm.GetUseBanner() {
				if !isCompletionCommand() {
					fmt.Println(banner)
				}
			}
			cmd.Execute()
		}
	} else {
		if cm.GetUseBanner() {
			if !isCompletionCommand() {
				fmt.Println(banner)
			}
		}
		cmd.Execute()
	}
}
