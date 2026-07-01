package core

import "gosteamrestarter/internal/platform"

type App struct {
	services Services
	cfg      Config
}

type Services struct {
	ConfigStore ConfigStore
	Platform    platform.Interface
}

type Result struct{}

func NewApp(services Services) *App {
	return &App{services: services, cfg: DefaultConfig()}
}
