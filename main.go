package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/atotto/clipboard"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/shirou/gopsutil/process"
	"gopkg.in/yaml.v3"
)

type Config struct {
	App struct {
		Location  string `yaml:"location"`
		Arguments string `yaml:"arguments"`
	} `yaml:"app"`
}

var AppVersion string = "development"

var WEBSITE string = "https://github.com/SummonHIM"

var CONFIG Config

var CONFIG_PATH string

const APP_NAME string = "Steam"

var APP_PROC_NAME string = "steam"

var APP_PROC_LOCATION string = APP_PROC_NAME

// 使用 PsUtil 查询进程信息。
func process_search(name string) (*process.Process, error) {
	// 获取所有进程
	processes, err := process.Processes()
	if err != nil {
		return nil, fmt.Errorf("错误：获取进程列表失败: %v", err)
	}

	// 遍历所有进程，查找目标进程
	for _, proc := range processes {
		procName, err := proc.Name()
		if err == nil && strings.EqualFold(procName, name) { // 忽略大小写匹配
			return proc, nil
		}
	}

	return nil, fmt.Errorf("错误：未找到进程: %s", name)
}

// 杀死指定进程
func process_kill(name string, output *tview.TextView) bool {
	// 查询进程信息
	fmt.Fprintf(output, "正在查询 %s 的相关进程信息…\n", name)
	proc, err := process_search(name)
	if err != nil {
		fmt.Fprintf(output, "%v\n", err)
		return false
	} else {
		fmt.Fprintf(output, "找到进程 %s (PID: %d)。准备杀死…\n", name, proc.Pid)
	}

	// 查找进程具体路径，存入 PROC_LOCATION 中
	exePath, err := proc.Exe()
	if err != nil {
		fmt.Fprintf(output, "警告：无法获取进程 %s 的文件具体位置。\n", name)
	} else {
		fmt.Fprintf(output, "顺便查询到 %s 的具体位置：%s\n", name, exePath)
		APP_PROC_LOCATION = exePath
	}

	// 杀死进程
	err = proc.Kill()
	if err != nil {
		fmt.Fprintf(output, "错误：尝试杀死进程失败: %v\n", err)
		return false
	} else {
		fmt.Fprintf(output, "进程 %s 已被成功杀死。\n", name)
		return true
	}
}

func process_start(path string, output *tview.TextView) {
	args := strings.Fields(CONFIG.App.Arguments)
	fmt.Fprintf(output, "正在尝试启动 %s %s …\n", path, CONFIG.App.Arguments)
	cmd := exec.Command(path, args...)

	err := cmd.Start()
	if err != nil {
		fmt.Fprintf(output, "错误：启动 %s 失败: %v", path, err)
	} else {
		fmt.Fprintf(output, "启动 %s 成功！", path)
	}
}

// 清理 DNS 缓存
func flush_dns(output *tview.TextView) {
	fmt.Fprintf(output, "正在尝试清理系统 DNS 缓存...\n")

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("ipconfig", "/flushdns")
	case "darwin": // macOS
		cmd = exec.Command("sudo", "dscacheutil", "-flushcache")
	case "linux":
		cmd = exec.Command("sudo", "systemctl", "restart", "systemd-resolved")
	default:
		fmt.Fprintf(output, "错误：不支持的操作系统\n")
		return
	}

	// 执行命令并捕获输出
	outputBytes, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Fprintf(output, "错误：清理 DNS 缓存失败: %v\n", err)
	} else {
		fmt.Fprintf(output, "成功清理 DNS 缓存:\n%s\n", string(outputBytes))
	}
}

// 输出页面
func render_output_page(app *tview.Application, name string, action func(app *tview.Application, output *tview.TextView, done chan bool)) {
	// 创建一个 TextView 来显示输出
	output := tview.NewTextView().
		SetDynamicColors(true). // 支持动态颜色
		SetScrollable(true).    // 支持滚动
		SetChangedFunc(func() {
			// 自动滚动到底部
			app.Draw()
		})

	// 设置边框和标题
	output.SetBorder(true).SetTitle(" " + name + " ")

	// 创建一个 Flex 布局，初始只有输出区域
	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(output, 0, 1, true) // 输出区域

	// 设置 Flex 为当前页面
	app.SetRoot(flex, true)

	// 创建一个信号通道，用于通知 action 完成
	done := make(chan bool, 1)

	// 执行操作
	go func() {
		action(app, output, done) // 执行传入的 action
		done <- true              // 通知完成
	}()

	// 监听完成信号并动态添加返回按钮
	go func() {
		<-done // 等待 action 完成
		app.QueueUpdateDraw(func() {
			// 在 Flex 中添加返回按钮
			return_button := tview.NewButton("按任意键返回").SetSelectedFunc(func() {
				app.SetRoot(render_main_menu_page(app), true)
			})

			flex.AddItem(return_button, 3, 1, false)

			app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
				app.SetRoot(render_main_menu_page(app), true)
				return nil
			})
		})
	}()
}

