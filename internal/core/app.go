package core

func (a *App) SetConfigStore(store ConfigStore) {
	a.services.ConfigStore = store
}

func (a *App) ConfigStore() ConfigStore {
	store, err := a.resolveConfigStore()
	if err != nil {
		return a.services.ConfigStore
	}
	return store
}

func (a *App) LoadConfig() (Config, error) {
	store, err := a.resolveConfigStore()
	if err != nil {
		return Config{}, err
	}
	cfg, err := store.Load()
	if err != nil {
		return Config{}, err
	}
	a.cfg = cfg
	return cfg, nil
}

func (a *App) SaveConfig(cfg Config) error {
	store, err := a.resolveConfigStore()
	if err != nil {
		return err
	}
	if err := store.Save(cfg); err != nil {
		return err
	}
	a.cfg = cfg
	return nil
}

func (a *App) resolveConfigStore() (ConfigStore, error) {
	if a.services.ConfigStore != (ConfigStore{}) {
		return a.services.ConfigStore, nil
	}
	return a.defaultConfigStore()
}

func (a *App) defaultConfigStore() (ConfigStore, error) {
	dir, err := ConfigDir()
	if err != nil {
		return ConfigStore{}, err
	}
	return NewConfigStore(dir), nil
}

func (a *App) GetConfig() Config {
	return a.cfg
}
