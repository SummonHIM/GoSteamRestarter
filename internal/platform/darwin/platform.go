package darwin

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Platform struct{}

func New() Platform {
	return Platform{}
}

func (Platform) DefaultSteamPath() string {
	return "/Applications/Steam.app"
}

func (p Platform) FindSteamPath() (string, error) {
	candidates := []string{
		"/Applications/Steam.app",
	}
	if home, err := os.UserHomeDir(); err == nil {
		candidates = append(candidates, home+"/Applications/Steam.app")
	}
	for _, c := range candidates {
		if _, err := os.Stat(c); err == nil {
			return c, nil
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
		allArgs := append([]string{"-a", path, "--args"}, fields...)
		cmd = exec.Command("open", allArgs...)
	} else {
		cmd = exec.Command("open", "-a", path)
	}
	return cmd.Start()
}

func (Platform) FlushDNS() error {
	if out, err := exec.Command("dscacheutil", "-flushcache").CombinedOutput(); err != nil {
		return fmt.Errorf("flush dns cache: %s", strings.TrimSpace(string(out)))
	}
	_ = exec.Command("killall", "-HUP", "mDNSResponder").Run()
	return nil
}
