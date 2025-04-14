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

	"github.com/ByteTheCookies/cookiefarm-client/internal/logger"
	"github.com/ByteTheCookies/cookiefarm-client/internal/utils"
)

func init() {
	logger.SetLevel(logger.DebugLevel)
}

func run_exploit() {
	var flags []string
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
			flags = append(flags, scanner.Text())
		}
	}()

	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			fmt.Println("[stderr]", scanner.Text())
		}
	}()

	go func() {
		time.Sleep(5 * time.Second)
		logger.Debug("Len flags before %d", len(flags))
		SendFlag(flags...)
		flags = []string{}
		logger.Debug("Len flags after %d", len(flags))
	}()

	if err := cmd.Wait(); err != nil {
		fmt.Println("Errore comando:", err)
	}
}

func SendFlag(flags ...string) {
	fmt.Println("Invio flag:", flags)
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
