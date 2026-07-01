package main

import (
	"errors"
	"fmt"
	"os"

	"gosteamrestarter/internal/core"
	"gosteamrestarter/internal/platform"
	"gosteamrestarter/internal/ui/cli"
)

func main() {
	app, err := bootstrap(core.Services{Platform: platform.New()})
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if err := cli.Run(app); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func bootstrap(services core.Services) (*core.App, error) {
	app := core.NewApp(services)
	if _, err := app.LoadConfig(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}
	if _, err := app.FindOrConfirmSteamPath(); err != nil && !errors.Is(err, core.ErrSteamNotFound) {
		return nil, err
	}
	return app, nil
}
