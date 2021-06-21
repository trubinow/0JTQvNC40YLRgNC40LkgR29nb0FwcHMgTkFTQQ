package handlers

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"nasa/middlewares"
	"nasa/parsers"
	"nasa/utils"
	"net/http"
)

//Response wraps http response
type Response struct {
	URLS  []string `json:"urls,omitempty"`
	Error string   `json:"error,omitempty"`
}

type ErrResponse struct {
	Error string `json:"error"`
}

//Runner interface
type Runner interface {
	Run(dates []string) ([]string, error)
	Blocked() bool
	Error() error
}

//Handler is handler for http requests
type Handler struct {
	logger *logrus.Entry
	runner Runner
}

//NewHandler creates new http handler
func NewHandler(logger *logrus.Entry, runner Runner) *Handler {
	return &Handler{
		logger: logger,
		runner: runner,
	}
}

//Status returns status of the service
func (h *Handler) Status(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("OK"))
	if err != nil {
		h.logger.WithError(err).Warnf("status response write error")
	}
}

//Pictures returns pictures result data set
func (h *Handler) Pictures(w http.ResponseWriter, r *http.Request) {

	props, ok := r.Context().Value("props").(middlewares.Request)
	if !ok {
		w.WriteHeader(400)
		_, err := w.Write([]byte("Something wrong with props"))
		if err != nil {
			h.logger.WithError(err).Warnf("pictures response write error")
		}
		return
	}

	dates, err := utils.Interval(props.StartDate, props.EndDate)
	if err != nil {
		w.WriteHeader(400)
		_, err = w.Write([]byte(err.Error()))
		if err != nil {
			h.logger.WithError(err).Warnf("pictures response write error")
		}
		return
	}

	h.logger.Infof("Dates: %v", dates)
	images, err := h.runner.Run(dates)
	var response = Response{URLS: images}

	if h.runner.Blocked() {
		w.WriteHeader(429)
		response.Error = parsers.ErrOverRateLimit.Error()
	} else if err != nil {
		w.WriteHeader(400)
		response.Error = parsers.ErrUnknown.Error()
	}

	res, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(400)
		res, _ = json.Marshal(ErrResponse{Error: err.Error()})
	}

	_, err = w.Write(res)
	if err != nil {
		h.logger.WithError(err).Warnf("pictures response write error")
	}
}
