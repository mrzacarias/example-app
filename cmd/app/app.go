package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/mrzacarias/example-app/config"
	emoji "github.com/mrzacarias/example-app/internal/emoji"
	"github.com/mrzacarias/example-app/internal/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

// EmojiClient is a (mockable) Emoji client
var EmojiClient emoji.Contract

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetReportCaller(true)
	EmojiClient = emoji.NewClient()
}

// RootHandler is example-app root friendly page
func RootHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "<html><head><title>example-app</title></head><body><p align=\"center\"><h1 align=\"center\">example-app was created by Bootstrapper!</h1></p></body></html>")
}

// HealthCheckHandler is example-app endpoint for livenessProbe
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// EmojiHandler is example-app endpoint for the internal package Emoji
func EmojiHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	metrics.RequestsTotal.With(prometheus.Labels{"endpoint": "emoji"}).Inc()

	// Total request time start
	requestStart := time.Now()

	// Get the name parameter
	nameParams, ok := r.URL.Query()["name"]
	if !ok || len(nameParams[0]) < 1 {
		metrics.RequestsErrors.With(prometheus.Labels{"endpoint": "emoji", "type": "r.URL.Query"}).Inc()
		log.WithFields(log.Fields{"endpoint": "emoji"}).Errorln("Url Param 'name' is missing")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.WithFields(log.Fields{"endpoint": "emoji", "request": nameParams}).Infoln("Request started")

	// Make emoji.GetFromGithub request
	emojiRes, err := EmojiClient.GetFromGithub(emoji.Request{Name: nameParams[0]})
	if err != nil {
		metrics.RequestsErrors.With(prometheus.Labels{"endpoint": "emoji", "type": "emoji.GetFromGithub"}).Inc()
		log.WithFields(log.Fields{"endpoint": "emoji", "error": err}).Errorln("Error on emoji.GetFromGithub")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	metrics.RequestDurationTotal.With(prometheus.Labels{"endpoint": "emoji"}).Set(time.Since(requestStart).Seconds())
	w.WriteHeader(http.StatusOK)

	if emojiRes != nil {
		fmt.Fprintf(w, "<html><head><title>example-app</title></head><body><p align=\"center\"><h1 align=\"center\">Emoji Found: <img src=\"%s\"/></h1></p></body></html>", emojiRes.EmojiURL)
	} else {
		fmt.Fprintf(w, "<html><head><title>example-app</title></head><body><p align=\"center\"><h1 align=\"center\">Emoji Not Found :/</h1></p></body></html>")
	}
}

// serveHTTP will start the HTTP server and it's endpoints
func serveHTTP(srv *http.Server, errChan chan error) {
	http.HandleFunc("/", RootHandler)
	http.HandleFunc("/healthcheck", HealthCheckHandler)
	http.HandleFunc("/healthz", HealthCheckHandler)
	http.HandleFunc("/emoji", EmojiHandler)

	log.WithFields(log.Fields{"address": srv.Addr}).Info("`example-app` listening")
	err := srv.ListenAndServe()
	if err != nil {
		errChan <- fmt.Errorf("cannot listen to address: %v", err)
		return
	}
}

// serveMetrics will start the HTTP server for metrics, that will be consumed by telegraf
func serveMetrics(srv *http.Server, errChan chan error) {
	mux := http.NewServeMux()
	srv.Handler = mux
	mux.Handle("/metrics", promhttp.Handler())

	log.WithFields(log.Fields{"address": srv.Addr}).Info("`example-app` Metrics listening")
	err := srv.ListenAndServe()
	if err != nil {
		errChan <- fmt.Errorf("cannot listen to address: %v", err)
		return
	}
}

// Main thread
func main() {
	log.Info("example-app Initialized!")
	cfg := config.GetConfig()

	// Preparing channels to listen to hard errors and signals
	errChan := make(chan error)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan)

	// HTTP logic
	proxySrv := &http.Server{Addr: fmt.Sprintf(":%s", cfg.Port)}
	metricsSrv := &http.Server{Addr: fmt.Sprintf(":%s", cfg.MetricsPort)}
	go serveHTTP(proxySrv, errChan)
	go serveMetrics(metricsSrv, errChan)

	// Blocking the main thread execution while waits for a hard error or signal
	select {
	case err := <-errChan:
		log.WithFields(log.Fields{"error": err}).Error("Failed to start server")
	case sig := <-sigChan:
		if sigMsg := sig.String(); sigMsg == "interrupt" || sigMsg == "terminated" {
			log.WithFields(log.Fields{"signal_message": sigMsg}).Error("Received termination signal, shutting down servers...")
		}
	}
}
