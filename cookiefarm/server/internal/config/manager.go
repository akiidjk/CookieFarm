package config

import (
	"sharedconfig"
	"sync"
)

var (
	instance *ConfigManager
	once     sync.Once
)

func GetInstance() *ConfigManager {
	once.Do(func() {
		instance = &ConfigManager{}
	})
	return instance
}

func (cm *ConfigManager) GetURLFlagChecker() string {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.cfg.Server.URLFlagChecker
}

func (cm *ConfigManager) SetURLFlagChecker(v string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.cfg.Server.URLFlagChecker = v
}

func (cm *ConfigManager) GetTeamToken() string {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.cfg.Server.TeamToken
}

func (cm *ConfigManager) SetTeamToken(v string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.cfg.Server.TeamToken = v
}

func (cm *ConfigManager) GetSubmitFlagCheckerTime() uint {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.cfg.Server.SubmitFlagCheckerTime
}

func (cm *ConfigManager) SetSubmitFlagCheckerTime(v uint) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.cfg.Server.SubmitFlagCheckerTime = v
}

func (cm *ConfigManager) GetMaxFlagBatchSize() uint {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.cfg.Server.MaxFlagBatchSize
}

func (cm *ConfigManager) SetMaxFlagBatchSize(v uint) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.cfg.Server.MaxFlagBatchSize = v
}

func (cm *ConfigManager) GetProtocol() string {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.cfg.Server.Protocol
}

func (cm *ConfigManager) SetProtocol(v string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.cfg.Server.Protocol = v
}

func (cm *ConfigManager) GetTickTime() uint {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.cfg.Server.TickTime
}

func (cm *ConfigManager) SetTickTime(v uint) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.cfg.Server.TickTime = v
}

func (cm *ConfigManager) GetFlagTTL() uint64 {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.cfg.Server.FlagTTL
}

func (cm *ConfigManager) SetFlagTTL(v uint64) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.cfg.Server.FlagTTL = v
}

func (cm *ConfigManager) GetStartTime() string {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.cfg.Server.StartTime
}

func (cm *ConfigManager) SetStartTime(v string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.cfg.Server.StartTime = v
}

func (cm *ConfigManager) GetEndTime() string {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.cfg.Server.EndTime
}

func (cm *ConfigManager) SetEndTime(v string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.cfg.Server.EndTime = v
}

func (cm *ConfigManager) SetConfigured(value bool) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.cfg.Configured = value
}

func (cm *ConfigManager) GetConfigured() bool {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.cfg.Configured
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

func (cm *ConfigManager) GetShared() sharedconfig.Shared {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.cfg.Shared
}

func (cm *ConfigManager) SetShared(v sharedconfig.Shared) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.cfg.Shared = v
}

func (cm *ConfigManager) SetConfig(c Config) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.cfg.Server = c
}

func (cm *ConfigManager) GetConfig() Config {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	return cm.cfg.Server
}

func (cm *ConfigManager) GetFullConfig() FullConfig {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.cfg
}

func (cm *ConfigManager) SetFullConfig(f FullConfig) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.cfg = f
}
