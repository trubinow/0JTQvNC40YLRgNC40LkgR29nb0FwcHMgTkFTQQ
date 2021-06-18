package helpers

import (
	"errors"
	"time"
)

var StartDateConversionError = errors.New("start date conversion error")
var EndDateConversionError = errors.New("end date conversion error")
var EndDateBeforeStartDateError = errors.New("end date should not go before start date")

//Interval
func Interval(start string, end string) ([]string, error) {
	startDate, err := time.Parse("2006-01-02", start)
	if err != nil {
		return []string{}, StartDateConversionError
	}

	endDate, err := time.Parse("2006-01-02", end)
	if err != nil {
		return []string{}, EndDateConversionError
	}

	if endDate.Before(startDate) {
		return []string{}, EndDateBeforeStartDateError
	}

	days := int(endDate.Sub(startDate).Hours() / 24)

	var res []string
	for i := 0; i <= days; i++ {
		res = append(res, startDate.AddDate(0, 0, i).Format("2006-01-02"))
	}

	return res, nil
}
