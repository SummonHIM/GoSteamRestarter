//go:build !nofyne

package desktop

import (
	"gosteamrestarter/internal/core"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func Run(coreApp *core.App) error {
	if coreApp == nil {
		return nil
	}

	a := app.New()
	w := a.NewWindow("GoSteamRestarter")
	w.Resize(fyne.NewSize(400, 300))

	killBtn := widget.NewButton("强制结束 Steam", func() {
		if err := coreApp.KillSteam(); err != nil {
			dialog.ShowError(err, w)
		} else {
			dialog.ShowInformation("成功", "Steam 已结束", w)
		}
	})

	restartBtn := widget.NewButton("重启 Steam", func() {
		if err := coreApp.RestartSteam(); err != nil {
			dialog.ShowError(err, w)
		} else {
			dialog.ShowInformation("成功", "Steam 已重启", w)
		}
	})

	flushBtn := widget.NewButton("刷新 DNS", func() {
		if err := coreApp.FlushDNS(); err != nil {
			dialog.ShowError(err, w)
		} else {
			dialog.ShowInformation("成功", "DNS 缓存已清理", w)
		}
	})

	settingsBtn := widget.NewButton("设置", func() {
		showSettings(w, coreApp)
	})

	content := container.NewVBox(killBtn, restartBtn, flushBtn, settingsBtn)
	w.SetContent(content)
	w.ShowAndRun()
	return nil
}

func showSettings(w fyne.Window, coreApp *core.App) {
	cfg := coreApp.GetConfig()

	pathEntry := widget.NewEntry()
	pathEntry.SetText(cfg.SteamPath)
	pathEntry.SetPlaceHolder("Steam 可执行文件路径")

	argsEntry := widget.NewEntry()
	argsEntry.SetText(cfg.SteamArgs)
	argsEntry.SetPlaceHolder("启动参数（可选）")

	items := []*widget.FormItem{
		widget.NewFormItem("Steam 路径", pathEntry),
		widget.NewFormItem("Steam 参数", argsEntry),
	}

	d := dialog.NewForm("设置", "保存", "取消", items, func(ok bool) {
		if ok {
			newCfg := core.Config{
				SteamPath: pathEntry.Text,
				SteamArgs: argsEntry.Text,
			}
			if err := coreApp.SaveConfig(newCfg); err != nil {
				dialog.ShowError(err, w)
			} else {
				dialog.ShowInformation("成功", "设置已保存", w)
			}
		}
	}, w)
	d.Resize(fyne.NewSize(400, 200))
	d.Show()
}
