package models

type Response struct {
	URLS []string `json:"urls"`
	Error string `json:"error,omitempty"`
}

type ErrResponse struct {
	Error string `json:"error"`
}
