package days

import (
	"errors"
	"time"
)

var ErrWrongRange = errors.New("from bigger then to")

func DaysBetween(from time.Time, to time.Time) ([]time.Time, error) {
	if from.After(to) {
		return nil, ErrWrongRange
	}

	days := make([]time.Time, 0, to.Day()-from.Day()+1)
	current := ToDay(from)
	end := ToDay(to)

	for !current.After(end) {
		days = append(days, current)
		current = current.AddDate(0, 0, 1) // Increment by one day.
	}

	return days, nil
}

func ToDay(timestamp time.Time) time.Time {
	return time.Date(timestamp.Year(), timestamp.Month(), timestamp.Day(), 0, 0, 0, 0, time.UTC).Round(0)
}

func Date(year, month, day int) time.Time {
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC).Round(0)
}
