package linux

import "testing"

func TestDefaultSteamPath(t *testing.T) {
	p := New()
	if p.DefaultSteamPath() == "" {
		t.Fatal("expected a default path")
	}
}
