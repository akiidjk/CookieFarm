package config

import (
	"fmt"
	"logger"
	"os"
	"path/filepath"
	"sync"

	"gopkg.in/yaml.v3"
)

var (
	instance *ConfigManager
	once     sync.Once
)

const configTemplate = `
host: "localhost"
port: 8080
https: false
username: "cookieguest"
`

func GetInstance() *ConfigManager {
	once.Do(func() {
		instance = &ConfigManager{}
	})
	return instance
}

func (cm *ConfigManager) GetHost() string {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.cfg.Host
}

func (cm *ConfigManager) SetHost(host string) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	cm.cfg.Host = host
}

func (cm *ConfigManager) GetPort() uint16 {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.cfg.Port
}

func (cm *ConfigManager) SetPort(port uint16) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	cm.cfg.Port = port
}

func (cm *ConfigManager) GetHTTPS() bool {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.cfg.HTTPS
}

func (cm *ConfigManager) SetHTTPS(https bool) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	cm.cfg.HTTPS = https
}

func (cm *ConfigManager) GetToken() string {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.token
}

func (cm *ConfigManager) SetToken(token string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.token = token
}

func (cm *ConfigManager) GetUsername() string {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.cfg.Username
}

func (cm *ConfigManager) SetUsername(username string) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	cm.cfg.Username = username
}

func (cm *ConfigManager) GetConfig() Config {
	return cm.cfg
}

func (cm *ConfigManager) GetSession() (string, error) {
	cm.mu.RLock()
	sessionPath := filepath.Join(DefaultPath, "session")
	cm.mu.RUnlock()

	data, err := os.ReadFile(sessionPath)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (cm *ConfigManager) Logout() (string, error) {
	cm.mu.RLock()
	sessionPath := filepath.Join(DefaultPath, "session")
	cm.mu.RUnlock()

	err := os.Remove(sessionPath)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Error removing session file")
		return "", err
	}

	return "Logout successfully", nil
}

func (cm *ConfigManager) ShowLocalConfigContent() (string, error) {
	cm.mu.RLock()
	configPath := filepath.Join(DefaultPath, "config.yml")
	cm.mu.RUnlock()

	content, err := os.ReadFile(configPath)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Error reading configuration file")
		return "", fmt.Errorf("error reading configuration file: %w", err)
	}

	return string(content), nil
}

func (cm *ConfigManager) MapPortToService(port uint16) string {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	for service, serviceport := range cm.cfg.services {
		if serviceport == port {
			return service
		}
	}

	return ""
}

func (cm *ConfigManager) MapServiceToPort(serviceName string) uint16 {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	port, ok := cm.cfg.services[serviceName]
	if !ok {
		return 0
	}

	return port
}

func (cm *ConfigManager) Read() error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if _, err := os.Stat(DefaultPath); os.IsNotExist(err) {
		logger.Log.Warn().Msgf("Config directory does not exist at %s, creating it", DefaultPath)
		err = os.MkdirAll(DefaultPath, 0o755)
		if err != nil {
			return fmt.Errorf("error creating config directory: %w", err)
		}

		if _, err := os.Create(filepath.Join(DefaultPath, "config.yml")); err != nil {
			return fmt.Errorf("error creating default config file: %w", err)
		}
	}

	configFileContent, err := os.ReadFile(filepath.Join(DefaultPath, "config.yml"))
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("config file does not exist at %s", DefaultPath)
		}
		return fmt.Errorf("error reading config file: %w", err)
	}

	var tmp Config
	err = yaml.Unmarshal(configFileContent, &tmp)
	if err != nil {
		return fmt.Errorf("error unmarshalling config: %w", err)
	}

	cm.mu.Unlock()
	cm.SetHTTPS(tmp.HTTPS)
	cm.SetHost(tmp.Host)
	cm.SetUsername(tmp.Username)
	cm.SetPort(tmp.Port)
	cm.mu.Lock()

	return nil
}

func (cm *ConfigManager) Reset() error {
	err := cm.Read()
	if err != nil {
		logger.Log.Warn().Err(err).Msg("Could not load existing config, proceeding with empty exploits")
	}
	cm.mu.Lock()
	defer cm.mu.Unlock()

	err = os.MkdirAll(DefaultPath, 0o755)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Error creating config directory")
		return err
	}

	configPath := filepath.Join(DefaultPath, "config.yml")

	file, err := os.Create(configPath)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Error opening configuration file")
		return err
	}
	defer file.Close()

	var tmp Config
	err = yaml.Unmarshal([]byte(configTemplate), &tmp)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Error unmarshalling default configuration")
		return err
	}

	cm.mu.Unlock()
	cm.SetHTTPS(tmp.HTTPS)
	cm.SetHost(tmp.Host)
	cm.SetUsername(tmp.Username)
	cm.SetPort(tmp.Port)
	cm.mu.Lock()

	err = yaml.NewEncoder(file).Encode(cm.GetConfig())
	if err != nil {
		logger.Log.Error().Err(err).Msg("Error encoding configuration to YAML")
		return err
	}

	return nil
}

func (cm *ConfigManager) Write() error {
	cm.mu.RLock()
	config := cm.GetConfig()
	cm.mu.RUnlock()

	configFilePath := filepath.Join(DefaultPath, "config.yml")
	configFileContent, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("error marshalling config: %w", err)
	}

	err = os.WriteFile(configFilePath, configFileContent, 0o644)
	if err != nil {
		return fmt.Errorf("error writing config file: %w", err)
	}

	return nil
}

func (cm *ConfigManager) WriteTemplate() error {
	return nil
}
