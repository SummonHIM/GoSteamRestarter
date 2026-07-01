package cli

import (
	"fmt"

	"github.com/rivo/tview"
	"gosteamrestarter/internal/core"
)

func Run(app *core.App) error {
	tApp := tview.NewApplication()
	pages := tview.NewPages()

	mainMenu := createMainMenu(tApp, pages, app)
	pages.AddPage("menu", mainMenu, true, true)

	settingsForm := createSettingsForm(tApp, pages, app)
	pages.AddPage("settings", settingsForm, true, false)

	return tApp.SetRoot(pages, true).EnableMouse(true).Run()
}

func showMessage(tApp *tview.Application, pages *tview.Pages, returnPage, title, message string) {
	modal := tview.NewModal().
		SetText(fmt.Sprintf("[%s]\n%s", title, message)).
		AddButtons([]string{"确定"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			pages.RemovePage("modal")
			pages.SwitchToPage(returnPage)
		})
	pages.AddAndSwitchToPage("modal", modal, true)
}
