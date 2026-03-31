package platform

import (
	"fmt"
	"sync"

	"cross-border-admin/config"
)

var (
	registry = make(map[string]Platform)
	mu       sync.RWMutex
)

// Register 注册平台实现
func Register(name string, platform Platform) {
	mu.Lock()
	defer mu.Unlock()
	registry[name] = platform
}

// Get 获取平台实例
func Get(name string) (Platform, error) {
	mu.RLock()
	defer mu.RUnlock()

	platform, exists := registry[name]
	if !exists {
		return nil, fmt.Errorf("平台 %s 未注册", name)
	}

	return platform, nil
}

// GetAll 获取所有已注册的平台
func GetAll() map[string]Platform {
	mu.RLock()
	defer mu.RUnlock()

	result := make(map[string]Platform)
	for k, v := range registry {
		result[k] = v
	}
	return result
}

// GetEnabled 获取所有启用的平台
func GetEnabled() map[string]Platform {
	mu.RLock()
	defer mu.RUnlock()

	result := make(map[string]Platform)
	for name, platform := range registry {
		if cfg, exists := config.GetPlatformConfig(name); exists && cfg.Enabled {
			result[name] = platform
		}
	}
	return result
}

// IsEnabled 检查平台是否启用
func IsEnabled(name string) bool {
	cfg, exists := config.GetPlatformConfig(name)
	return exists && cfg.Enabled
}
