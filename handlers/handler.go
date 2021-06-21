package handlers

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"nasa/middlewares"
	"nasa/models"
	"nasa/parsers"
	"nasa/utils"
	"net/http"
)

type Response struct {
	URLS []string `json:"urls"`
	Error string `json:"error,omitempty"`
}

type ErrResponse struct {
	Error string `json:"error"`
}

type Runner interface {
	Run(dates []string) ([]string, error)
	Blocked() bool
}

//Handler type
type Handler struct {
	logger     *logrus.Entry
	runner     Runner
}

//NewHandler constructor
func NewHandler(logger *logrus.Entry, runner Runner) *Handler {
	return &Handler{
		logger:     logger,
		runner:     runner,
	}
}

//Status function
func (h *Handler) Status(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("OK"))
	if err != nil {
		h.logger.WithError(err).Warnf("status response write error")
	}
}

//Status function
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
	var response = models.Response{URLS: images}
	if err != nil {
		response.Error = err.Error()
	}

	if h.runner.Blocked() {
		response.Error = parsers.ErrOverRateLimit.Error()
	}

	res, err := json.Marshal(response)
	if err != nil {
		res, _ = json.Marshal(models.ErrResponse{Error: err.Error()})
	}

	_, err = w.Write(res)
	if err != nil {
		h.logger.WithError(err).Warnf("pictures response write error")
	}
}
