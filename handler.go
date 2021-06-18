package main

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"nasa/helpers"
	"nasa/interfaces"
	"nasa/models"
	"nasa/parsers"
	"net/http"
	"sync"
)

//Handler type
type Handler struct {
	logger     *logrus.Entry
	concurrent *chan int
	parser     interfaces.Parser
}

//NewHandler constructor
func NewHandler(logger *logrus.Entry, concurrent *chan int, parser interfaces.Parser) *Handler {
	return &Handler{
		logger:     logger,
		concurrent: concurrent,
		parser:     parser,
	}
}

//Status function
func (h *Handler) Status(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

//Status function
func (h *Handler) Pictures(w http.ResponseWriter, r *http.Request) {

	startList, ok := r.URL.Query()["start_date"]
	if !ok || len(startList[0]) < 1 {
		h.logger.Warnf("param start is missing")
		w.WriteHeader(400)
		w.Write([]byte("param start is missing"))
		return
	}

	endList, ok := r.URL.Query()["end_date"]
	if !ok || len(endList[0]) < 1 {
		h.logger.Warnf("param end is missing")
		w.WriteHeader(400)
		w.Write([]byte("param end is missing"))
		return
	}

	dates, err := helpers.Interval(startList[0], endList[0])
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
	}

	h.logger.Infof("Dates: %v", dates)
	var isBlocked bool

	wg := sync.WaitGroup{}
	var images []string
	for _, d := range dates {

		if isBlocked {
			break
		}

		*h.concurrent <- 1
		wg.Add(1)
		go func(dateParam string) {
			defer wg.Done()
			val, parseErr := h.parser.Parse(dateParam)
			if parseErr != nil {
				h.logger.WithError(parseErr).Warnf("date: %s", dateParam)
				if parseErr == parsers.OverRateLimit {
					isBlocked = true
				}
			} else {
				h.logger.Infof("Date: %s - Output image: %s",dateParam, val)
				images = append(images, val)
			}

			<-*h.concurrent
		}(d)
	}
	wg.Wait()

	var response = models.Response{URLS: images}
	if isBlocked {
		response.Error = parsers.OverRateLimit.Error()
	}

	res, err := json.Marshal(response)
	if err != nil {
		res, _ = json.Marshal(models.ErrResponse{Error: err.Error()})
	}

	w.Write(res)
	return
}
