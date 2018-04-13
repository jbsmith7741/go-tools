package trial

import (
	"fmt"
	"strings"
	"time"
)

// ContainsFn uses the strings.Contain method to compare two interfaces.
// both interfaces need to be strings or implementer the stringer method.
func ContainsFn(i1, i2 interface{}) bool {
	s1, ok := i1.(string)
	if !ok {
		s1 = i1.(fmt.Stringer).String()
	}
	s2, ok := i2.(string)
	if !ok {
		s2 = i2.(fmt.Stringer).String()
	}
	return strings.Contains(s1, s2)
}

func Args(args ...interface{}) interface{} {
	return args
}

// IntP returns a pointer to a defined int
func IntP(i int) *int {
	return &i
}

// Int8P returns a pointer to a defined int8
func Int8P(i int8) *int8 {
	return &i
}

// Int16P returns a pointer to a defined int16
func Int16P(i int16) *int16 {
	return &i
}

// Int32P returns a pointer to a defined int32
func Int32P(i int32) *int32 {
	return &i
}

// Int64P returns a pointer to a defined int64
func Int64P(i int64) *int64 {
	return &i
}

// Float32P returns a pointer to a defined float32
func Float32P(f float32) *float32 {
	return &f
}

// Float64P returns a pointer to a defined float64
func Float64P(f float64) *float64 {
	return &f
}

// Time is a panic wrapper for the time.Parse method
// it returns a time.Time for the given layout and value
func Time(layout, value string) time.Time {
	t, err := time.Parse(layout, value)
	if err != nil {
		panic(err)
	}
	return t
}

func Times(layout string, values ...string) []time.Time {
	times := make([]time.Time, len(values))
	for i, v := range values {
		times[i] = Time(layout, v)
	}
	return times
}

// TimeP return a pointers to a time.Time for the given layout and value.
// it panics on error
func TimeP(layout, value string) *time.Time {
	t := Time(layout, value)
	return &t
}
