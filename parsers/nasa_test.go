package parsers

import (
	"bytes"
	"io/ioutil"
	"nasa/mocks"
	"net/http"
	"testing"
)

func TestNasaParser(t *testing.T) {
	json := `{"copyright":"Amir H. Abolfath","date":"2019-12-06","explanation":"This frame.","hdurl":"https://apod.nasa.gov/apod/image/1912/TaurusAbolfath.jpg","media_type":"image","service_version":"v1","title":"Pleiades to Hyades","url":"https://apod.nasa.gov/apod/image/1912/TaurusAbolfath1024.jpg"}`
	r := ioutil.NopCloser(bytes.NewReader([]byte(json)))

	client := &mocks.MockClient{}
	client.GetFunc = func(url string) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body: r,
		}, nil
	}

	parser := NewNasaParser("", "", client)
	url, _ := parser.Parse("2012-06-11")

	if  url != "https://apod.nasa.gov/apod/image/1912/TaurusAbolfath1024.jpg" {
		t.Errorf("url: %s", url)
	}
}
