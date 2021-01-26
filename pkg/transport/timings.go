package transport

import "time"

// EventTimings captures timing information corresponding to a single sample as:
//	.Sourced: 	value read time, natively/as close as possible to the actual peripheral
//	.Updated: 	value update time on transport
//	.Accessed:	value access time on transport
type EventTimings struct {
	Sourced  time.Time
	Updated  time.Time
	Accessed time.Time
}

// yes it's dumb to have this in its own file but the generator duplicates it for some reason
// if it's in transport.go (bug?)
