package mocks

import "net/http"

type MockClient struct {
	GetFunc func(url string) (*http.Response, error)
}

func (m *MockClient) Get(url string) (*http.Response, error) {
	return m.GetFunc(url)
}
