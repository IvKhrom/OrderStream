package config

import (
	"os"
	"testing"
)

func TestGetEnv_Default(t *testing.T) {
	key := "TEST_CFG_KEY"
	_ = os.Unsetenv(key)
	if got := getEnv(key, "def"); got != "def" {
		t.Fatalf("ожидали def, получили %q", got)
	}
}


