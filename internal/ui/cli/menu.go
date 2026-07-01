package cli

import (
	"fmt"
	"io"
)

func RenderMainMenu(w io.Writer, admin bool) {
	_, _ = fmt.Fprintln(w, "GoSteamRestarter")
	_, _ = fmt.Fprintln(w, "1. 强制结束 Steam 客户端")
	_, _ = fmt.Fprintln(w, "2. 强制重启 Steam 客户端")
	_, _ = fmt.Fprintln(w, "3. 清理系统 DNS 缓存")
	_, _ = fmt.Fprintln(w, "4. 选项")
	_, _ = fmt.Fprintln(w, "0. 退出")
	if admin {
		_, _ = fmt.Fprintln(w, "当前权限: 管理员")
	}
}
