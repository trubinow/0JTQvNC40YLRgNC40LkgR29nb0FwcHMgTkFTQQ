package interfaces

import "net/http"

// HTTPClient interface
type HTTPClient interface {
	Get(url string) (*http.Response, error)
}
