package models

import "time"

type Policy struct {
	MaxRequests int
	Interval    time.Duration
}
