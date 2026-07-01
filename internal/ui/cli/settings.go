package cli

import (
	"github.com/rivo/tview"
	"gosteamrestarter/internal/core"
)

func createSettingsForm(tApp *tview.Application, pages *tview.Pages, app *core.App) *tview.Form {
	cfg := app.GetConfig()

	var form *tview.Form
	form = tview.NewForm().
		AddInputField("Steam 路径", cfg.SteamPath, 60, nil, nil).
		AddInputField("Steam 参数", cfg.SteamArgs, 60, nil, nil).
		AddButton("保存", func() {
			pathItem := form.GetFormItemByLabel("Steam 路径").(*tview.InputField)
			argsItem := form.GetFormItemByLabel("Steam 参数").(*tview.InputField)
			newCfg := core.Config{
				SteamPath: pathItem.GetText(),
				SteamArgs: argsItem.GetText(),
			}
			if err := app.SaveConfig(newCfg); err != nil {
				showMessage(tApp, pages, "settings", "错误", err.Error())
			} else {
				showMessage(tApp, pages, "settings", "成功", "设置已保存")
			}
		}).
		AddButton("返回", func() {
			pages.SwitchToPage("menu")
		})

	form.SetTitle(" 设置 ").SetBorder(true)
	return form
}
