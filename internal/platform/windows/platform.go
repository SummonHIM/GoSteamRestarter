package windows

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Platform struct{}

func New() Platform {
	return Platform{}
}

func (Platform) DefaultSteamPath() string {
	return `C:\Program Files (x86)\Steam\steam.exe`
}

func (p Platform) FindSteamPath() (string, error) {
	cmd := exec.Command("reg", "query", `HKLM\SOFTWARE\WOW6432Node\Valve\Steam`, "/v", "InstallPath")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err == nil {
		for _, line := range strings.Split(out.String(), "\n") {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "InstallPath") {
				parts := strings.SplitN(line, "REG_SZ", 2)
				if len(parts) == 2 {
					dir := strings.TrimSpace(parts[1])
					exe := filepath.Join(dir, "steam.exe")
					if _, err := os.Stat(exe); err == nil {
						return exe, nil
					}
				}
			}
		}
	}
	def := p.DefaultSteamPath()
	if _, err := os.Stat(def); err == nil {
		return def, nil
	}
	return "", fmt.Errorf("steam not found")
}

func (Platform) KillSteam() error {
	cmd := exec.Command("taskkill", "/F", "/IM", "steam.exe")
	output, err := cmd.CombinedOutput()
	if err != nil {
		out := string(output)
		if strings.Contains(out, "not found") || strings.Contains(out, "没有找到") {
			return nil
		}
		return fmt.Errorf("kill steam: %s", strings.TrimSpace(out))
	}
	return nil
}

func (Platform) StartSteam(path, args string) error {
	var cmd *exec.Cmd
	if args != "" {
		fields := strings.Fields(args)
		cmd = exec.Command(path, fields...)
	} else {
		cmd = exec.Command(path)
	}
	return cmd.Start()
}

func (Platform) FlushDNS() error {
	cmd := exec.Command("ipconfig", "/flushdns")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("flush dns: %s", strings.TrimSpace(string(output)))
	}
	return nil
}
