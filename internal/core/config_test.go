package core

import "testing"

func TestConfigDefaults(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.SteamPath != "" {
		t.Fatalf("expected empty SteamPath, got %q", cfg.SteamPath)
	}
	if cfg.SteamArgs != "" {
		t.Fatalf("expected empty SteamArgs, got %q", cfg.SteamArgs)
	}
}
