package emoji_test

import (
	"testing"

	emoji "github.com/mrzacarias/example-app/internal/emoji"
	"github.com/mrzacarias/example-app/internal/mock"
	"github.com/spf13/viper"
)

var EmojiClient emoji.Contract

func init() {
	EmojiClient = emoji.NewClient()
}

// Helper for setup/teardown
func setup(t *testing.T) func(t *testing.T) {
	server := mock.GithubEmojiServerStub()
	viper.Set("github_emoji_url", server.URL)

	// Returning teardown function
	return func(t *testing.T) {
		server.Close()
	}
}

func TestGetFromGithub(test *testing.T) {
	teardown := setup(test)
	defer teardown(test)

	test.Run("Emoji requested exists", func(t *testing.T) {
		got, err := EmojiClient.GetFromGithub(emoji.Request{Name: "100"})
		if err != nil {
			t.Fatal("Error on GetFromGithub: ", err)
		}

		if got == nil || got.EmojiURL != mock.FoundEmojiURL {
			t.Fatalf("Should be %v, but it was %v", mock.FoundEmojiURL, got)
		}
	})

	test.Run("Emoji requested do not exist", func(t *testing.T) {
		got, err := EmojiClient.GetFromGithub(emoji.Request{Name: "foo"})
		if err != nil {
			t.Fatal("Error on GetFromGithub: ", err)
		}

		if got != nil {
			t.Fatalf("Should be nil, but it was %v", got)
		}
	})
}
