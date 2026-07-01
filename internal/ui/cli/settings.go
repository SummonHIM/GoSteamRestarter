package cli

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"gosteamrestarter/internal/core"
)

func promptSettings(in *bufio.Reader, out io.Writer, app *core.App) error {
	_, _ = fmt.Fprintln(out, "设置")
	_, _ = fmt.Fprintln(out, "请输入 Steam 路径:")
	path, err := in.ReadString('\n')
	if err != nil {
		return err
	}
	_, _ = fmt.Fprintln(out, "请输入 Steam 启动参数:")
	args, err := in.ReadString('\n')
	if err != nil {
		return err
	}

	cfg := core.Config{
		SteamPath: strings.TrimSpace(path),
		SteamArgs: strings.TrimSpace(args),
	}
	if err := app.SaveConfig(cfg); err != nil {
		return err
	}
	_, _ = fmt.Fprintln(out, "设置已保存")
	return nil
}
