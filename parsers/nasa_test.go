package parsers

import (
	"bytes"
	"errors"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"nasa/mocks"
	"net/http"
	"testing"
)

var logger *logrus.Logger
var contextLogger *logrus.Entry

func init() {
	logger = logrus.New()
	contextLogger = logger.WithFields(logrus.Fields{
		"cmd": "test-nasa-parser",
	})
}

func TestNasaParser(t *testing.T) {
	json := `{"copyright":"Amir H. Abolfath","date":"2019-12-06","explanation":"This frame.","hdurl":"https://apod.nasa.gov/apod/image/1912/TaurusAbolfath.jpg","media_type":"image","service_version":"v1","title":"Pleiades to Hyades","url":"https://apod.nasa.gov/apod/image/1912/TaurusAbolfath1024.jpg"}`
	r := ioutil.NopCloser(bytes.NewReader([]byte(json)))

	client := &mocks.MockClient{}
	client.GetFunc = func(url string) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}

	parser := NewNasaParser(contextLogger, "", "", client)
	url, _ := parser.Parse("2012-06-11")

	if url != "https://apod.nasa.gov/apod/image/1912/TaurusAbolfath1024.jpg" {
		t.Errorf("url: %s", url)
	}
}

func TestNasaParserWrongInterval(t *testing.T) {
	json := `{code: 400, msg: "Date must be between Jun 16, 1995 and Jun 20, 2021.",service_version: "v1"}`
	r := ioutil.NopCloser(bytes.NewReader([]byte(json)))

	client := &mocks.MockClient{}
	client.GetFunc = func(url string) (*http.Response, error) {
		return &http.Response{
			StatusCode: 400,
			Body:       r,
		}, nil
	}

	inputDate := "1995-05-11"
	parser := NewNasaParser(contextLogger, "", "", client)
	_, err := parser.Parse(inputDate)

	if !errors.Is(err, ErrWrongDateInterval) {
		t.Error("Wrong date interval error expected")
	}
}

func TestNasaParserOverRateLimit(t *testing.T) {
	json := `<html>
  <body>
    <h1>OVER_RATE_LIMIT</h1>
    <p>You have exceeded your rate limit. Try again later or contact us at https:&#x2F;&#x2F;api.nasa.gov:443&#x2F;contact&#x2F; for assistance</p>
  </body>
</html>`
	r := ioutil.NopCloser(bytes.NewReader([]byte(json)))

	client := &mocks.MockClient{}
	client.GetFunc = func(url string) (*http.Response, error) {
		return &http.Response{
			StatusCode: 429,
			Body:       r,
		}, nil
	}

	inputDate := "2012-05-11"
	parser := NewNasaParser(contextLogger, "", "", client)
	_, err := parser.Parse(inputDate)

	if !errors.Is(err, ErrOverRateLimit) {
		t.Error("Over rate limit error expected")
	}
}
