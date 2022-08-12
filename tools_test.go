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

const (
	IByte = 1
	KByte = IByte * 1000
	MByte = KByte * 1000
	GByte = MByte * 1000
	TByte = GByte * 1000
	PByte = TByte * 1000
	EByte = PByte * 1000
)

func TestBytes(t *testing.T) {
	type testStruct struct {
		name     string
		actual   string
		expected string
	}
	tests := []testStruct{
		{"bytes(0)", Bytes(0), "0 B"},
		{"bytes(1)", Bytes(1), "1 B"},
		{"bytes(803)", Bytes(803), "803 B"},
		{"bytes(999)", Bytes(999), "999 B"},

		{"bytes(1024)", Bytes(1024), "1.0 kB"},
		{"bytes(9999)", Bytes(9999), "10 kB"},
		{"bytes(1MB - 1)", Bytes(MByte - 2014), "1000 kB"},

		{"bytes(1MB)", Bytes(1024 * 1024), "1.0 MB"},
		{"bytes(1GB - 1K)", Bytes(GByte - KByte), "1000 MB"},

		{"bytes(1GB)", Bytes(GByte), "1.0 GB"},
		{"bytes(1TB - 1M)", Bytes(TByte - MByte), "1000 GB"},
		{"bytes(10MB)", Bytes(9999 * 1000), "10 MB"},

		{"bytes(1TB)", Bytes(TByte), "1.0 TB"},
		{"bytes(1PB - 1T)", Bytes(PByte - TByte), "999 TB"},

		{"bytes(1PB)", Bytes(PByte), "1.0 PB"},
		{"bytes(1PB - 1T)", Bytes(EByte - PByte), "999 PB"},

		{"bytes(1EB)", Bytes(EByte), "1.0 EB"},
	}
	for _, v := range tests {
		if v.actual != v.expected {
			t.Errorf("FAIL: %v %q!=%q", v.name, v.actual, v.expected)
		}
	}
}
