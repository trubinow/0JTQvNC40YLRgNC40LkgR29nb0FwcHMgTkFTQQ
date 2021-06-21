package services

import (
	"github.com/sirupsen/logrus"
	"nasa/parsers"
	"net/http"
	"net/http/httptest"
	"testing"
)

var logger *logrus.Logger
var contextLogger *logrus.Entry

func init() {
	logger = logrus.New()
	contextLogger = logger.WithFields(logrus.Fields{
		"cmd": "test-runner",
	})
}

//BenchmarkRunner_Run runner Run method benchmark test
func BenchmarkRunner_Run(b *testing.B) {
	srv := serverMock()
	defer srv.Close()

	client := &http.Client{}
	concurrent := make(chan struct{}, 5)
	parser := parsers.NewNasaParser(contextLogger, "DEMO_KEY", srv.URL+"/planetary/apod?api_key=%s&date=%s", client)
	runner := NewRunner(contextLogger, concurrent, parser)

	for i := 0; i < b.N; i++ {
		_, err := runner.Run([]string{"2019-12-06", "2019-12-07", "2019-12-10", "2019-12-11", "2019-12-12", "2019-12-14"})
		if err != nil {
			b.Fatalf("Error: %v", err)
		}
	}
}

//serverMock http server mock
func serverMock() *httptest.Server {
	handler := http.NewServeMux()
	handler.HandleFunc("/planetary/apod", picturesMock)
	srv := httptest.NewServer(handler)
	return srv
}

//picturesMock pictures http endpoint mock
func picturesMock(w http.ResponseWriter, r *http.Request) {
	json := `{"copyright":"Amir H. Abolfath","date":"2019-12-06","explanation":"This frame.","hdurl":"https://apod.nasa.gov/apod/image/1912/TaurusAbolfath.jpg","media_type":"image","service_version":"v1","title":"Pleiades to Hyades","url":"https://apod.nasa.gov/apod/image/1912/TaurusAbolfath1024.jpg"}`
	w.WriteHeader(200)
	_, _ = w.Write([]byte(json))
}
