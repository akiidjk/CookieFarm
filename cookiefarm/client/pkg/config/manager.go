package config

import (
	"fmt"
	"logger"
	"os"
	"path/filepath"
	"sync"
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
	return cm.cfg.host
}

func (cm *ConfigManager) SetHost(host string) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	cm.cfg.host = host
}

func (cm *ConfigManager) GetPort() uint16 {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.cfg.port
}

func (cm *ConfigManager) SetPort(port uint16) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	cm.cfg.port = port
}

func (cm *ConfigManager) GetHTTPS() bool {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.cfg.https
}

func (cm *ConfigManager) SetHTTPS(https bool) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	cm.cfg.https = https
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
	return cm.cfg.username
}

func (cm *ConfigManager) SetUsername(username string) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	cm.cfg.username = username
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
	return nil
}

func (cm *ConfigManager) Reset() error {
	return nil
}

func (cm *ConfigManager) Write() error {
	return nil
}

func (cm *ConfigManager) WriteTemplate() error {
	return nil
}
