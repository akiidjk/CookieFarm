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
	config.UseTUI = true
	config.UseBanner = true
	for _, arg := range os.Args {
		switch arg {
		case "--no-tui", "-N", "-h", "--help":
			config.UseTUI = false
			logger.SetTUI(false)
		case "--no-banner", "-B":
			config.UseBanner = false
		case "-v", "--version":
			config.UseTUI = false
			logger.SetTUI(false)
			config.UseBanner = false
		}
	}

	if config.UseTUI && os.Getenv("COOKIECLIENT_NO_TUI") == "" {
		if err := tui.StartTUI(banner); err != nil {
			fmt.Printf("Error starting TUI: %v\nFalling back to CLI mode\n", err)
			if config.UseBanner {
				if !isCompletionCommand() {
					fmt.Println(banner)
				}
			}
			cmd.Execute()
		}
	} else {
		if config.UseBanner {
			if !isCompletionCommand() {
				fmt.Println(banner)
			}
		}
		cmd.Execute()
	}
}
