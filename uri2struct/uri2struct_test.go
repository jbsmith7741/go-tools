package uri2struct

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

type testStruct struct {
	// basic types
	String string
	Int    int
	//IntP    *int
	Int32   int32
	Int64   int64
	Float32 float32
	Float64 float64
	//FloatP  *float32

	// struct
	Time  time.Time
	TimeP *time.Time

	// alias
	//Duration time.Duration
}

func TestConvert(t *testing.T) {
	tm, _ := time.Parse(time.RFC3339, "2017-10-10T12:12:12Z")
	cases := []struct {
		msg       string
		uri       string
		shouldErr bool
		expected  testStruct
	}{
		{
			msg:      "simple string",
			uri:      "?String=hello",
			expected: testStruct{String: "hello"},
		},
		{
			msg:      "integers",
			uri:      "?Int=10&Int32=32&Int64=64",
			expected: testStruct{Int: 10, Int32: 32, Int64: 64},
		},
		/*	{
			msg:      "pointer int",
			uri:      "?IntP=77",
			expected: testStruct{IntP: newInt(77)},
		}, */
		{
			msg:       "invalid integer",
			uri:       "?Int=abc",
			shouldErr: true,
		},
		{
			msg:      "floats",
			uri:      "?Float32=12.2&Float64=33.3",
			expected: testStruct{Float32: 12.2, Float64: 33.3},
		},
		{
			msg:      "time.Time",
			uri:      "?Time=2017-10-10T12:12:12Z",
			expected: testStruct{Time: tm},
		},
		/*{
			msg:      "Duration",
			uri:      "?Duration=1h",
			expected: testStruct{Duration: time.Hour},
		},*/
	}

	for _, test := range cases {
		var d testStruct
		err := Convert(&d, test.uri)
		if err != nil != test.shouldErr {
			t.Errorf("FAIL: %v error mismatch %v", test.msg, err)
		} else if !test.shouldErr && !cmp.Equal(d, test.expected) {
			t.Errorf("FAIL: %v values did not match %s", test.msg, cmp.Diff(d, test.expected))
		} else {
			t.Logf("PASS: %v", test.msg)
		}
	}
}

func newInt(i int) *int {
	return &i
}
