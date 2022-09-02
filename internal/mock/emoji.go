package mock

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/mrzacarias/example-app/internal/emoji"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetReportCaller(true)
}

// --------------- EMOJI MOCK ---------------

// FoundEmojiURL will mock a successful emoji URL
const FoundEmojiURL = "https://github.githubassets.com/images/icons/emoji/unicode/1f4af.png?v8"

// EmojiClient will mock when an emoji was found
type EmojiClient struct{}

// GetFromGithub will mock EmojiClient GetFromGithub
func (ec *EmojiClient) GetFromGithub(req emoji.Request) (*emoji.Response, error) {
	if req.Name == "100" {
		log.WithFields(log.Fields{"mock": true, "package": "emoji"}).Info("GetFromGithub found")
		return &emoji.Response{EmojiURL: FoundEmojiURL}, nil
	}
	log.WithFields(log.Fields{"mock": true, "package": "emoji"}).Info("GetFromGithub not found")
	return nil, nil
}

// --------------- GITHUB EMOJI SERVER MOCK ---------------

// GithubEmojiResponse will be the Emoji stubbed server result
var GithubEmojiResponse = fmt.Sprintf(`{"100": "%s"}`, FoundEmojiURL)

// GithubEmojiServerStub will stub github emoji API requests and return a valid response
func GithubEmojiServerStub() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.WithFields(log.Fields{"mock": true, "package": "emoji"}).Info("GithubEmojiServerStub")
		w.Write([]byte(GithubEmojiResponse))
	}))
}
