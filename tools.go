package gtools

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

// PreciseFloat truncates a float to the specified number of digits
func PreciseFloat(f float64, digits int) float64 {
	pow := math.Pow(10.0, float64(digits))
	return math.Trunc(pow*f) / pow
}

const (
	durDay  = 24 * time.Hour
	durYear = 365 * durDay
)

// PrintDuration converts a duration into a human friendly time
func PrintDuration(d time.Duration) (s string) {
	if d == 0 {
		return "0s"
	}
	// deal with negatives
	if d < 0 {
		s = "-"
		d *= -1
	}
	if d < time.Microsecond {
		return s + strconv.FormatInt(int64(d), 10) + "ns"
	}
	if d < time.Millisecond {

		return s + fDuration(d, time.Microsecond) + "µs"
	}
	if d < time.Second {
		return s + fDuration(d, time.Millisecond) + "ms"
	}
	// add year
	y := int(d / durYear)
	if y != 0 {
		d -= time.Duration(y) * durYear
		s += strconv.Itoa(y) + "y"
	}

	// add day
	day := int(d / durDay)
	if day != 0 {
		d -= time.Duration(day) * durDay
		s += strconv.Itoa(day) + "d"
	}

	// add hour
	h := int(d / time.Hour)
	if h != 0 {
		d -= time.Duration(h) * time.Hour
		s += strconv.Itoa(h) + "h"
	}

	// add minute
	m := int(d / time.Minute)
	if m != 0 {
		d -= time.Duration(m) * time.Minute
		s += strconv.Itoa(m) + "m"
	}
	// add second
	if sec := fDuration(d, time.Second) + "s"; sec != "0s" {
		s += sec
	}
	return s
}

// fDuration returns a duration as a two digit approximation of the precision provided
func fDuration(d time.Duration, precision time.Duration) string {
	i := int(d / (precision / 100)) // truncate to 3 digits
	s := strconv.Itoa(i / 100)      // first digit is whole number
	i %= 100
	if i == 0 {
		return s
	}
	s += "." + strconv.Itoa(i/10) // first decimal
	i %= 10
	if i == 0 {
		return s
	}
	return s + strconv.Itoa(i) // second decimal
}

type ints interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}
type floats interface {
	~float32 | ~float64
}

type numbers interface {
	ints | floats
}

// Join like strings.Join but for any slice Type
func Join[T numbers](elems []T, sep string) string {
	s := make([]string, len(elems))
	for i, v := range elems {
		s[i] = fmt.Sprint(v)
	}
	return strings.Join(s, sep)
}
