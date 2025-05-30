package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ByteTheCookies/cookieclient/internal/filesystem"
	"gopkg.in/yaml.v3"
)

func GetSession() (string, error) {
	sessionPath := filepath.Join(filesystem.GetExecutableDir(), "session")
	data, err := os.ReadFile(sessionPath)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func LoadLocalConfig() error {
	expandendPath, err := filesystem.ExpandTilde(DefaultConfigPath)
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

	err = yaml.Unmarshal(configFileContent, &ArgsConfigInstance)
	if err != nil {
		return err
	}

	return nil
}

// MapPortToService maps a port to a service name.
func MapPortToService(port uint16) string {
	for _, service := range Current.ConfigClient.Services {
		if service.Port == port {
			return service.Name
		}
	}
	return ""
}
