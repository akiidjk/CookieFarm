// Package utils provides utility functions for the CookieFarm client.
package utils

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"

	"github.com/ByteTheCookies/cookieclient/internal/config"
	"github.com/ByteTheCookies/cookieclient/internal/models"
	"gopkg.in/yaml.v3"
)

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

var pathRegex = regexp.MustCompile(`(~)([^/]*)(/?.*)`)

// Detach detaches the current process from the terminal re executing itself.
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
		fmt.Println(Red+"[ERROR]"+Reset+" | Error during detach:", err)
		os.Exit(1)
	}

	fmt.Println(Yellow+"[WARN]"+Reset+" | Process detached with PID:", cmd.Process.Pid)
	os.Exit(0)
}

// MapPortToService maps a port to a service name.
func MapPortToService(port uint16) string {
	for _, service := range config.Current.ConfigClient.Services {
		if service.Port == port {
			return service.Name
		}
	}
	return ""
}

// GetExecutableDir returns the directory of the executable.
func GetExecutableDir() string {
	exePath, err := os.Executable()
	if err != nil {
		panic("impossible to determine the binary path: " + err.Error())
	}
	return filepath.Dir(exePath)
}

// ValidateArgs validates the arguments passed to the program.
func ValidateArgs(args models.ArgsAttack) error {
	if args.TickTime < 1 {
		return errors.New("tick time must be at least 1")
	}

	exploitPath, err := filepath.Abs(args.ExploitPath)
	if err != nil {
		return fmt.Errorf("error resolving exploit path: %v", err)
	}

	if info, err := os.Stat(exploitPath); err == nil && info.Mode()&0o111 == 0 {
		return errors.New("exploit file is not executable")
	}

	if _, err := os.Stat(exploitPath); os.IsNotExist(err) {
		return errors.New("exploit not found in the exploits directory")
	}

	return nil
}

func IsPath(pathExploit string) bool {
	if strings.HasPrefix(pathExploit, "/") || strings.HasPrefix(pathExploit, ".") || strings.HasPrefix(pathExploit, "~") {
		return true
	}
	return false
}

func IsValid(fp string) bool {
	if _, err := os.Stat(fp); err == nil {
		return true
	}

	var d []byte
	if err := os.WriteFile(fp, d, 0o644); err == nil {
		os.Remove(fp)
		return true
	}

	return false
}

// Code by @prep on Github https://github.com/prep/tilde
func ExpandTilde(p string) (string, error) {
	if len(p) < 1 || p[0] != '~' {
		return p, nil
	}

	var tildePath string

	results := pathRegex.FindStringSubmatch(p)[2:]
	switch results[0] {
	case "":
		u, err := user.Current()
		if err != nil {
			return "", err
		}

		tildePath = u.HomeDir
	case "+":
		pwd, err := os.Getwd()
		if err != nil {
			return "", err
		}

		tildePath = pwd
	default:
		u, err := user.Lookup(results[0])
		if err != nil {
			return "", err
		}

		tildePath = u.HomeDir
	}

	return path.Join(tildePath, results[1]), nil
}

func NormalizeNamePathExploit(name string) (string, error) {
	if !strings.HasSuffix(name, ".py") {
		name += ".py"
	}

	var err error
	if strings.HasPrefix(name, "~") {
		name, err = ExpandTilde(name)
		if err != nil {
			return "", err
		}
	}

	return name, nil
}

func GetSession() (string, error) {
	sessionPath := filepath.Join(GetExecutableDir(), "session")
	data, err := os.ReadFile(sessionPath)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func LoadLocalConfig() error {
	expandendPath, err := ExpandTilde(config.DefaultConfigPath)
	configPath := filepath.Join(expandendPath, "config.yml")
	if err != nil {
		return err
	}
	configFileContent, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("config file does not exist at %s", configPath)
		}
		return fmt.Errorf("error reading config file: %w", err)
	}

	fmt.Println(string(configFileContent))

	err = yaml.Unmarshal(configFileContent, &config.ArgsConfig)
	if err != nil {
		return err
	}

	return nil
}
