package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/ByteTheCookies/cookiefarm-client/internal/api"
	"github.com/ByteTheCookies/cookiefarm-client/internal/config"
	"github.com/ByteTheCookies/cookiefarm-client/internal/logger"
	"github.com/ByteTheCookies/cookiefarm-client/internal/models"
	"github.com/ByteTheCookies/cookiefarm-client/internal/utils"
	"github.com/google/uuid"
	"github.com/spf13/pflag"
)

func Flag(stdoutFlag models.StdoutFormat) models.Flag {
	randomService := config.Current.ConfigClient.Services[utils.RandInt(0, len(config.Current.ConfigClient.Services))]
	id, _ := uuid.NewV7()
	return models.Flag{
		ID:           id.String(),
		FlagCode:     stdoutFlag.FlagCode,
		ServiceName:  randomService.Name,
		ServicePort:  stdoutFlag.ServicePort,
		SubmitTime:   uint64(time.Now().Unix()),
		ResponseTime: 0,
		Status:       "UNSUBMITTED",
		TeamID:       stdoutFlag.TeamId,
	}
}

var (
	exploitPath     = pflag.StringP("exploit", "e", "", "Path to the exploit to execute")
	debug           = pflag.Bool("debug", false, "Enable debug log level")
	password        = pflag.StringP("password", "p", "", "Password for authentication")
	base_url_server = pflag.StringP("base_url_server", "b", "", "Base URL of the target server (e.g. http://localhost:8080)")
	detach          = pflag.BoolP("detach", "d", false, "Run the exploit in the background") // alias -d
	threadsNumber   = pflag.IntP("threads", "t", 1, "Number of threads to use")
	tickTime        = pflag.IntP("tick", "i", 120, "Interval in seconds between run exploits ")
)

func init() {
	pflag.Parse()

	if *detach {
		utils.Detach()
	}

	if *exploitPath == "" {
		fmt.Println("Errore: devi specificare il percorso dell'exploit con --exploit <path>")
		os.Exit(1)
	}

	if *base_url_server == "" {
		fmt.Println("Errore: devi specificare il base_url_server con --base_url_server <url>")
		os.Exit(1)
	}

	config.Current.ConfigClient.BaseUrlServer = *base_url_server

	if *password == "" {
		fmt.Println("Errore: devi specificare la password con --password <password>")
		os.Exit(1)
	}

	if *debug {
		logger.SetLevel(logger.DebugLevel)
	} else {
		logger.SetLevel(logger.InfoLevel)
	}

	var err error
	config.Token, err = api.Login(*password)
	if err != nil {
		fmt.Println("Errore login:", err)
		os.Exit(1)
	}

}

func main() {
	var flags []models.Flag
	cmd := exec.Command(*exploitPath, config.Current.ConfigClient.BaseUrlServer, *password, strconv.Itoa(*tickTime), strconv.Itoa(*threadsNumber))

	config.Current = api.GetConfig()

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
			flagJson := scanner.Text()
			flagStdout := models.StdoutFormat{}
			json.Unmarshal([]byte(flagJson), &flagStdout)
			fmt.Println("[stdout]", flagJson)
			flag := Flag(flagStdout)
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
			time.Sleep(time.Duration(config.Current.ConfigClient.SubmitFlagServerTime) * time.Second)
			api.SendFlag(flags...)
			flags = []models.Flag{}
		}
	}()

	if err := cmd.Wait(); err != nil {
		fmt.Println("Errore comando:", err)
	}
}
