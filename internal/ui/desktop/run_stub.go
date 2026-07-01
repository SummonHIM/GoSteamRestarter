//go:build nofyne

package desktop

import "gosteamrestarter/internal/core"

func Run(coreApp *core.App) error {
	if coreApp == nil {
		return nil
	}
	_ = NewWindow(coreApp)
	return nil
}
