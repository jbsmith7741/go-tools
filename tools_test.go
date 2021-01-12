package gtools

import (
	"testing"
	"time"

	"github.com/hydronica/trial"
)

func TestPreciseFloat(t *testing.T) {
	type input struct {
		float  float64
		digits int
	}
	fn := func(i trial.Input) (interface{}, error) {
		in := i.Interface().(input)
		return PreciseFloat(in.float, in.digits), nil
	}
	cases := trial.Cases{
		"3.45_": {
			Input:    input{3.45678, 2},
			Expected: 3.45,
		},
		"no rounding": {
			Input:    input{float64(3.99999999999999), 4},
			Expected: 3.9999,
		},
		"unnecessary truncation": {
			Input:    input{12.0, 2},
			Expected: 12.0,
		},
		"-3.45": {
			Input:    input{-3.45678, 2},
			Expected: -3.45,
		},
		"negative digit": {
			Input:    input{1234.5678, -2},
			Expected: 1200.0,
		},
		"zero digit": {
			Input:    input{1234.5678, 0},
			Expected: 1234.0,
		},
	}

	trial.New(fn, cases).SubTest(t)
}

func TestPrintDuration(t *testing.T) {
	fn := func(in trial.Input) (interface{}, error) {
		return PrintDuration(in.Interface().(time.Duration)), nil
	}

	cases := trial.Cases{
		"nano": {
			Input:    12 * time.Nanosecond,
			Expected: "12ns",
		},
		"micro": {
			Input:    17 * time.Microsecond,
			Expected: "17µs",
		},
		"milli": {
			Input:    28 * time.Millisecond,
			Expected: "28ms",
		},
		"milli+micro+nano": {
			Input:    28*time.Millisecond + 17*time.Microsecond + 12*time.Nanosecond,
			Expected: "28.01ms",
		},
		"second": {
			Input:    1234 * time.Millisecond,
			Expected: "1.23s",
		},
		"minute": {
			Input:    12*time.Minute + (35600 * time.Millisecond),
			Expected: "12m35.6s",
		},
		"hour": {
			Input:    12 * time.Hour,
			Expected: "12h",
		},
		"hour+sec": {
			Input:    2*time.Hour + 12*time.Second,
			Expected: "2h12s",
		},
		"day": {
			Input:    72 * time.Hour,
			Expected: "3d",
		},
		"year": {
			Input:    2 * durYear,
			Expected: "2y",
		},
		"full year": {
			Input:    3*durYear + 120*durDay + 7*time.Hour + 31*time.Minute + 12*time.Second + 120*time.Millisecond,
			Expected: "3y120d7h31m12.12s",
		},
		"neg full year": {
			Input:    -3*durYear - 120*durDay - 7*time.Hour - 31*time.Minute - 12*time.Second - 120*time.Millisecond,
			Expected: "-3y120d7h31m12.12s",
		},
		"neg milli": {
			Input:    -28 * time.Millisecond,
			Expected: "-28ms",
		},
		"neg day": {
			Input:    -72 * time.Hour,
			Expected: "-3d",
		},
		"neg hour+sec": {
			Input:    -2*time.Hour - 12*time.Second,
			Expected: "-2h12s",
		},
	}
	trial.New(fn, cases).Test(t)
}
