package platform

import (
	"runtime"

	"gosteamrestarter/internal/platform/darwin"
	"gosteamrestarter/internal/platform/linux"
	"gosteamrestarter/internal/platform/windows"
)

type Interface interface {
	FindSteamPath() (string, error)
	KillSteam() error
	StartSteam(path, args string) error
	FlushDNS() error
}

func New() Interface {
	switch runtime.GOOS {
	case "windows":
		return windows.New()
	case "darwin":
		return darwin.New()
	case "linux":
		return linux.New()
	default:
		return linux.New()
	}
}
