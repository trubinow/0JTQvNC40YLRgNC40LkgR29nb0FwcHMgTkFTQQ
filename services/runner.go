package services

import (
	"github.com/sirupsen/logrus"
	"nasa/parsers"
	"sync"
)

//Parser obtains nasa api data
type Parser interface {
	Parse(targetDate string) (string, error)
}

//Runner starts parsing process
type Runner struct {
	blocked    bool
	err	error
	logger     *logrus.Entry
	parser     Parser
	concurrent chan struct{}
}

//NewRunner creates new parser runner
func NewRunner(logger *logrus.Entry, concurrent chan struct{}, parser Parser) *Runner {
	return &Runner{
		blocked:    false,
		logger:     logger,
		concurrent: concurrent,
		parser:     parser,
	}
}

func (r *Runner) Blocked() bool {
	return r.blocked
}

func (r *Runner) SetBlocked(b bool) {
	r.blocked = b
}

func (r *Runner) Error() error {
	return r.err
}

func (r *Runner) SetError(err error) {
	r.err = err
}

//Run starts parse process
func (r *Runner) Run(dates []string) ([]string, error) {
	r.SetBlocked(false)
	wg := sync.WaitGroup{}
	var images []string
	for _, d := range dates {

		if r.Blocked() || r.Error() != nil {
			break
		}

		r.concurrent <- struct{}{}
		wg.Add(1)
		go func(dateParam string) {
			defer wg.Done()
			val, parseErr := r.parser.Parse(dateParam)
			if parseErr != nil {
				r.logger.WithError(parseErr).Warnf("date: %s", dateParam)
				if parseErr == parsers.ErrOverRateLimit {
					r.SetBlocked(true)
				}
				r.SetError(parseErr)
			} else {
				r.logger.Infof("date: %s - output image: %s", dateParam, val)
				images = append(images, val)
			}

			<-r.concurrent
		}(d)
	}
	wg.Wait()

	return images, r.Error()
}
