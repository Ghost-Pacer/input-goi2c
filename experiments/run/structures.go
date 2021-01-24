package run

import (
	"github.com/google/uuid"
	"gonum.org/v1/gonum/spatial/r3"
	"time"
)

type Route struct {
	id      uuid.UUID
	name    string
	created time.Time

	globalStart r3.Vec
	nodes       []r3.Vec
}

type Activity struct {
	route   *Route
	id      uuid.UUID
	name    string
	created time.Time

	timeTagPeriod time.Duration
	timeTags      []TimeTag
}

type TimeTag struct {
	segment    int
	completion float64
}

/*

"Global" - latitude, longitude, altitude relative to sea level / wall time
"Absolute" - meters relative to start point / time since run start


In route mode there needs to be some idea of a path independent of time so you can compare performance
across multiple activities run along the same route.

If time is sampled as a segment identifier plus a [0,1) segment completion percentage every N milliseconds,
saving user live time is simple provided that the user position is interpolated on demand from the last
known GPS location, heading, & speed. There would be some inherent inaccuracy in that interpolation
but as long as the GPS update period is a small multiple of N this should not be a problem.

Avatar position in route mode can be determined simply by indexing into the appropriate location in
the time sample array (time started since run) and interpolating between the nearest neighbors.

Avatar position, route mode:
	1. Find the time tag interval [t_a, t_b), which is made up of the time tags directly before and after
		the current absolute time. Since time tags are sampled with a known period, simply
		floor and ceil the absolute time to find the beginning and end of the interval.
	2. Compute the effective time tag t_e (segment and completion) of the avatar at the current time.
		This is significantly easier if t_a and t_b refer to the same segment.
		Provided that the segments are on average relatively long and the time tag period is relatively
		short, this should be true maybe 80% of the time.
		I. When t_a and t_b refer to the same segment.
			In this case, t_e's segment is already known, and its completion can be directly interpolated
			since there is a 1:1 correspondence between completion and distance within a single segment.
			t_e.completion = t_a.completion + (t_b.completion - t_a.completion)/(t_b - t_a) * (t_e - t_a)
			Then we can simply obtain the start and end points of the segment s_0 and s_1 and calculate
			the current position as (s_1 - s_0) * t_e.completion.
		II. When t_a and t_b are on different segments.
			This is significantly more complicated because the unit of completion does not necessarily
			correspond to the same quantity of distance across different segments.
			At a high level, you need to compute the total path distance from t_a to t_b along all
			included segments, then divide out t_b-t_a to get the speed along the time interval.
			You can then get the offset distance along the path from s_0, and then you have to
			seek segment by segment subtracting the distance traveled at each node until the
			target segment is found. The position is the position of the start node of that segment
			plus the normalized segment direction vector times the calculated distance along the segment.

User position, route mode:
	Invariant: Maintain two pointers/indices B (backside) and F (frontside) into the nodes array, such that the path
		enclosed by [B, F] has a minimum length w ~= 15 meters, at least two nodes, and contains the user's
		current position roughly at the center (beginning and ending areas excepting).
	Initialization: B is nodes[0], F is the first node such that norm(BF) > w.
	1. Receive the new absolute position p and perform a linear search within [F, B] to find the two nearest
		nodes j, k, such that the vector J.K has a positive component along the vector BF. Then the current segment
		is described by j and k.
	2. Compute the new path position p* as j + proj_J.K(p).
	3. Maintain a timer that ticks off after at least (time tag period) time has elapsed. Each time this ticks,
		store a time tag with { segment: (index of j), completion: (comp_J.K(p) / norm(J.K)) }.
	4. Repeatedly advance F until the distance between the current path position and F norm(P*.F) is greater than
		w/2, maintaining a bounds check on the end of the node array.
	5. Repeatedly advance B until the distance between the current path position and B norm(B.P*) is as small as
		possible while still exceeding w/2, maintaining a bounds check on the beginning of the node array.

Explore mode:
	Avatar position is simply maintained as distance = rate * time, projected onto the user's forward vector.

User position, explore mode:
	1. When a new GPS sample is available, append it on the end of the route under construction at index n.
	2. Optional route optimization step: Attempt to remove sample n-1 if the segment [n-2, n] sufficiently captures
		the information added by holding sample n-1. This can be defined as checking whether the vector
		N-2.N-1 has a direction close enough to the vector N-2.N, by computing the perpendicular component
		of N-2.N-1 with respect to N-2.N norm(N-2.N - proj_N-2.N(N-2.N-1)) and comparing with a fixed threshold y.
		If the component is greater than y, sample n-1 should be maintained. If the component is less than y,
		sample n-1 should be replaced by sample n. If the replacement occurs, all time tags referring to segment
		n-2 need to have their completion recomputed as (old completion) * norm(N-2.N) / norm(N-2.N-1). There may be
		zero, one, or many of these tags but they will be contiguously located at the end of the activity under
		construction.
		You could skip this step and have a constant-time sampled activity and do a post-processing step to do this
		instead.
	3. Perform time tagging on the activity as described in user position, route mode, step 3.


*/
