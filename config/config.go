package config

import (
	"github.com/spf13/viper"
)

// AppConfig will contain example-app configurable information
type AppConfig struct {
	Port        string
	MetricsPort string

	// Emoji endpoint
	GithubEmojiURL string
}

func init() {
	viper.SetEnvPrefix("EXPAPP")
	viper.AutomaticEnv()
	viper.SetDefault("port", "8080")
	viper.SetDefault("metrics_port", "8081")

	// Emoji endpoint
	viper.SetDefault("github_emoji_url", "https://api.github.com/emojis")
}

// GetConfig will generate the standard AppConfig
func GetConfig() AppConfig {
	return AppConfig{
		Port:        viper.GetString("port"),
		MetricsPort: viper.GetString("metrics_port"),

		// Emoji endpoint
		GithubEmojiURL: viper.GetString("github_emoji_url"),
	}
}
