package main

import (
	_ "embed"
	"fmt"
	"os"

	"github.com/ByteTheCookies/cookieclient/cmd"
	"github.com/ByteTheCookies/cookieclient/internal/config"
	"github.com/ByteTheCookies/cookieclient/internal/tui"
)

//go:embed banner.txt
var banner string

func main() {
	config.NoTUI = true
	config.UseBanner = true
	for _, arg := range os.Args {
		if arg == "--no-tui" || arg == "-N" {
			config.NoTUI = false
		}
		if arg == "--no-banner" || arg == "-B" {
			config.UseBanner = false
		}
	}

	if config.NoTUI && os.Getenv("COOKIECLIENT_NO_TUI") == "" {
		if err := tui.StartTUI(banner); err != nil {
			fmt.Printf("Error starting TUI: %v\nFalling back to CLI mode\n", err)
			if config.UseBanner {
				fmt.Println(banner)
			}
			cmd.Execute()
		}
	} else {
		if config.UseBanner {
			fmt.Println(banner)
		}
		cmd.Execute()
	}
}
