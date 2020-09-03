package gtools

import (
	"testing"

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
