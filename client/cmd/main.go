package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/ByteTheCookies/cookiefarm-client/internal/api"
	"github.com/ByteTheCookies/cookiefarm-client/internal/config"
	"github.com/ByteTheCookies/cookiefarm-client/internal/logger"
	"github.com/ByteTheCookies/cookiefarm-client/internal/models"
	"github.com/ByteTheCookies/cookiefarm-client/internal/utils"
	"github.com/google/uuid"
)

func GenerateFakeFlag(flagCode string) models.Flag {
	randomService := config.FAKE_SERVICES[utils.RandInt(1, len(config.FAKE_SERVICES))]
	id, _ := uuid.NewV7()
	return models.Flag{
		ID:           id.String(),
		FlagCode:     flagCode,
		ServiceName:  randomService.Name,
		ServicePort:  randomService.Port,
		SubmitTime:   uint64(time.Now().UnixNano()),
		ResponseTime: 0,
		Status:       "UNSUBMITTED",
		TeamID:       uint16(utils.RandInt(1, 40)),
	}
}

var exploitPath *string

func init() {
	exploitPath = flag.String("exploit", "", "Percorso all'exploit da eseguire")
	debug := flag.Bool("debug", false, "Abilita il livello di log debug")

	flag.Parse()

	if *exploitPath == "" {
		fmt.Println("Errore: devi specificare il percorso dell'exploit con --exploit <path>")
		os.Exit(1)
	}

	if *debug {
		logger.SetLevel(logger.DebugLevel)
	} else {
		logger.SetLevel(logger.InfoLevel)
	}

}

func main() {
	var flags []models.Flag
	cmd := exec.Command(*exploitPath)

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
			flagCode := scanner.Text()
			fmt.Println("[stdout]", flagCode)
			flag := GenerateFakeFlag(flagCode)
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
			time.Sleep(time.Duration(config.CYCLE_TIME) * time.Second)
			api.SendFlag(flags...)
			flags = []models.Flag{}
		}
	}()

	if err := cmd.Wait(); err != nil {
		fmt.Println("Errore comando:", err)
	}
}
