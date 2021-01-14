package transport

import "time"

type EventTimings struct {
	Sourced  time.Time
	Updated  time.Time
	Accessed time.Time
}
