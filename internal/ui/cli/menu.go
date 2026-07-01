package cli

import (
	"github.com/rivo/tview"
	"gosteamrestarter/internal/core"
)

func createMainMenu(tApp *tview.Application, pages *tview.Pages, app *core.App) *tview.List {
	list := tview.NewList().
		AddItem("强制结束 Steam 客户端", "终止所有 Steam 进程", '1', func() {
			if err := app.KillSteam(); err != nil {
				showMessage(tApp, pages, "menu", "错误", err.Error())
			} else {
				showMessage(tApp, pages, "menu", "成功", "Steam 已结束")
			}
		}).
		AddItem("重启 Steam 客户端", "结束并重新启动 Steam", '2', func() {
			if err := app.RestartSteam(); err != nil {
				showMessage(tApp, pages, "menu", "错误", err.Error())
			} else {
				showMessage(tApp, pages, "menu", "成功", "Steam 已重启")
			}
		}).
		AddItem("刷新 DNS 缓存", "清理系统 DNS 缓存", '3', func() {
			if err := app.FlushDNS(); err != nil {
				showMessage(tApp, pages, "menu", "错误", err.Error())
			} else {
				showMessage(tApp, pages, "menu", "成功", "DNS 缓存已清理")
			}
		}).
		AddItem("设置", "配置 Steam 路径和启动参数", '4', func() {
			pages.SwitchToPage("settings")
		}).
		AddItem("退出", "退出程序", 'q', func() {
			tApp.Stop()
		})

	list.SetTitle(" GoSteamRestarter ").SetBorder(true)
	return list
}
