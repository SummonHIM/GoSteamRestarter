package cli

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"gosteamrestarter/internal/core"
)

func Run(in io.Reader, out io.Writer, app *core.App) error {
	reader := bufio.NewReader(in)

	for {
		RenderMainMenu(out, false)
		_, _ = fmt.Fprintln(out, "请输入选项:")

		choice, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		switch strings.TrimSpace(choice) {
		case "1":
			if err := app.KillSteam(); err != nil {
				return err
			}
			_, _ = fmt.Fprintln(out, "Steam 已结束")
		case "2":
			if err := app.RestartSteam(); err != nil {
				return err
			}
			_, _ = fmt.Fprintln(out, "Steam 已重启")
		case "3":
			if err := app.FlushDNS(); err != nil {
				return err
			}
			_, _ = fmt.Fprintln(out, "DNS 缓存已清理")
		case "4":
			if err := promptSettings(reader, out, app); err != nil {
				return err
			}
		case "0":
			_, _ = fmt.Fprintln(out, "已退出")
			return nil
		default:
			_, _ = fmt.Fprintln(out, "无效选项")
		}
	}
}