// 渲染选项1： 强制结束客户端
func render_kill_process(app *tview.Application) {
	render_output_page(app, "强制结束 "+APP_NAME+" 客户端", func(app *tview.Application, output *tview.TextView, done chan bool) {
		process_kill(APP_PROC_NAME, output)
	})
}

// 渲染选项2： 强制重启客户端
func render_restart_process(app *tview.Application) {
	render_output_page(app, "强制重启 "+APP_NAME+" 客户端", func(app *tview.Application, output *tview.TextView, done chan bool) {
		if process_kill(APP_PROC_NAME, output) {
			if CONFIG.App.Location != "" {
				process_start(CONFIG.App.Location, output)
			} else {
				process_start(APP_PROC_LOCATION, output)
			}
		}
	})
}

// 渲染选项3： 清理 DNS 缓存
func render_flush_dns(app *tview.Application) {
	render_output_page(app, "清理系统 DNS 缓存", func(app *tview.Application, output *tview.TextView, done chan bool) {
		flush_dns(output)
	})
}

// 显示编辑设置项对话框
func show_modal(app *tview.Application, title string, initialValue string, onSave func(newValue string)) {
	// 创建输入框
	inputField := tview.NewInputField().
		SetText(initialValue).
		SetLabel(title + ": ")

	when_save := func() {
		newValue := inputField.GetText()
		onSave(newValue)   // 调用保存回调函数
		render_config(app) // 返回配置界面
	}
	when_cancel := func() {
		render_config(app) // 返回配置界面
	}

	btn_save := tview.NewButton("(C-s) 保存").SetSelectedFunc(when_save)
	btn_cancel := tview.NewButton("(C-q) 取消").SetSelectedFunc(when_cancel)

	// 创建 Modal
	modal := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(inputField, 0, 1, true).
		AddItem(tview.NewFlex().
			AddItem(btn_save, 0, 1, false).
			AddItem(btn_cancel, 0, 1, false), 3, 1, false)

	// 设置 Modal 为当前页面
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlS:
			when_save()
		case tcell.KeyCtrlQ:
			when_cancel()
		case tcell.KeyTAB:
			switch {
			case inputField.HasFocus():
				app.SetFocus(btn_save)
			case btn_save.HasFocus():
				app.SetFocus(btn_cancel)
			case btn_cancel.HasFocus():
				app.SetFocus(inputField)
			default:
				app.SetFocus(inputField)
			}
		default:
			return event
		}
		return event
	})
	app.SetRoot(modal, true)
}

func render_config(app *tview.Application) {
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		return event
	})

	originalConfig := CONFIG
	var about_list *tview.List

	about_list = tview.NewList().
		AddItem("Steam 启动路径", "（留空自动识别）  "+CONFIG.App.Location, '1', func() {
			// 弹出编辑 Modal 界面
			show_modal(app, "编辑 Steam 启动路径", CONFIG.App.Location, func(newValue string) {
				CONFIG.App.Location = newValue
				// 更新列表项
				about_list.SetItemText(0, "Steam 启动路径", "（留空自动识别）  "+CONFIG.App.Location)
			})
		}).
		AddItem("Steam 启动选项", CONFIG.App.Arguments, '2', func() {
			// 弹出编辑 Modal 界面
			show_modal(app, "编辑 Steam 启动选项", CONFIG.App.Arguments, func(newValue string) {
				CONFIG.App.Arguments = newValue
				// 更新列表项
				about_list.SetItemText(1, "Steam 启动选项", CONFIG.App.Arguments)
			})
		}).
		AddItem("保存并退出", "", 's', func() {
			// 检查启动路径是否是文件
			if _, err := os.Stat(CONFIG.App.Location); CONFIG.App.Location != "" && err != nil {
				show_error_modal(app, fmt.Sprintf("找不到 %s 这个文件。", CONFIG.App.Location))
				return
			}

			// 保存配置到 YAML 文件
			file, err := os.Create(CONFIG_PATH) // 创建/覆盖文件
			if err != nil {
				show_error_modal(app, fmt.Sprintf("保存失败: %v", err))
				return
			}
			defer file.Close()

			// 将配置写入文件
			encoder := yaml.NewEncoder(file)
			defer encoder.Close()
			if err := encoder.Encode(&CONFIG); err != nil {
				// 显示错误信息
				show_error_modal(app, fmt.Sprintf("保存失败: %v", err))
				return
			}

			// 显示保存成功信息并返回主菜单
			modal := tview.NewModal().
				SetText("保存成功！").
				AddButtons([]string{"确定"}).
				SetDoneFunc(func(buttonIndex int, buttonLabel string) {
					app.SetRoot(render_main_menu_page(app), true)
				})
			app.SetRoot(modal, true)
		}).
		AddItem("撤销并退出", "", 'q', func() {
			// 恢复原始配置
			CONFIG = originalConfig
			app.SetRoot(render_main_menu_page(app), true)
		})
	about_list.SetBorder(true).SetTitle(" 选项 ")

	app.SetRoot(about_list, true)
}

