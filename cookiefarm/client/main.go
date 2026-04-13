package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"logger"
	"os"
	"sharedconfig"
	"strings"

	"client/api"
	"client/cmd"
	"client/config"
)

var VERSION = sharedconfig.GetVersion()

func configCheck(cm *config.ConfigManager) error {
	remoteCfg, err := api.GetConfig()
	if err != nil {
		return fmt.Errorf("error receiving shared configs: %w", err)
	}

	remoteConfigStr, err := json.Marshal(remoteCfg)
	if err != nil {
		return err
	}

	currentShc, err := json.Marshal(cm.Get().Shared)
	if err != nil {
		return err
	}

	if strings.TrimSpace(string(currentShc)) != strings.TrimSpace(string(remoteConfigStr)) {
		return errors.New("Shared configuration has changed. Updating local config doing `ckc config login`")
	}

	return nil
}

func main() {
	cm := config.GetInstance()
	cm.Read()
	checkErr := configCheck(cm)
	if checkErr != nil {
		if strings.Contains(checkErr.Error(), "Shared configuration has changed") {
			fmt.Fprintf(os.Stderr, "\n\033[1;33m[!] %v\033[0m\n\n", checkErr)
		} else {
			fmt.Fprintf(os.Stderr, "\n\033[1;31m[+] Error checking configuration: %v\033[0m\n\n", checkErr)
		}
	}

	cmd.ParseArgs(VERSION, logger.CookieCLIColorSchema)
}
