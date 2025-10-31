package main

import (
	"time"
)

func ConvertTime(date string) (time.Time, error) {
	return time.Parse(DefaultTimeFormat, date)
}
