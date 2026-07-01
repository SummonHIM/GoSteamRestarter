package core

import (
	"errors"
	"testing"
)

type stubPlatform struct{}

func (stubPlatform) FindSteamPath() (string, error) { return "", nil }
func (stubPlatform) KillSteam() error               { return nil }
func (stubPlatform) StartSteam(path, args string) error {
	return nil
}
func (stubPlatform) FlushDNS() error { return nil }

func TestNewAppWiresConfigAndPlatform(t *testing.T) {
	platformStub := stubPlatform{}
	services := Services{
		ConfigStore: NewConfigStore(t.TempDir()),
		Platform:    platformStub,
	}

	app := NewApp(services)
	if app == nil {
		t.Fatal("expected app")
	}
	if app.services.ConfigStore != services.ConfigStore {
		t.Fatal("expected config store to be wired into app services")
	}
	if app.services.Platform == nil {
		t.Fatal("expected platform dependency to be preserved on app services")
	}
}

func TestTypedErrorsAreDefined(t *testing.T) {
	if ErrSteamNotFound == nil {
		t.Fatal("expected ErrSteamNotFound")
	}
	if ErrSteamNotRunning == nil {
		t.Fatal("expected ErrSteamNotRunning")
	}
	if ErrPermissionDenied == nil {
		t.Fatal("expected ErrPermissionDenied")
	}
	if errors.Is(ErrSteamNotFound, ErrSteamNotRunning) {
		t.Fatal("expected distinct typed errors")
	}
}

func TestResultTypeExists(t *testing.T) {
	result := Result{}
	if result != (Result{}) {
		t.Fatal("expected zero-value result")
	}
}
