// Package config provides functions to manage the CookieFarm client configuration globally.
package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ByteTheCookies/cookieclient/internal/filesystem"
	"github.com/ByteTheCookies/cookieclient/internal/logger"
	"gopkg.in/yaml.v3"
)

var (
	Token              string     // Token stores the global authentication token.
	ArgsAttackInstance ArgsAttack // Struct holding runtime arguments
	ArgsConfigInstance ArgsConfig // Struct holding configuration arguments
	Current            Config     // Current holds the global configuration for the client.
	UseTUI             bool       // UseTUI indicates whether to use the TUI mode or not
	PID                int        // PID is the process ID of the running exploit
	ExploitName        string     // ExploitName is the name of the exploit being run
	UseBanner          bool       // NoBanner indicates whether to disable the banner on startup
)

var DefaultConfigPath, _ = filesystem.ExpandTilde("~/.config/cookiefarm")

var ExploitTemplate = []byte(`#!/usr/bin/env python3
from cookiefarm_exploiter import exploit_manager

@exploit_manager
def exploit(ip, port, name_service):
    # Run your exploit here
    flag = ""

    # Just print the flag to stdout
    print(flag)
`)

var ConfigTemplate = []byte(`address: "localhost"
port: 8080
https: false
username: "cookieguest"
`)

// ResetConfigFunc resets the configuration to defaults
func GetConfig() (string, error) {
	var err error
	err = os.MkdirAll(DefaultConfigPath, 0o755)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Error creating config directory")
		return "", err
	}

	configPath := filepath.Join(DefaultConfigPath, "config.yml")

	file, err := os.Create(configPath)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Error opening configuration file")
		return "", err
	}
	defer file.Close()

	err = yaml.Unmarshal(ConfigTemplate, &ArgsConfigInstance)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Error unmarshalling default configuration")
		return "", err
	}

	err = yaml.NewEncoder(file).Encode(ArgsConfigInstance)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Error encoding configuration to YAML")
		return "", err
	}

	return "Config reset successfully", nil
}

// Reset resets the configuration to defaults
func Reset() (string, error) {
	var err error
	err = os.MkdirAll(DefaultConfigPath, 0o755)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Error creating config directory")
		return "", err
	}

	configPath := filepath.Join(DefaultConfigPath, "config.yml")

	file, err := os.Create(configPath)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Error opening configuration file")
		return "", err
	}
	defer file.Close()

	err = yaml.Unmarshal(ConfigTemplate, &ArgsConfigInstance)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Error unmarshalling default configuration")
		return "", err
	}

	err = yaml.NewEncoder(file).Encode(ArgsConfigInstance)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Error encoding configuration to YAML")
		return "", err
	}

	return "Config reset successfully", nil
}

// Update updates the configuration with new values
func Update(configuration ArgsConfig) (string, error) {
	var err error
	err = os.MkdirAll(DefaultConfigPath, 0o755)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Error creating config directory")
		return "", err
	}

	configPath := filepath.Join(DefaultConfigPath, "config.yml")

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		logger.Log.Warn().Msg("Configuration file does not exist, creating a new one with default settings")
		os.WriteFile(configPath, ConfigTemplate, 0o644)
	} else if err != nil {
		logger.Log.Error().Err(err).Msg("Error checking configuration file")
		return "", err
	}

	file, err := os.Create(configPath)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Error creating or opening configuration file")
		return "", err
	}
	defer file.Close()

	err = yaml.NewEncoder(file).Encode(configuration)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Error encoding configuration to YAML")
		return "", err
	}

	return configPath, nil
}

// Logout handles user logout
func Logout() (string, error) {
	sessionPath := filepath.Join(DefaultConfigPath, "session")
	err := os.Remove(sessionPath)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Error removing session file")
		return "", err
	}
	return "Logout successfully", nil
}

// Show displays the current configuration
func Show() (string, error) {
	configPath := filepath.Join(DefaultConfigPath, "config.yml")

	content, err := os.ReadFile(configPath)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Error reading configuration file")
		return "", fmt.Errorf("%s", "Error reading configuration file: "+err.Error())
	}

	return string(content), nil
}
