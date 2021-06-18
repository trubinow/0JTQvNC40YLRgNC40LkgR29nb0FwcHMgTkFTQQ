package helpers

import (
	"errors"
	"reflect"
	"testing"
)

func TestInterval(t *testing.T) {
	cases := []struct {
		Start  string
		End string
		Result []string
		Err   error
	}{
		{
			"2021-06-11", "2021-06-12", []string{"2021-06-11", "2021-06-12"}, nil,
		},

		{
			"2021-06-11", "2021-06-11", []string{"2021-06-11"}, nil,
		},

		{
			"2021-06-11", "2021-06-10", []string{}, EndDateBeforeStartDateError,
		},

		{
			"2021-0a-11", "2021-06-10", []string{}, StartDateConversionError,
		},

		{
			"2021-06-11", "2021-13-10", []string{}, EndDateConversionError,
		},

	}

	for _, d := range cases {
		result, err := Interval(d.Start, d.End)
		if  !errors.Is(err, d.Err) || !reflect.DeepEqual(result, d.Result) {
			t.Errorf("Error: %v Interval(%s,%s) = %v; want %v", err, d.Start, d.End, result, d.Result)
		}
	}
}
