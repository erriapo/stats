package stats

import (
	"fmt"
	"math"
	"os"
	"reflect"
	"testing"
)

// a very small number
var EPSILON float64

func init() {
	EPSILON = math.Nextafter(1, 2) - 1
}

func TestMain(m *testing.M) {
	fmt.Println("testmain")
	fmt.Println()
	os.Exit(m.Run())
}

func TestTypes(t *testing.T) {
	var s Sink
	s = NewSink()
	fmt.Println(reflect.TypeOf(s))

	_, ok := s.(Sink)
	fmt.Printf("%v#\n", ok)
	if !ok {
		t.Error("not Sink")
	}
	_, ok = s.(Variability)
	if !ok {
		t.Error("not Variability")
	}

	got := s.Name()
	if got != welford {
		t.Errorf("got '%v' ; want '%v'", got, welford)
	}
}

func isCloseTo(x, target float64) bool {
	fmt.Printf("%g // %g\n\n", x-target, EPSILON)
	if math.IsInf(x, 0) && math.IsInf(target, 0) {
		return x == target
	} else {
		return math.Abs(x-target) < EPSILON
	}
}

func TestZeroValue(t *testing.T) {
	var v Variability
	v = NewSink()
	fmt.Println(reflect.TypeOf(v))

	_, ok := v.(Variability)
	if !ok {
		t.Error("not Variability")
	}
	min := v.Min()
	if !isCloseTo(min, math.Inf(-1)) {
		t.Errorf("got '%v' ; want '%v'", min, 0.00)
	}
}

func TestValidValuesAndCounts(t *testing.T) {
	s := NewSink()
	if err := s.Push(math.NaN()); err == nil {
		t.Errorf("expected %v ; got nil", err)
	}
}

func TestMean(t *testing.T) {
	s := NewSink()
	s.Push(5.0)
	s.Push(3.0)
	s.Push(1.0)
	mean := s.Mean()
	fmt.Printf("mean is %g\n", mean)
	if !isCloseTo(mean, 3.0) {
		t.Errorf("got '%v' ; want '%v'", mean, 3.0)
	}

	if !isCloseTo(s.Min(), 1.0) {
		t.Errorf("got '%v' ; want '%v'", s.Min(), 1.0)
	}
}
