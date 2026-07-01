package desktop

import "gosteamrestarter/internal/core"

type Callbacks struct {
	Restart      func() error
	Kill         func() error
	FlushDNS     func() error
	SaveSettings func(core.Config) error
}

type Window struct {
	App       *core.App
	Callbacks Callbacks
}

func NewWindow(app *core.App) Window {
	window := Window{App: app}
	if app == nil {
		return window
	}
	window.Callbacks = Callbacks{
		Restart: func() error {
			return app.RestartSteam()
		},
		Kill: func() error {
			return app.KillSteam()
		},
		FlushDNS: func() error {
			return app.FlushDNS()
		},
		SaveSettings: func(cfg core.Config) error {
			return app.SaveConfig(cfg)
		},
	}
	return window
}
