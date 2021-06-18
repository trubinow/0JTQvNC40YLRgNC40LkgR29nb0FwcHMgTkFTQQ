package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"nasa/helpers"
	"nasa/interfaces"
	"nasa/parsers"
	"net/http"
	"strconv"
	"time"
)

func main() {
	logger := logrus.New()
	contextLogger := logger.WithFields(logrus.Fields{
		"cmd": "url-collector",
	})

	apiKey := helpers.GetEnv("API_KEY", "DEMO_KEY")
	apiUrlTemplate := helpers.GetEnv("API_URL", "https://api.nasa.gov/planetary/apod?api_key=%s&date=%s")
	concurrentRequests := helpers.GetEnv("CONCURRENT_REQUESTS", "5")
	port := helpers.GetEnv("PORT", "8080")

	var concurrentLimit int64
	if v, err := strconv.ParseInt(concurrentRequests, 10, 32); err == nil && v > 0 {
		concurrentLimit = v
	} else {
		logger.Fatalf("CONCURRENT_REQUESTS conversion to int error(%v) or equal to 0", err)
	}

	concurrent := make(chan int, concurrentLimit)
	var parser interfaces.Parser
	var httpClient = &http.Client{}
	parser = parsers.NewNasaParser(apiKey, apiUrlTemplate, httpClient)

	handler := NewHandler(contextLogger, &concurrent, parser)
	r := mux.NewRouter()
	r.HandleFunc("/status", handler.Status)
	r.HandleFunc("/pictures", handler.Pictures)

	srv := &http.Server{
		Handler:      r,
		Addr:         fmt.Sprintf("0.0.0.0:%s", port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	logger.Infof("Server is listening port: %s", port)
	logger.Warn(srv.ListenAndServe())

	return
}
