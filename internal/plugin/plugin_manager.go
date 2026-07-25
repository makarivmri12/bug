package plugin

import (
	"fmt"
	"plugin"
	"sync"

	"go.uber.org/zap"
)

// PluginManager manages the loading and execution of scanner plugins
type PluginManager struct {
	plugins map[string]interface{}
	mu      sync.RWMutex
	logger  *zap.Logger
}

// NewPluginManager creates a new plugin manager
func NewPluginManager(logger *zap.Logger) *PluginManager {
	return &PluginManager{
		plugins: make(map[string]interface{}),
		logger:  logger,
	}
}

// LoadPlugin dynamically loads a plugin from a shared object file
func (pm *PluginManager) LoadPlugin(path string, name string) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	p, err := plugin.Open(path)
	if err != nil {
		pm.logger.Error("Failed to open plugin", zap.Error(err), zap.String("path", path))
		return err
	}

	symbol, err := p.Lookup("Plugin")
	if err != nil {
		pm.logger.Error("Failed to lookup plugin symbol", zap.Error(err), zap.String("name", name))
		return err
	}

	pm.plugins[name] = symbol
	pm.logger.Info("Plugin loaded", zap.String("name", name), zap.String("path", path))

	return nil
}

// GetPlugin retrieves a loaded plugin
func (pm *PluginManager) GetPlugin(name string) (interface{}, error) {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	if p, ok := pm.plugins[name]; ok {
		return p, nil
	}

	return nil, fmt.Errorf("plugin not found: %s", name)
}

// ListPlugins returns a list of loaded plugins
func (pm *PluginManager) ListPlugins() []string {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	var names []string
	for name := range pm.plugins {
		names = append(names, name)
	}

	return names
}

// UnloadPlugin removes a plugin from the manager
func (pm *PluginManager) UnloadPlugin(name string) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	delete(pm.plugins, name)
	pm.logger.Info("Plugin unloaded", zap.String("name", name))
}
