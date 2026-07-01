package platform

import "testing"

func TestNewReturnsPlatformImplementation(t *testing.T) {
	p := New()
	if p == nil {
		t.Fatal("expected a platform implementation")
	}
}
