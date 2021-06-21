package middlewares

import (
	"context"
	"github.com/sirupsen/logrus"
	"net/http"
)

//Middleware is http middleware
type Middleware struct {
	logger    *logrus.Entry
}

//NewMiddleware creates new http middleware
func NewMiddleware(logger *logrus.Entry) *Middleware {
	return &Middleware{
		logger: logger,
	}
}

//Request wraps user query parameters
type Request struct {
	StartDate string `json:"start_date"`
	EndDate string `json:"end_date"`
}

//ValidateParameters validates start_date, end_date query parameters(/pictures endpoint only)
//and stores them to request context.props
func(m *Middleware) ValidateParameters() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "pictures" {
				startList, ok := r.URL.Query()["start_date"]
				if !ok || len(startList[0]) < 1 {
					w.WriteHeader(400)
					_, err := w.Write([]byte("param start is missing"))
					if err != nil {
						m.logger.WithError(err).Warnf("validate parameters response write error")
					}
					return
				}

				endList, ok := r.URL.Query()["end_date"]
				if !ok || len(endList[0]) < 1 {
					w.WriteHeader(400)
					_, err := w.Write([]byte("param end is missing"))
					if err != nil {
						m.logger.WithError(err).Warnf("validate parameters response write error")
					}
					return
				}

				ctx := context.WithValue(r.Context(), "props", Request{
					StartDate: startList[0],
					EndDate:   endList[0],
				})

				next.ServeHTTP(w, r.WithContext(ctx))
			} else {
				next.ServeHTTP(w, r)
			}
		})
	}
}
