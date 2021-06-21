package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"nasa/handlers"
	"nasa/middlewares"
	"nasa/parsers"
	"nasa/services"
	"nasa/utils"
	"net/http"
	"strconv"
	"time"
)

func main() {
	logger := logrus.New()
	contextLogger := logger.WithFields(logrus.Fields{
		"cmd": "url-collector",
	})

	apiKey := utils.GetEnv("API_KEY", "DEMO_KEY")
	apiUrlTemplate := utils.GetEnv("API_URL", "https://api.nasa.gov/planetary/apod?api_key=%s&date=%s")
	concurrentRequests := utils.GetEnv("CONCURRENT_REQUESTS", "5")
	port := utils.GetEnv("PORT", "8080")

	var concurrentLimit int64
	if v, err := strconv.ParseInt(concurrentRequests, 10, 32); err == nil && v > 0 {
		concurrentLimit = v
	} else {
		logger.Fatalf("CONCURRENT_REQUESTS conversion to int error(%v) or equal to 0", err)
	}

	var httpClient = &http.Client{}
	concurrent := make(chan struct{}, concurrentLimit)
	parser := parsers.NewNasaParser(contextLogger, apiKey, apiUrlTemplate, httpClient)
	runner := services.NewRunner(contextLogger, concurrent, parser)
	handler := handlers.NewHandler(contextLogger, runner)
	mw := middlewares.NewMiddleware(contextLogger)

	r := mux.NewRouter()
	r.Use(mw.ValidateParameters())
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
}
