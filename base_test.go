package stats

import (
	"math"
	"os"
	"testing"
)

// a very small number
var EPSILON float64

func init() {
	EPSILON = math.Nextafter(1, 2) - 1
}

func TestMain(m *testing.M) {
	//fmt.Println("testmain")
	os.Exit(m.Run())
}

func TestTypes(t *testing.T) {
	var s Sink = NewSink()

	if _, ok := s.(Sink); !ok {
		t.Error("not Sink")
	}
	if _, ok := s.(Variability); !ok {
		t.Error("not Variability")
	}

	got := s.Name()
	if got != welford {
		t.Errorf("got '%v' ; want '%v'", got, welford)
	}
}

func isCloseTo(x, target float64) bool {
	//fmt.Printf("%g // %g\n\n", x-target, EPSILON)
	if math.IsInf(x, 0) && math.IsInf(target, 0) {
		return x == target
	}
	return math.Abs(x-target) < EPSILON
}

// round will round up the nearest integer.
// Source: https://github.com/golang/go/issues/4594#issuecomment-135336012
// Author: James Hartig https://github.com/fastest963
func round(f float64) int {
	if math.Abs(f) < 0.5 {
		return 0
	}
	return int(f + math.Copysign(0.5, f))
}

// toFixed truncates a float into a specified precision
// Source: https://stackoverflow.com/a/29786394
// Author: David Calhoun (Apr 22 15)
// License: http://creativecommons.org/licenses/by-sa/3.0/
func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}

func TestZeroValue(t *testing.T) {
	var v Variability = NewSink()

	if _, ok := v.(Variability); !ok {
		t.Error("not Variability")
	}
	min := v.Min()
	if !isCloseTo(min, math.Inf(-1)) {
		t.Errorf("got '%v' ; want '%v'", min, 0.00)
	}
	max := v.Max()
	if !isCloseTo(max, math.Inf(1)) {
		t.Errorf("got '%v' ; want '%v'", max, 0.00)
	}
}

func TestValidValuesAndCounts(t *testing.T) {
	s := NewSink()
	if err := s.Push(math.NaN()); err == nil {
		t.Errorf("expected %v ; got nil", ErrInvalidValue)
	}
	if err := s.Push(math.Inf(1)); err == nil {
		t.Errorf("expected %v ; got nil", ErrInvalidValue)
	}
}

type fixtureTest struct {
	value1       float64
	value2       float64
	value3       float64
	expectedMean float64
	expectedMin  float64
	expectedMax  float64
	expectedSd   float64
}

var fixtures = []fixtureTest{
	{5.0, 3.0, 1.0, 3.0, 1.0, 5.0, 2.0},
	{2.0, 0.0, 4.0, 2.0, 0.0, 4.0, 2.0},
	{-1.0, 1.0, 12.0, 4.0, -1.0, 12.0, 7.0},
	{-1.1813782, -0.2449577, 0.8799429, -0.182, -1.181, 0.880, 1.032},
	{-0.02735934, -0.62978372, -0.0863240, -0.248, -0.630, -0.027, 0.332},
}

func TestHappyPath(t *testing.T) {
	for _, f := range fixtures {
		s := NewSink()
		_ = s.Push(f.value1)
		_ = s.Push(f.value2)
		_ = s.Push(f.value3)
		//fmt.Printf("mean %v // min %v // max %v // sd %v\n", s.Mean(), s.Min(), s.Max(), s.StandardDeviation())
		//fmt.Printf("truncated sd to 3 decimal places %v\n", toFixed(s.StandardDeviation(), 3))
		if !isCloseTo(toFixed(s.Mean(), 3), f.expectedMean) {
			t.Errorf("got '%v' ; want '%v'", s.Mean(), f.expectedMean)
		}

		if !isCloseTo(toFixed(s.Min(), 3), f.expectedMin) {
			t.Errorf("got '%v' ; want '%v'", s.Min(), f.expectedMin)
		}

		if !isCloseTo(toFixed(s.Max(), 3), f.expectedMax) {
			t.Errorf("got '%v' ; want '%v'", s.Max(), f.expectedMax)
		}

		if !isCloseTo(toFixed(s.StandardDeviation(), 3), f.expectedSd) {
			t.Errorf("got '%v' ; want '%v'", s.StandardDeviation(), f.expectedSd)
		}
	}
}
