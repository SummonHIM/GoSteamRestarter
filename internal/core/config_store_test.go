package core

import (
	"os"
	"path/filepath"
	"testing"
)

func TestConfigStoreRoundTrip(t *testing.T) {
	dir := t.TempDir()
	store := NewConfigStore(dir)
	want := Config{SteamPath: filepath.Join(dir, "Steam.exe"), SteamArgs: "-bigpicture"}
	if err := store.Save(want); err != nil {
		t.Fatalf("save: %v", err)
	}
	got, err := store.Load()
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if got != want {
		t.Fatalf("got %+v, want %+v", got, want)
	}
	if _, err := os.Stat(filepath.Join(dir, "config.json")); err != nil {
		t.Fatalf("config file missing: %v", err)
	}
}

func TestConfigDirUsesUserConfigDir(t *testing.T) {
	oldUserConfigDir := userConfigDir
	t.Cleanup(func() { userConfigDir = oldUserConfigDir })

	base := filepath.Join(t.TempDir(), "profile")
	userConfigDir = func() (string, error) {
		return base, nil
	}

	got, err := ConfigDir()
	if err != nil {
		t.Fatalf("config dir: %v", err)
	}

	want := filepath.Join(base, appName)
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestAppSaveAndLoadConfigUseDefaultConfigDir(t *testing.T) {
	oldUserConfigDir := userConfigDir
	t.Cleanup(func() { userConfigDir = oldUserConfigDir })

	base := filepath.Join(t.TempDir(), "profile")
	userConfigDir = func() (string, error) {
		return base, nil
	}

	app := App{}
	want := Config{
		SteamPath: filepath.Join(base, "Steam.exe"),
		SteamArgs: "-silent",
	}

	if err := app.SaveConfig(want); err != nil {
		t.Fatalf("save config: %v", err)
	}

	got, err := app.LoadConfig()
	if err != nil {
		t.Fatalf("load config: %v", err)
	}

	if got != want {
		t.Fatalf("got %+v, want %+v", got, want)
	}

	configPath := filepath.Join(base, appName, configFileName)
	if _, err := os.Stat(configPath); err != nil {
		t.Fatalf("config file missing at default path: %v", err)
	}
}
