package uri2struct

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

type testStruct struct {
	// basic types
	String   string
	Int      int
	IntP     *int
	Int32    int32
	Int32P   *int32
	Int64    int64
	Int64P   *int64
	Float32  float32
	Float32P *float32
	Float64  float64
	Float64P *float64

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
			msg:      "string",
			uri:      "?String=hello",
			expected: testStruct{String: "hello"},
		},
		{
			msg:      "integer: int, int32, int64",
			uri:      "?Int=10&Int32=32&Int64=64",
			expected: testStruct{Int: 10, Int32: 32, Int64: 64},
		},
		{
			msg:      "pointer: *int, *int32, *int64",
			uri:      "?IntP=77&Int32P=11&Int64P=222",
			expected: testStruct{IntP: newInt(77), Int32P: newInt32(11), Int64P: newInt64(222)},
		},
		{
			msg:       "invalid integer",
			uri:       "?Int=abc",
			shouldErr: true,
		},
		{
			msg:      "float32, float64",
			uri:      "?Float32=12.2&Float64=33.3",
			expected: testStruct{Float32: 12.2, Float64: 33.3},
		},
		{
			msg:      "pointer: *float32, *float64",
			uri:      "?Float32P=12.2&Float64P=33.3",
			expected: testStruct{Float32P: newFloat32(12.2), Float64P: newFloat64(33.3)},
		},
		{
			msg:      "time.Time",
			uri:      "?Time=2017-10-10T12:12:12Z",
			expected: testStruct{Time: tm},
		},
		{
			msg:      "*time.Time",
			uri:      "?TimeP=2017-10-10T12:12:12Z",
			expected: testStruct{TimeP: &tm},
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
func newInt32(i int32) *int32 {
	return &i
}
func newInt64(i int64) *int64 {
	return &i
}
func newFloat32(f float32) *float32 {
	return &f
}
func newFloat64(f float64) *float64 {
	return &f
}
