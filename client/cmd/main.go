package main

//   ____            _    _      _____
//  / ___|___   ___ | | _(_) ___|  ___|_ _ _ __ _ __ ___
// | |   / _ \ / _ \| |/ / |/ _ \ |_ / _` | '__| '_ ` _ \
// | |__| (_) | (_) |   <| |  __/  _| (_| | |  | | | | | |
//  \____\___/ \___/|_|\_\_|\___|_|__\__,_|_|  |_| |_| |_|
//  / ___| |   |_ _| ____| \ | |_   _|
// | |   | |    | ||  _| |  \| | | |
// | |___| |___ | || |___| |\  | | |
//  \____|_____|___|_____|_| \_| |_|

import (
	"bufio"
	"fmt"
	"os/exec"
	"time"

	"github.com/ByteTheCookies/cookiefarm-client/internal/api"
	"github.com/ByteTheCookies/cookiefarm-client/internal/config"
	"github.com/ByteTheCookies/cookiefarm-client/internal/logger"
	"github.com/ByteTheCookies/cookiefarm-client/internal/models"
	"github.com/ByteTheCookies/cookiefarm-client/internal/utils"
)

func init() {
	logger.SetLevel(logger.DebugLevel)
}

func main() {
	var flags []models.Flag
	cmd := exec.Command("../tests/exploit.py")

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("Errore pipe stdout:", err)
		return
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		fmt.Println("Errore pipe stderr:", err)
		return
	}

	if err := cmd.Start(); err != nil {
		fmt.Println("Errore start:", err)
		return
	}

	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			flagCode := scanner.Text()
			fmt.Println("[stdout]", flagCode)
			flag := utils.GenerateFakeFlag(flagCode)
			flags = append(flags, flag)
			logger.Debug("Generated flag: %v", flag)
		}
	}()

	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			fmt.Println("[stderr]", scanner.Text())
		}
	}()

	go func() {
		for {
			time.Sleep(config.CYCLE_TIME * time.Second)
			api.SendFlag(flags...)
			flags = []models.Flag{}
		}
	}()

	if err := cmd.Wait(); err != nil {
		fmt.Println("Errore comando:", err)
	}
}
