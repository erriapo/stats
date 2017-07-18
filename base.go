package stats

import (
	"errors"
	"math"
)

type Sink interface {
	Push(x float64) error
	Name() string
}

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

// @TODO make this goroutine-safe
type WelfordSink struct {
	Sink
	Variability
	observationMin float64
	count          int
	mOld           float64
	mNew           float64
	nNew           float64
	nOld           float64
}

func NewSink() *WelfordSink {
	return &WelfordSink{}
}

func (s *WelfordSink) Push(x float64) error {
	if math.IsNaN(x) || math.IsInf(x, 0) {
		return ErrInvalidValue
	}

	s.count++

	if s.count == 1 {
		s.mOld = x
		s.mNew = x
		s.nOld = 0.0
	} else {
		s.mNew = s.mOld + (x-s.mOld)/float64(s.count)
		s.nNew = s.nOld + (x-s.mOld)*(x-s.mNew)

		s.mOld = s.mNew
		s.nOld = s.mNew
	}

	return nil
}

func (s *WelfordSink) Name() string {
	return welford
}

func (s *WelfordSink) Min() float64 {
	return s.observationMin
}

func (s *WelfordSink) Count() int {
	return s.count
}

func (s *WelfordSink) Mean() float64 {
	if s.count > 0 {
		return s.mNew
	} else {
		return float64(s.count)
	}
}
