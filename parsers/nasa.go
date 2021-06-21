package parsers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

// HTTPClient interface
type HTTPClient interface {
	Get(url string) (*http.Response, error)
}

type NasaParser struct {
	logger *logrus.Entry
	apiKey string
	apiURL string
	httpClient  HTTPClient
}

var (
	ErrOverRateLimit = errors.New("over rate limit exceeded")
	ErrWrongDateInterval = errors.New("wrong date interval")
)

type NasaPicture struct {
	Copyright      string `json:"contacts"`
	Date           string `json:"date"`
	Explanation    string `json:"explanation"`
	HdURL          string `json:"hdurl"`
	MediaType      string `json:"media_type"`
	ServiceVersion string `json:"service_version"`
	Title          string `json:"title"`
	URL            string `json:"url"`
}

func NewNasaParser(logger *logrus.Entry, apiKey string, apiTemplateURL string, httpClient HTTPClient) *NasaParser {

	return &NasaParser{
		logger: logger,
		apiKey: apiKey,
		apiURL: apiTemplateURL,
		httpClient:  httpClient,
	}
}

func(p *NasaParser) Parse(targetDate string) (string, error) {
	resp, err := p.httpClient.Get(fmt.Sprintf(p.apiURL, p.apiKey, targetDate))
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	defer func() {
		closeErr := resp.Body.Close()
		if closeErr != nil {
			p.logger.WithError(closeErr).Warnf("body close error")
		}
	}()

	if resp.StatusCode == 429 {
		return "", ErrOverRateLimit
	}

	if resp.StatusCode == 400 {
		return "", ErrWrongDateInterval
	}

	var picture NasaPicture
	err = json.Unmarshal(body, &picture)
	if err != nil {
		return "", err
	}

	return picture.URL, nil
}
