package config

import (
	"os"
	"testing"
)

func TestLoadUsesEnv(t *testing.T) {
	os.Setenv("API_PORT", "9999")
	defer os.Unsetenv("API_PORT")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}
	if cfg.ApiPort != "9999" {
		t.Fatalf("expected ApiPort=9999 got %s", cfg.ApiPort)
	}
}
