package main

import (
	"github.com/ByteTheCookies/cookiefarm-client/internal/logger"
)

func init() {
	logger.SetLevel(logger.DebugLevel)
}

func main() {
	logger.Info("Start client")
}
