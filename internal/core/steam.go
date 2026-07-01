package core

func (a *App) FindOrConfirmSteamPath() (string, error) {
	if a.cfg.SteamPath != "" {
		return a.cfg.SteamPath, nil
	}

	path, err := a.services.Platform.FindSteamPath()
	if err != nil {
		return "", err
	}
	if path == "" {
		return "", ErrSteamNotFound
	}

	cfg := a.cfg
	cfg.SteamPath = path
	if err := a.SaveConfig(cfg); err != nil {
		return "", err
	}
	return path, nil
}

func (a *App) RestartSteam() error {
	if err := a.KillSteam(); err != nil {
		return err
	}
	return a.StartSteam()
}

func (a *App) KillSteam() error {
	return a.services.Platform.KillSteam()
}

func (a *App) StartSteam() error {
	cfg := a.cfg
	if cfg.SteamPath == "" {
		return ErrSteamNotFound
	}
	return a.services.Platform.StartSteam(cfg.SteamPath, cfg.SteamArgs)
}

func (a *App) FlushDNS() error {
	return a.services.Platform.FlushDNS()
}
