package stats

import (
	"errors"
	"math"
)

// Sink accepts a stream of observations.
type Sink interface {
	Push(x float64) error
	Name() string
}

// Variability is a measure of how spread out is your dataset.
type Variability interface {
	Min() float64
	Max() float64
	Mean() float64
	StandardDeviation() float64
	Count() int
}

const welford = "B. P. Welford's 1962 algorithm"

// ErrInvalidValue means that the item you supplied is a NaN or an Inf
// which is considered illegal.
var ErrInvalidValue = errors.New("NaN or +-Inf are not allowed")

// WelfordSink is an implementation of a Sink and Variability
// @TODO make this goroutine-safe
type WelfordSink struct {
	Sink
	Variability
	observationMin float64
	observationMax float64
	count          int
	mOld           float64
	mNew           float64
	nNew           float64
	nOld           float64
}

// NewSink returns a pointer to an implementation of interfaces Sink & Variability
func NewSink() *WelfordSink {
	return &WelfordSink{
		observationMin: math.Inf(-1),
		observationMax: math.Inf(1)}
}

func (s *WelfordSink) variance() float64 {
	if s.count > 1 {
		return s.nNew / float64(s.count-1)
	}
	return 0.0
}

// StandardDeviation returns the running sd of the observations
func (s *WelfordSink) StandardDeviation() float64 {
	return math.Sqrt(s.variance())
}

// Push accepts an observation.
// On failure, it returns an ErrInvalidValue.
func (s *WelfordSink) Push(x float64) error {
	if math.IsNaN(x) || math.IsInf(x, 0) {
		return ErrInvalidValue
	}

	s.count++

	if s.count == 1 {
		s.mOld = x
		s.mNew = x
		s.nOld = 0.0
		s.nNew = 0.0
	} else {
		s.mNew = s.mOld + (x-s.mOld)/float64(s.count)
		s.nNew = s.nOld + (x-s.mOld)*(x-s.mNew)

		s.mOld = s.mNew
		s.nOld = s.nNew
	}

	// @TODO move to internal function
	// is this a candidate for minimum?
	if math.IsInf(s.observationMin, 0) {
		if x > s.observationMin {
			s.observationMin = x
		}
	} else {
		if x < s.observationMin {
			s.observationMin = x
		}
	}

	if math.IsInf(s.observationMax, 0) {
		if x < s.observationMax {
			s.observationMax = x
		}
	} else {
		if x > s.observationMax {
			s.observationMax = x
		}
	}

	return nil
}

// Name returns a descriptive comment
func (s *WelfordSink) Name() string {
	return welford
}

// Min returns the smallest observation seen.
func (s *WelfordSink) Min() float64 {
	return s.observationMin
}

// Max returns the largest observation seen.
func (s *WelfordSink) Max() float64 {
	return s.observationMax
}

// Count returns the number of observations that was seen thus far.
func (s *WelfordSink) Count() int {
	return s.count
}

// Mean returns the arithmetic mean of all the observations.
func (s *WelfordSink) Mean() float64 {
	if s.count > 0 {
		return s.mNew
	}
	return float64(s.count)
}
