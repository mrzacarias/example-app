package emoji

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/mrzacarias/example-app/config"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetReportCaller(true)
}

// Contract is the interface will define the service methods
type Contract interface {
	GetFromGithub(Request) (*Response, error)
}

// Client will be our base struct
type Client struct{}

// NewClient will return an emoji client pointer
func NewClient() *Client {
	return &Client{}
}

// Request will formalize the endpoint request
type Request struct {
	Name string `json:"name"`
}

// Response will formalize the endpoint response
type Response struct {
	EmojiURL string `json:"image_url"`
}

// GetFromGithub will reach the github emoji API and return a list of emojis
func (ec *Client) GetFromGithub(req Request) (*Response, error) {
	cfg := config.GetConfig()
	log.WithFields(log.Fields{
		"package": "emoji",
		"request": req,
		"apiURL":  cfg.GithubEmojiURL,
	}).Infoln("Requesting emoji from GitHub")

	// Building request to github
	githubReq, err := http.NewRequest("GET", cfg.GithubEmojiURL, strings.NewReader(""))
	if err != nil {
		return nil, err
	}

	// Making request to github
	client := &http.Client{}
	resp, err := client.Do(githubReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Parsing body and storing on list of emojis
	var emojisList map[string]string
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &emojisList)
	if err != nil {
		return nil, err
	}

	// Return requested emoji URL, if existent
	if url, found := emojisList[req.Name]; found {
		return &Response{EmojiURL: url}, nil
	}

	return nil, nil
}
