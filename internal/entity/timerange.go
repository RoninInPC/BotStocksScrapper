package entity

import (
	"errors"
	"time"
)

type TimeRange struct {
	start time.Time
	end   time.Time
}

func NewTimeRange(start, end time.Time) (*TimeRange, error) {
	if start.IsZero() || end.IsZero() {
		return nil, errors.New("start or end time is zero")
	}
	if start.After(end) {
		return nil, errors.New("start time must be before end time")
	}
	start = start.UTC()
	end = end.UTC()
	return &TimeRange{start: start, end: end}, nil
}

func (t *TimeRange) Start() time.Time {
	return t.start
}

func (t *TimeRange) End() time.Time {
	return t.end
}

func (t *TimeRange) Contains(timePoint time.Time) bool {
	timePoint = timePoint.UTC()
	if t.start.Equal(timePoint) || t.end.Equal(timePoint) {
		return true
	}
	if t.start.Before(timePoint) && t.end.After(timePoint) {
		return true
	} else {
		return false
	}
}
