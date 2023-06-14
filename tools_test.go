package gtools

import (
	"fmt"
	"testing"
	"time"

	"github.com/hydronica/trial"
)

func TestPreciseFloat(t *testing.T) {
	type input struct {
		float  float64
		digits int
	}
	fn := func(in input) (float64, error) {
		return PreciseFloat(in.float, in.digits), nil
	}
	cases := trial.Cases[input, float64]{
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
	fn := func(d time.Duration) (string, error) {
		return PrintDuration(d), nil
	}

	cases := trial.Cases[time.Duration, string]{
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

func TestJoin(t *testing.T) {
	fn := func(in any) (string, error) {
		switch v := in.(type) {
		case []int:
			return Join(v, ","), nil
		case []int16:
			return Join(v, ","), nil
		case []int32:
			return Join(v, ","), nil
		case []int64:
			return Join(v, ","), nil
		case []uint:
			return Join(v, ","), nil
		case []uint16:
			return Join(v, ","), nil
		case []uint32:
			return Join(v, ","), nil
		case []uint64:
			return Join(v, ","), nil
		case []float32:
			return Join(v, ","), nil
		case []float64:
			return Join(v, ","), nil
		default:
			return "", fmt.Errorf("%T is not a supported type", in)
		}
	}
	cases := trial.Cases[any, string]{
		"ints": {
			Input:    []int{1, 2, 3, 4, 5},
			Expected: "1,2,3,4,5",
		},
		"int16": {
			Input:    []int16{1, 2, 3, 4, 5},
			Expected: "1,2,3,4,5",
		},
		"uint64": {
			Input:    []uint64{1, 2, 3, 4, 5},
			Expected: "1,2,3,4,5",
		},
		"float32": {
			Input:    []float32{1.1, 1.2, 4.5},
			Expected: "1.1,1.2,4.5",
		},
		"float64": {
			Input:    []float64{1.1, 1.2, 4.5},
			Expected: "1.1,1.2,4.5",
		},
	}
	trial.New(fn, cases).SubTest(t)
}
