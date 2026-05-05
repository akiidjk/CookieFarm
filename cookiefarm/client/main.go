package main

import (
	"logger"

	"client/cmd"
)

var Version = "dev"

func main() {
	cmd.ParseArgs(Version, logger.CookieCLIColorSchema)
}
