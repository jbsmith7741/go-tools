package uri2struct

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

type testStruct struct {
	String  string
	Int     int
	Int32   int32
	Int64   int64
	Float32 float32
	Float64 float64
}

func TestConvert(t *testing.T) {
	cases := []struct {
		msg       string
		uri       string
		shouldErr bool
		expected  testStruct
	}{
		{
			msg:      "simple string",
			uri:      "?string=hello",
			expected: testStruct{String: "hello"},
		},
		{
			msg:      "integers",
			uri:      "?int=10&int32=32&int64=64",
			expected: testStruct{Int: 10, Int32: 32, Int64: 64},
		},
		{
			msg:       "invalid integer",
			uri:       "?int=abc",
			shouldErr: true,
		},
		{
			msg:      "floats",
			uri:      "?float32=12.2&float64=33.3",
			expected: testStruct{Float32: 12.2, Float64: 33.3},
		},
	}

	for _, test := range cases {
		var d testStruct
		err := Convert(&d, test.uri)
		if !(err == nil || test.shouldErr) {
			t.Errorf("FAIL: %v error mismatch %v", test.msg, err)
		} else if !test.shouldErr && !cmp.Equal(d, test.expected) {
			t.Errorf("FAIL: %v values did not match %s", test.msg, cmp.Diff(d, test.expected))
		} else {
			t.Logf("PASS: %v", test.msg)
		}
	}
}
