package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"syscall"

	"math/rand"

	"github.com/ByteTheCookies/cookieclient/internal/config"
	"github.com/ByteTheCookies/cookieclient/internal/models"
)

const regexUrl = `^http://(localhost|127\.0\.0\.1):[0-9]{1,5}$`

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

func GetExecutableDir() string {
	exePath, err := os.Executable()
	if err != nil {
		panic("impossible to determine the binary path: " + err.Error())
	}
	return filepath.Dir(exePath)
}

func ValidateArgs(args models.Args) error {

	if *args.ExploitName == "" {
		return fmt.Errorf("missing required --exploit argument")
	}
	if *args.BaseURLServer == "" {
		return fmt.Errorf("missing required --base_url_server argument")
	}
	if *args.Password == "" {
		return fmt.Errorf("missing required --password argument")
	}

	if !regexp.MustCompile(regexUrl).MatchString(*args.BaseURLServer) {
		return fmt.Errorf("invalid base URL server")
	}

	if *args.TickTime < 1 {
		return fmt.Errorf("tick time must be at least 1")
	}

	exploitPath := filepath.Join(GetExecutableDir(), "..", "exploits", *args.ExploitName)

	if _, err := os.Stat(exploitPath); os.IsNotExist(err) {
		return fmt.Errorf("exploit not found in the exploits directory")
	}

	return nil
}
