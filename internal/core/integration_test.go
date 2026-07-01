package core

import "testing"

func TestDefaultConfigIsSerializable(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.SteamPath != "" || cfg.SteamArgs != "" {
		t.Fatal("unexpected defaults")
	}
}
