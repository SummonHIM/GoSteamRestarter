package desktop

import "gosteamrestarter/internal/core"

func Run(app *core.App) error {
	_ = NewWindow(app)
	return nil
}
