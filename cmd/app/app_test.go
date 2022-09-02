package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mrzacarias/example-app/internal/mock"
)

func init() {
	EmojiClient = &mock.EmojiClient{}
}

func TestHealthCheckHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/healthcheck", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Using ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HealthCheckHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestRootHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Using ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(RootHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	bodyBytes, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		t.Fatal("Error when reading RootHandler response body")
	}

	want := "<html><head><title>example-app</title></head><body><p align=\"center\"><h1 align=\"center\">example-app was created by Bootstrapper!</h1></p></body></html>"
	if got := string(bodyBytes); got != want {
		t.Fatalf("Should be %v, but it was %v", want, got)
	}
}

func TestEmojiHandler(test *testing.T) {
	test.Run("Emoji requested exists", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/emoji?name=100", nil)
		if err != nil {
			t.Fatal(err)
		}

		// Using ResponseRecorder to record the response
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(EmojiHandler)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		bodyBytes, err := ioutil.ReadAll(rr.Body)
		if err != nil {
			t.Fatal("Error when reading EmojiHandler response body")
		}

		want := "<html><head><title>example-app</title></head><body><p align=\"center\"><h1 align=\"center\">Emoji Found: <img src=\"https://github.githubassets.com/images/icons/emoji/unicode/1f4af.png?v8\"/></h1></p></body></html>"
		if got := string(bodyBytes); got != want {
			t.Fatalf("Should be %v, but it was %v", want, got)
		}
	})

	test.Run("Emoji requested do not exist", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/emoji?name=foo", nil)
		if err != nil {
			t.Fatal(err)
		}

		// Using ResponseRecorder to record the response
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(EmojiHandler)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		bodyBytes, err := ioutil.ReadAll(rr.Body)
		if err != nil {
			t.Fatal("Error when reading EmojiHandler response body")
		}

		want := "<html><head><title>example-app</title></head><body><p align=\"center\"><h1 align=\"center\">Emoji Not Found :/</h1></p></body></html>"
		if got := string(bodyBytes); got != want {
			t.Fatalf("Should be %v, but it was %v", want, got)
		}
	})
}
