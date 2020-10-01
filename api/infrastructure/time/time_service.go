package time

import (
	"time"
)

type ITimeInterface interface {
	GetStartOfDayUTC() time.Time
}

type TimeService struct {
}

func NewTimeService() *TimeService {
	m := new(TimeService)
	return m
}

func (e *TimeService) GetStartOfDayUTC() time.Time {
	now := time.Now()
	return time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		0,
		0,
		0,
		0,
		time.UTC)
}
