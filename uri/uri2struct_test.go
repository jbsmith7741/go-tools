package uri

import (
	"log"
	"reflect"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

type testStruct struct {
	// basic types
	String   string
	Bool     bool
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

	// slice
	Strings  []string
	Ints     []int
	IntsP    []*int
	Ints32   []int32
	Ints64   []int64
	Floats32 []float32
	Floats64 []float64

	// struct
	Time       time.Time
	TimeP      *time.Time
	Unmarshal  unmarshalStruct
	UnmarshalP *unmarshalStruct

	// alias
	Dessert dessert
}

type unmarshalStruct struct {
	Data string
}

func (s *unmarshalStruct) UnmarshalText(text []byte) error {
	s.Data = string(text)
	return nil
}

func TestUnmarshal(t *testing.T) {
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
			msg:       "invalid float",
			uri:       "?Float32=abc",
			shouldErr: true,
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
		{
			msg:       "invalid time",
			uri:       "?Time=2017-10-",
			shouldErr: true,
		},
		{
			msg: "struct with UnMarshalText",
			uri: "?Unmarshal=abc&UnmarshalP=def",
			expected: testStruct{
				Unmarshal:  unmarshalStruct{Data: "abc"},
				UnmarshalP: &unmarshalStruct{Data: "def"},
			},
		},
		{
			msg: "bool",
			uri: "?Bool=true",
			expected: testStruct{
				Bool: true,
			},
		},
		{
			msg: "bool implicit true",
			uri: "?Bool&Test",
			expected: testStruct{
				Bool: true,
			},
		},
		{
			msg: "slice of string",
			uri: "?Strings=a&Strings=b&Strings=c",
			expected: testStruct{
				Strings: []string{"a", "b", "c"},
			},
		},
		{
			msg: "slice: int, int32, int64",
			uri: "?Ints=1&Ints=2&Ints=3&Ints32=4,5,6&Ints64=7,8,9",
			expected: testStruct{
				Ints:   []int{1, 2, 3},
				Ints32: []int32{4, 5, 6},
				Ints64: []int64{7, 8, 9},
			},
		},
		{
			msg: "slice: float32, float64",
			uri: "?Floats32=1.1&Floats32=2.2&Floats32=3.3&Floats64=4.4,5.5,6.6",
			expected: testStruct{
				Floats32: []float32{1.1, 2.2, 3.3},
				Floats64: []float64{4.4, 5.5, 6.6},
			},
		},
		{
			msg: "slice of *int",
			uri: "?IntsP=1,2,3",
			expected: testStruct{
				IntsP: []*int{newInt(1), newInt(2), newInt(3)},
			},
		},
		{
			msg:      "alias type (dessert)",
			uri:      "?Dessert=brownie",
			expected: testStruct{Dessert: brownie},
		},
		{
			msg:       "invalid alias type",
			uri:       "?Dessert=cat",
			shouldErr: true,
		},
	}
	for _, test := range cases {
		var d testStruct
		err := Unmarshal(&d, test.uri)
		if err != nil != test.shouldErr {
			t.Errorf("FAIL: %v error mismatch %v", test.msg, err)
		} else if !test.shouldErr && !cmp.Equal(d, test.expected) {
			t.Errorf("FAIL: %v values did not match %s", test.msg, cmp.Diff(d, test.expected))
		} else {
			t.Logf("PASS: %v", test.msg)
		}
	}
}

type (
	testScheme struct {
		Schema string `uri:"scheme"`
	}
	testHost struct {
		Host string `uri:"host"`
	}
	testPath struct {
		Path string `uri:"path"`
	}
	testOrigin struct {
		Origin string `uri:"Origin"`
	}
	testAuthority struct {
		Authority string `uri:"authority"`
	}
	testCustom struct {
		OldInt int `uri:"NewInt"`
		Host   string
	}
	testPrivate struct {
		int    int    `uri:"int"`
		String string `uri:"string"`
	}
	testPrivate2 struct {
		int int
		Int int `uri:"int"`
	}
)

func TestTags(t *testing.T) {
	cases := []struct {
		msg      string
		uri      string
		expected interface{}
	}{
		{
			msg:      "Scheme uri tag",
			uri:      "https://localhost:8080/usr/bin",
			expected: &testScheme{Schema: "https"},
		},
		{
			msg:      "Host uri tag",
			uri:      "https://localhost:8080/usr/bin",
			expected: &testHost{Host: "localhost:8080"},
		},
		{
			msg:      "Path uri tag",
			uri:      "https://localhost:8080/usr/bin",
			expected: &testPath{Path: "/usr/bin"},
		},
		{
			msg:      "Authority uri tag",
			uri:      "https://localhost:8080/usr/bin",
			expected: &testAuthority{Authority: "https://localhost:8080"},
		},
		{
			msg:      "Origin uri tag",
			uri:      "https://localhost:8080/usr/bin",
			expected: &testOrigin{Origin: "https://localhost:8080/usr/bin"},
		},
		{
			msg:      "Origin uri tag without authority",
			uri:      "/usr/bin",
			expected: &testOrigin{Origin: "/usr/bin"},
		},
		{
			msg:      "Custom int name",
			uri:      "?NewInt=10",
			expected: &testCustom{OldInt: 10},
		},
		{
			msg:      "Var named Host without tag",
			uri:      "https://local/usr/bin?Host=hello",
			expected: &testCustom{Host: "hello"},
		},
	}
	for _, test := range cases {
		v := reflect.New(reflect.TypeOf(test.expected).Elem()).Interface()
		Unmarshal(v, test.uri)
		if !cmp.Equal(v, test.expected) {
			t.Errorf("FAIL: %v values did not match %s", test.msg, cmp.Diff(v, test.expected))
		} else {
			t.Logf("PASS: %v", test.msg)
		}
	}
}

func TestValidate(t *testing.T) {
	cases := []struct {
		msg       string
		uri       string
		data      interface{}
		shouldErr bool
	}{
		{
			msg:       "Cannot write to struct",
			data:      struct{}{},
			shouldErr: true,
		},
		{
			msg:       "Invalid uri",
			uri:       "://",
			data:      &struct{}{},
			shouldErr: true,
		},
		{
			msg:  "private variables",
			uri:  "?string=hello&int=1",
			data: &testPrivate{},
		},
		{
			msg:  "private variables",
			uri:  "int=1",
			data: &testPrivate2{Int: 1},
		},
	}
	for _, test := range cases {
		err := Unmarshal(test.data, test.uri)
		if err != nil != test.shouldErr {
			t.Errorf(test.msg)
		} else {
			log.Printf("PASS: %q data: %v", test.msg, test.data)
		}
	}
}
