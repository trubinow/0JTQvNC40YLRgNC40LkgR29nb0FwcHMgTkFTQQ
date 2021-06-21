package utils

import (
	"errors"
	"time"
)

var (
	ErrStartDateConversion    = errors.New("start date conversion error")
	ErrEndDateConversion      = errors.New("end date conversion error")
	ErrEndDateBeforeStartDate = errors.New("end date should not go before start date")
	ErrWrongInterval          = errors.New("date interval must be between Jun 16, 1995 and today")
)

var bottomDateLimit time.Time
var err error

func init() {
	bottomDateLimit, err = time.Parse("2006-01-02", "1995-06-16")
	if err != nil {
		panic(err)
	}
}

//Interval returns dates slice in string format between start-date and end-date. Date output and input format is 2006-01-02
func Interval(start string, end string) ([]string, error) {
	startDate, err := time.Parse("2006-1-2", start)
	if err != nil {
		return []string{}, ErrStartDateConversion
	}

	endDate, err := time.Parse("2006-1-2", end)
	if err != nil {
		return []string{}, ErrEndDateConversion
	}

	if startDate.Before(bottomDateLimit) || startDate.After(time.Now()) || endDate.Before(bottomDateLimit) || endDate.After(time.Now()) {
		return []string{}, ErrWrongInterval
	}

	if endDate.Before(startDate) {
		return []string{}, ErrEndDateBeforeStartDate
	}

	days := int(endDate.Sub(startDate).Hours() / 24)

	var res []string
	for i := 0; i <= days; i++ {
		res = append(res, startDate.AddDate(0, 0, i).Format("2006-01-02"))
	}

	return res, nil
}
