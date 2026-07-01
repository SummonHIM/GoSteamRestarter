package main

import (
	"io"
	"path/filepath"
	"testing"

	"gosteamrestarter/internal/core"
)

type bootstrapPlatform struct {
	findCalls int
	findPath  string
	startPath string
	startArgs string
}

func (p *bootstrapPlatform) FindSteamPath() (string, error) {
	p.findCalls++
	return p.findPath, nil
}

func (p *bootstrapPlatform) KillSteam() error {
	return nil
}

func (p *bootstrapPlatform) StartSteam(path, args string) error {
	p.startPath = path
	p.startArgs = args
	return nil
}

func (p *bootstrapPlatform) FlushDNS() error {
	return nil
}

func TestRunDesktopPassesLoadedConfigToRunner(t *testing.T) {
	store := core.NewConfigStore(t.TempDir())
	want := core.Config{SteamPath: "/saved/Steam.exe", SteamArgs: "-silent"}
	if err := store.Save(want); err != nil {
		t.Fatalf("seed config: %v", err)
	}

	platform := &bootstrapPlatform{findPath: "/discovered/Steam.exe"}
	called := false
	code := runDesktop(io.Discard, core.Services{Platform: platform, ConfigStore: store}, func(app *core.App) error {
		called = true
		return app.StartSteam()
	})

	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	if !called {
		t.Fatal("expected desktop runner to be called")
	}
	if platform.findCalls != 0 {
		t.Fatalf("expected no discovery lookup, got %d", platform.findCalls)
	}
	if platform.startPath != want.SteamPath || platform.startArgs != want.SteamArgs {
		t.Fatalf("expected start args %q %q, got %q %q", want.SteamPath, want.SteamArgs, platform.startPath, platform.startArgs)
	}
}

func TestRunDesktopFindsSteamPathWhenConfigMissing(t *testing.T) {
	store := core.NewConfigStore(t.TempDir())
	wantPath := filepath.Join(t.TempDir(), "Steam.exe")
	platform := &bootstrapPlatform{findPath: wantPath}
	called := false

	code := runDesktop(io.Discard, core.Services{Platform: platform, ConfigStore: store}, func(app *core.App) error {
		called = true
		return app.StartSteam()
	})

	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	if !called {
		t.Fatal("expected desktop runner to be called")
	}
	if platform.findCalls != 1 {
		t.Fatalf("expected one discovery lookup, got %d", platform.findCalls)
	}

	cfg, err := store.Load()
	if err != nil {
		t.Fatalf("load config: %v", err)
	}
	if cfg.SteamPath != wantPath {
		t.Fatalf("expected persisted path %q, got %q", wantPath, cfg.SteamPath)
	}
	if platform.startPath != wantPath {
		t.Fatalf("expected start path %q, got %q", wantPath, platform.startPath)
	}
}
