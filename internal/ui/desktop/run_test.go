package desktop

import (
	"testing"

	"gosteamrestarter/internal/core"
)

func TestRunAcceptsNilApp(t *testing.T) {
	if err := Run(nil); err != nil {
		t.Fatal(err)
	}
}

func TestNewWindowWiresCallbacksToApp(t *testing.T) {
	store := core.NewConfigStore(t.TempDir())
	platform := &testPlatform{}
	app := core.NewApp(core.Services{Platform: platform, ConfigStore: store})
	if err := app.SaveConfig(core.Config{SteamPath: "/steam/Steam.exe", SteamArgs: "-silent"}); err != nil {
		t.Fatalf("save config: %v", err)
	}

	window := NewWindow(app)

	if err := window.Callbacks.Restart(); err != nil {
		t.Fatalf("restart callback: %v", err)
	}
	if err := window.Callbacks.Kill(); err != nil {
		t.Fatalf("kill callback: %v", err)
	}
	if err := window.Callbacks.FlushDNS(); err != nil {
		t.Fatalf("flush dns callback: %v", err)
	}

	want := core.Config{SteamPath: "/new/Steam.exe", SteamArgs: "-bigpicture"}
	if err := window.Callbacks.SaveSettings(want); err != nil {
		t.Fatalf("save settings callback: %v", err)
	}

	if platform.killCalls != 2 {
		t.Fatalf("expected 2 kill calls, got %d", platform.killCalls)
	}
	if platform.startCalls != 1 {
		t.Fatalf("expected 1 start call, got %d", platform.startCalls)
	}
	if platform.startPath != "/steam/Steam.exe" || platform.startArgs != "-silent" {
		t.Fatalf("expected restart to use saved config, got %q %q", platform.startPath, platform.startArgs)
	}
	if platform.flushCalls != 1 {
		t.Fatalf("expected 1 flush call, got %d", platform.flushCalls)
	}

	got, err := store.Load()
	if err != nil {
		t.Fatalf("load config: %v", err)
	}
	if got != want {
		t.Fatalf("expected saved config %+v, got %+v", want, got)
	}
}

type testPlatform struct {
	killCalls  int
	startCalls int
	flushCalls int
	startPath  string
	startArgs  string
}

func (p *testPlatform) FindSteamPath() (string, error) { return "", nil }

func (p *testPlatform) KillSteam() error {
	p.killCalls++
	return nil
}

func (p *testPlatform) StartSteam(path, args string) error {
	p.startCalls++
	p.startPath = path
	p.startArgs = args
	return nil
}

func (p *testPlatform) FlushDNS() error {
	p.flushCalls++
	return nil
}
