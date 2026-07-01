package core

import (
	"path/filepath"
	"testing"
)

type fakePlatform struct {
	find  func() (string, error)
	kill  func() error
	start func(path, args string) error
	flush func() error
}

func (p fakePlatform) FindSteamPath() (string, error) {
	if p.find != nil {
		return p.find()
	}
	return "", nil
}

func (p fakePlatform) KillSteam() error {
	if p.kill != nil {
		return p.kill()
	}
	return nil
}

func (p fakePlatform) StartSteam(path, args string) error {
	if p.start != nil {
		return p.start(path, args)
	}
	return nil
}

func (p fakePlatform) FlushDNS() error {
	if p.flush != nil {
		return p.flush()
	}
	return nil
}

func TestRestartSteamUsesSavedPathAndArgs(t *testing.T) {
	called := 0
	app := NewApp(Services{Platform: fakePlatform{start: func(path, args string) error {
		called++
		if path != "/steam/Steam.exe" || args != "-bigpicture" {
			t.Fatalf("unexpected start args: %q %q", path, args)
		}
		return nil
	}}})
	app.cfg = Config{SteamPath: "/steam/Steam.exe", SteamArgs: "-bigpicture"}
	if err := app.RestartSteam(); err != nil {
		t.Fatalf("restart: %v", err)
	}
	if called != 1 {
		t.Fatalf("start called %d times", called)
	}
}

func TestStartSteamReturnsSteamNotFoundWhenPathMissing(t *testing.T) {
	app := NewApp(Services{Platform: fakePlatform{}})

	err := app.StartSteam()
	if err != ErrSteamNotFound {
		t.Fatalf("expected ErrSteamNotFound, got %v", err)
	}
}

func TestFindOrConfirmSteamPathReturnsSavedPathWithoutLookup(t *testing.T) {
	lookups := 0
	app := NewApp(Services{Platform: fakePlatform{find: func() (string, error) {
		lookups++
		return "/discovered/Steam.exe", nil
	}}})
	app.cfg = Config{SteamPath: "/saved/Steam.exe", SteamArgs: "-silent"}

	path, err := app.FindOrConfirmSteamPath()
	if err != nil {
		t.Fatalf("find or confirm: %v", err)
	}
	if path != "/saved/Steam.exe" {
		t.Fatalf("expected saved path, got %q", path)
	}
	if lookups != 0 {
		t.Fatalf("expected no platform lookup, got %d", lookups)
	}
}

func TestFindOrConfirmSteamPathPersistsDiscoveredPath(t *testing.T) {
	dir := t.TempDir()
	store := NewConfigStore(dir)
	app := NewApp(Services{
		ConfigStore: store,
		Platform: fakePlatform{find: func() (string, error) {
			return filepath.Join(dir, "Steam.exe"), nil
		}},
	})
	app.cfg = Config{SteamArgs: "-bigpicture"}

	path, err := app.FindOrConfirmSteamPath()
	if err != nil {
		t.Fatalf("find or confirm: %v", err)
	}
	want := filepath.Join(dir, "Steam.exe")
	if path != want {
		t.Fatalf("expected discovered path %q, got %q", want, path)
	}
	if app.cfg.SteamPath != want {
		t.Fatalf("expected app config to cache %q, got %q", want, app.cfg.SteamPath)
	}

	persisted, err := store.Load()
	if err != nil {
		t.Fatalf("load persisted config: %v", err)
	}
	if persisted.SteamPath != want {
		t.Fatalf("expected persisted path %q, got %q", want, persisted.SteamPath)
	}
	if persisted.SteamArgs != "-bigpicture" {
		t.Fatalf("expected persisted args to be preserved, got %q", persisted.SteamArgs)
	}
}

func TestSaveConfigCachesPathAndArgsForStartSteam(t *testing.T) {
	started := 0
	app := NewApp(Services{
		ConfigStore: NewConfigStore(t.TempDir()),
		Platform: fakePlatform{start: func(path, args string) error {
			started++
			if path != "/steam/Steam.exe" || args != "-silent" {
				t.Fatalf("unexpected start args: %q %q", path, args)
			}
			return nil
		}},
	})

	cfg := Config{SteamPath: "/steam/Steam.exe", SteamArgs: "-silent"}
	if err := app.SaveConfig(cfg); err != nil {
		t.Fatalf("save config: %v", err)
	}
	if err := app.StartSteam(); err != nil {
		t.Fatalf("start steam: %v", err)
	}
	if started != 1 {
		t.Fatalf("start called %d times", started)
	}
}

func TestLoadConfigCachesPathAndArgs(t *testing.T) {
	store := NewConfigStore(t.TempDir())
	want := Config{SteamPath: "/steam/Steam.exe", SteamArgs: "-tenfoot"}
	if err := store.Save(want); err != nil {
		t.Fatalf("seed config: %v", err)
	}

	started := 0
	app := NewApp(Services{
		ConfigStore: store,
		Platform: fakePlatform{start: func(path, args string) error {
			started++
			if path != want.SteamPath || args != want.SteamArgs {
				t.Fatalf("unexpected start args: %q %q", path, args)
			}
			return nil
		}},
	})

	cfg, err := app.LoadConfig()
	if err != nil {
		t.Fatalf("load config: %v", err)
	}
	if cfg != want {
		t.Fatalf("expected loaded config %+v, got %+v", want, cfg)
	}
	if err := app.StartSteam(); err != nil {
		t.Fatalf("start steam: %v", err)
	}
	if started != 1 {
		t.Fatalf("start called %d times", started)
	}
}

func TestFlushDNSDelegatesToPlatform(t *testing.T) {
	called := 0
	app := NewApp(Services{Platform: fakePlatform{flush: func() error {
		called++
		return nil
	}}})

	if err := app.FlushDNS(); err != nil {
		t.Fatalf("flush dns: %v", err)
	}
	if called != 1 {
		t.Fatalf("flush called %d times", called)
	}
}
