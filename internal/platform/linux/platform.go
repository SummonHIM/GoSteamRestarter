package linux

import (
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
	return "/usr/bin/steam"
}

func (p Platform) FindSteamPath() (string, error) {
	candidates := []string{
		"/usr/bin/steam",
		"/usr/games/steam",
	}
	if home, err := os.UserHomeDir(); err == nil {
		candidates = append(candidates,
			filepath.Join(home, ".local", "share", "Steam", "steam.sh"),
			filepath.Join(home, ".steam", "steam.sh"),
		)
	}
	for _, c := range candidates {
		if _, err := os.Stat(c); err == nil {
			return c, nil
		}
	}
	if out, err := exec.Command("which", "steam").Output(); err == nil {
		p := strings.TrimSpace(string(out))
		if p != "" {
			return p, nil
		}
	}
	return "", fmt.Errorf("steam not found")
}

func (Platform) KillSteam() error {
	cmd := exec.Command("pkill", "-f", "steam")
	output, err := cmd.CombinedOutput()
	if err != nil {
		if cmd.ProcessState != nil && cmd.ProcessState.ExitCode() == 1 {
			return nil
		}
		return fmt.Errorf("kill steam: %s", strings.TrimSpace(string(output)))
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
	cmd.Stdout = nil
	cmd.Stderr = nil
	return cmd.Start()
}

func (Platform) FlushDNS() error {
	if _, err := exec.Command("systemd-resolve", "--flush-caches").CombinedOutput(); err == nil {
		return nil
	}
	if _, err := exec.Command("resolvectl", "flush-caches").CombinedOutput(); err == nil {
		return nil
	}
	if _, err := exec.Command("nscd", "-i", "hosts").CombinedOutput(); err == nil {
		return nil
	}
	return fmt.Errorf("flush dns: no supported resolver found")
}
