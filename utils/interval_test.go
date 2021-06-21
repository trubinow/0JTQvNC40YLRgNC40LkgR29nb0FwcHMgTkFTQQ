package utils

import (
	"errors"
	"reflect"
	"testing"
	"time"
)

func TestInterval(t *testing.T) {
	cases := []struct {
		Start  string
		End string
		Result []string
		Err   error
	}{
		{
			"2021-05-30", "2021-06-02", []string{"2021-05-30", "2021-05-31", "2021-06-01", "2021-06-02"}, nil,
		},

		{
			"2021-06-11", "2021-06-12", []string{"2021-06-11", "2021-06-12"}, nil,
		},

		{
			"2021-06-11", "2021-06-11", []string{"2021-06-11"}, nil,
		},

		{
			"2021-06-11", "2021-06-10", []string{}, ErrEndDateBeforeStartDate,
		},

		{
			"2021-0a-11", "2021-06-10", []string{}, ErrStartDateConversion,
		},

		{
			"2021-06-11", "2021-13-10", []string{}, ErrEndDateConversion,
		},

		{
			"2021-06-11", "", []string{}, ErrEndDateConversion,
		},

		{
			"1995-06-11", "2006-06-09", []string{}, ErrWrongInterval,
		},

		{
			"2020-06-11", time.Now().Add(24*time.Hour).Format("2006-01-02"), []string{}, ErrWrongInterval,
		},
	}

	for _, d := range cases {
		result, err := Interval(d.Start, d.End)
		if  !errors.Is(err, d.Err) || !reflect.DeepEqual(result, d.Result) {
			t.Errorf("Error: %v Interval(%s,%s) = %v; want %v", err, d.Start, d.End, result, d.Result)
		}
	}
}
