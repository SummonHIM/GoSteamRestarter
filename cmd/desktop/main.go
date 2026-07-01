package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"gosteamrestarter/internal/core"
	"gosteamrestarter/internal/platform"
	"gosteamrestarter/internal/ui/desktop"
)

func main() {
	os.Exit(runDesktop(os.Stderr, core.Services{Platform: platform.New()}, desktop.Run))
}

func runDesktop(stderr io.Writer, services core.Services, run func(*core.App) error) int {
	app, err := bootstrap(services)
	if err != nil {
		_, _ = fmt.Fprintln(stderr, err)
		return 1
	}
	if err := run(app); err != nil {
		_, _ = fmt.Fprintln(stderr, err)
		return 1
	}
	return 0
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
