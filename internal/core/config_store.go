package core

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type ConfigStore struct {
	dir string
}

var userConfigDir = os.UserConfigDir

func NewConfigStore(dir string) ConfigStore {
	return ConfigStore{dir: dir}
}

func ConfigDir() (string, error) {
	base, err := userConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(base, appName), nil
}

func (s ConfigStore) filePath() string {
	return filepath.Join(s.dir, configFileName)
}

func (s ConfigStore) Load() (Config, error) {
	data, err := os.ReadFile(s.filePath())
	if err != nil {
		return Config{}, err
	}

	cfg := DefaultConfig()
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func (s ConfigStore) Save(cfg Config) error {
	if err := os.MkdirAll(s.dir, 0o755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.filePath(), data, 0o644)
}