func render_about(app *tview.Application) {
	about_list := tview.NewList().
		AddItem("版本号", AppVersion, '1', func() {
			if err := clipboard.WriteAll(AppVersion); err != nil {
				show_error_modal(app, fmt.Sprintf("复制文本时发生错误：%s", err))
				return
			}
		}).
		AddItem("网站", WEBSITE, '2', func() {
			if err := clipboard.WriteAll(WEBSITE); err != nil {
				show_error_modal(app, fmt.Sprintf("复制文本时发生错误：%s", err))
				return
			}
		})
	about_list.SetBorder(true).SetTitle(" 关于一键强制关闭 " + APP_NAME + " 客户端 ")

	return_button := tview.NewButton("(q) 返回").SetSelectedFunc(func() {
		app.SetRoot(render_main_menu_page(app), true)
	})

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(about_list, 0, 1, true).
		AddItem(return_button, 3, 1, false)

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTAB {
			switch {
			case about_list.HasFocus():
				app.SetFocus(return_button)
			case return_button.HasFocus():
				app.SetFocus(about_list)
			default:
				app.SetFocus(return_button)
			}
		}
		if event.Rune() == 'q' {
			app.SetRoot(render_main_menu_page(app), true)
		}
		return event
	})

	app.SetRoot(flex, true)
}

// 渲染主菜单页面
func render_main_menu_page(app *tview.Application) *tview.Flex {
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		return event
	})

	list := tview.NewList().
		AddItem("强制结束 "+APP_NAME+" 客户端", "", '1', func() { render_kill_process(app) }).
		AddItem("强制重启 "+APP_NAME+" 客户端", "", '2', func() { render_restart_process(app) }).
		AddItem("清理系统 DNS 缓存", "", '3', func() { render_flush_dns(app) }).
		AddItem("选项", "", '4', func() { render_config(app) }).
		AddItem("关于", "", '5', func() { render_about(app) })

	flex := tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(list, 0, 12, true).
			AddItem(nil, 0, 1, false), 0, 2, true).
		AddItem(nil, 0, 1, false)
	flex.SetBorder(true).
		SetTitle(" 一键强制关闭 " + APP_NAME + " 客户端 ")

	return flex
}

// 错误对话框
func show_error_modal(app *tview.Application, errorMessage string) {
	modal := tview.NewModal().
		SetText(errorMessage).
		AddButtons([]string{"确定"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			render_config(app) // 返回配置界面
		})
	app.SetRoot(modal, true)
}

// 渲染页面
func render() {
	app := tview.NewApplication()

	if err := app.SetRoot(render_main_menu_page(app), true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

func main() {
	if runtime.GOOS == "windows" {
		APP_PROC_NAME = "steam.exe"
		APP_PROC_LOCATION = "steam.exe"
	}

	CONFIG_PATH = xdg.ConfigHome + "/GoSteamRestarter/config.yaml"
	config_file, err := os.Open(CONFIG_PATH)
	if err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(filepath.Dir(CONFIG_PATH), os.ModePerm); err != nil {
				fmt.Printf("错误：无法创建存放配置的文件夹。%v\n", err)
			}

			default_config := Config{}
			default_config.App.Location = ""
			default_config.App.Arguments = ""

			yaml_data, err := yaml.Marshal(&default_config)
			if err != nil {
				fmt.Printf("初始化默认配置失败: %v\n", err)
				return
			}

			create_config, err := os.Create(CONFIG_PATH)
			if err != nil {
				fmt.Printf("创建配置文件失败: %v\n", err)
				return
			}
			defer create_config.Close()

			if _, err := create_config.Write(yaml_data); err != nil {
				fmt.Printf("写入默认配置失败: %v\n", err)
				return
			}
		} else {
			fmt.Printf("错误：无法打开配置文件。%v\n", err)
			return
		}
	} else {
		// 确保文件被正确关闭
		defer config_file.Close()

		// 使用 yaml 解码器解析文件
		decoder := yaml.NewDecoder(config_file)
		if err := decoder.Decode(&CONFIG); err != nil {
			fmt.Printf("错误：该 YAML 文件无法解读 。%v\n", err)
			return
		}
	}

	render()
}
