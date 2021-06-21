package mocks

import "net/http"

//MockClient imitates httpClient
type MockClient struct {
	GetFunc func(url string) (*http.Response, error)
}

//Get imitates httpClient.Get
func (m *MockClient) Get(url string) (*http.Response, error) {
	return m.GetFunc(url)
}
