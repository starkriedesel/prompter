package prompter

import (
	"encoding/json"
	"sync"
)

// Config is used as a key/value store to persist state between commands.
type Config map[string]string

var currentConfig Config
var once sync.Once

// GetConfig returns the config. Creates it when it's called for the first time.
func GetConfig() Config {
	once.Do(func() {
		cfg := make(Config)
		currentConfig = cfg
	})
	return currentConfig
}

// Export marshals the config to JSON.
func (c Config) Export() ([]byte, error) {
	exported, err := json.Marshal(c)
	if err != nil {
		return []byte{}, err
	}
	return exported, nil
}

// ImportConfig converts a serialized JSON into the config.
func ImportConfig(input []byte) error {
	importedConfig := make(Config)
	err := json.Unmarshal(input, &importedConfig)
	if err != nil {
		return err
	}
	GetConfig() // force the singleton to be generated if its the first time
	currentConfig = importedConfig
	return nil
}

// Set assigns the value to the key in the config. It overwrites any previous
// values so if needed, check with Exists.
func (c Config) Set(key, value string) {
	c[key] = value
}

// Key returns the value of a key or "" if it does not exist in the config.
func (c Config) Key(key string) string {
	return c[key]
}

// Contains returns true if a key exists in the config.
func (c Config) Contains(key string) bool {
	_, exists := c[key]
	return exists
}
