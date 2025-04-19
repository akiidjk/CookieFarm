package utils

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"syscall"

	"math/rand"

	"github.com/ByteTheCookies/cookiefarm-client/internal/config"
)

func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func Detach() {
	cmd := exec.Command(os.Args[0], os.Args[1:]...)

	filteredArgs := []string{}
	for _, arg := range os.Args[1:] {
		if arg != "--detach" && arg != "-d" {
			filteredArgs = append(filteredArgs, arg)
		}
	}
	cmd = exec.Command(os.Args[0], filteredArgs...)

	// Scollega input/output/terminal
	cmd.Stdout = nil
	cmd.Stderr = nil
	cmd.Stdin = nil
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}

	err := cmd.Start()
	if err != nil {
		fmt.Println("Errore nel detach:", err)
		os.Exit(1)
	}

	fmt.Println("Process detached with PID:", cmd.Process.Pid)
	os.Exit(0)
}

const (
	Reset   = "\033[0m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	Gray    = "\033[37m"
	White   = "\033[97m"
)

func CleanGC() (uint64, uint64) {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	before := mem.Alloc / 1_048_576
	runtime.GC() //Cleaning garbage collector
	runtime.ReadMemStats(&mem)
	after := mem.Alloc / 1_048_576
	return before, after
}

func RandInt(min, max int) int {
	return min + rand.Intn(max-min)
}

func MapPortToService(port uint16) string {
	for _, service := range config.Current.ConfigClient.Services {
		if service.Port == port {
			return service.Name
		}
	}
	return ""
}
