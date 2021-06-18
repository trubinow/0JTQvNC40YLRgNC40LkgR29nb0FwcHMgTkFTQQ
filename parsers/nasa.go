package parsers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"nasa/interfaces"
	"nasa/models"
)

type NasaParser struct {
	apiKey string
	apiURL string
	httpClient  interfaces.HTTPClient
}

var OverRateLimit = errors.New("over rate limit exceeded")

func NewNasaParser(apiKey string, apiTemplateURL string, httpClient interfaces.HTTPClient) NasaParser {

	return NasaParser{
		apiKey: apiKey,
		apiURL: apiTemplateURL,
		httpClient:  httpClient,
	}
}

func(p NasaParser) Parse(targetDate string) (string, error) {
	resp, err := p.httpClient.Get(fmt.Sprintf(p.apiURL, p.apiKey, targetDate))
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 429 {
		return "", OverRateLimit
	}

	var picture models.NasaPicture

	err = json.Unmarshal(body, &picture)
	if err != nil {
		return "", err
	}

	return picture.URL, nil
}
