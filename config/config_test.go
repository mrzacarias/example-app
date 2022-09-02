package config_test

import (
	"testing"

	cfg "github.com/mrzacarias/example-app/config"
)

func TestMain(t *testing.T) {
	config := cfg.GetConfig()

	checkConfig(t, "Port", config.Port, "8080")
	checkConfig(t, "MetricsPort", config.MetricsPort, "8081")

	checkConfig(t, "GithubEmojiURL", config.GithubEmojiURL, "https://api.github.com/emojis")
}

// DRY helpers for checking a config attribute
func checkConfig(t *testing.T, attr string, got string, want string) {
	if want != got {
		t.Fatalf("Attribute '%s' should be %s, but it was %s", attr, want, got)
	}
}
