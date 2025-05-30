package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

func GetSession() (string, error) {
	sessionPath := filepath.Join(DefaultConfigPath, "session")
	data, err := os.ReadFile(sessionPath)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func LoadLocalConfig() error {
	configFileContent, err := os.ReadFile(DefaultConfigPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("config file does not exist at %s", DefaultConfigPath)
		}
		return fmt.Errorf("error reading config file: %w", err)
	}

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
