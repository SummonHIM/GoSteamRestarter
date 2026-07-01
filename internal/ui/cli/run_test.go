package cli

import (
	"bytes"
	"strings"
	"testing"

	"gosteamrestarter/internal/core"
)

func TestRenderMainMenu(t *testing.T) {
	var buf bytes.Buffer
	RenderMainMenu(&buf, false)
	got := buf.String()
	if !strings.Contains(got, "强制结束 Steam 客户端") {
		t.Fatal("missing menu item")
	}
}

func TestRunExecutesCoreActionsFromMenu(t *testing.T) {
	platform := &testPlatform{}
	app := core.NewApp(core.Services{Platform: platform, ConfigStore: core.NewConfigStore(t.TempDir())})
	if err := app.SaveConfig(core.Config{SteamPath: "/steam/Steam.exe", SteamArgs: "-silent"}); err != nil {
		t.Fatalf("save config: %v", err)
	}

	var out bytes.Buffer
	input := strings.NewReader("1\n2\n3\n0\n")

	if err := Run(input, &out, app); err != nil {
		t.Fatalf("run: %v", err)
	}
	if platform.killCalls != 2 {
		t.Fatalf("expected 2 kill calls, got %d", platform.killCalls)
	}
	if platform.startCalls != 1 {
		t.Fatalf("expected 1 start call, got %d", platform.startCalls)
	}
	if platform.flushCalls != 1 {
		t.Fatalf("expected 1 flush call, got %d", platform.flushCalls)
	}
}

func TestRunSettingsSavesConfigFromMenu(t *testing.T) {
	store := core.NewConfigStore(t.TempDir())
	app := core.NewApp(core.Services{Platform: &testPlatform{}, ConfigStore: store})

	var out bytes.Buffer
	input := strings.NewReader("4\n/steam/Steam.exe\n-bigpicture\n0\n")

	if err := Run(input, &out, app); err != nil {
		t.Fatalf("run: %v", err)
	}

	cfg, err := store.Load()
	if err != nil {
		t.Fatalf("load config: %v", err)
	}
	if cfg.SteamPath != "/steam/Steam.exe" {
		t.Fatalf("expected steam path to be saved, got %q", cfg.SteamPath)
	}
	if cfg.SteamArgs != "-bigpicture" {
		t.Fatalf("expected steam args to be saved, got %q", cfg.SteamArgs)
	}
	if !strings.Contains(out.String(), "设置已保存") {
		t.Fatal("expected saved settings message")
	}
}

type testPlatform struct {
	killCalls  int
	startCalls int
	flushCalls int
}

func (p *testPlatform) FindSteamPath() (string, error) { return "", nil }

func (p *testPlatform) KillSteam() error {
	p.killCalls++
	return nil
}

func (p *testPlatform) StartSteam(path, args string) error {
	p.startCalls++
	return nil
}

func (p *testPlatform) FlushDNS() error {
	p.flushCalls++
	return nil
}
