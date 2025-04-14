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
	"github.com/ByteTheCookies/cookiefarm-client/internal/logger"
	"github.com/ByteTheCookies/cookiefarm-client/internal/models"
	"github.com/ByteTheCookies/cookiefarm-client/internal/utils"
	"github.com/google/uuid"
)

const CYCLE_TIME = 15

func init() {
	logger.SetLevel(logger.DebugLevel)
}

func run_exploit() {
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
			fmt.Println("[stdout]", scanner.Text())
			flag := models.Flag{
				ID:           int64(uuid.New().ID()),
				FlagCode:     scanner.Text(),
				ServiceName:  "Pippo",
				ResponseTime: time.Now().UnixNano(),
				Status:       "BOH",
				TeamID:       0,
			}
			flags = append(flags, flag)
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
			time.Sleep(CYCLE_TIME * time.Second)
			api.SendFlag(flags...)
			flags = []models.Flag{}
		}
	}()

	if err := cmd.Wait(); err != nil {
		fmt.Println("Errore comando:", err)
	}
}

func main() {
	fmt.Printf(utils.Yellow + `
	   ______            __   _      ______
	  / ____/___  ____  / /__(_)__  / ____/___ __________ ___
	 / /   / __ \/ __ \/ //_/ / _ \/ /_  / __ ` + `/ ___/ __` + `__ \
	/ /___/ /_/ / /_/ / ,< / /  __/ __/ / /_/ / /  / / / / / /
	\____/\____/\____/_/|_/_/\___/_/    \__,_/_/  /_/ /_/ /_/
	  / ____/ (_)__  ____  / /_
	 / /   / / / _ \/ __ \/ __/
	/ /___/ / /  __/ / / / /_
	\____/_/_/\___/_/ /_/\__/
 ` + utils.Reset)
	run_exploit()
}
