package trial

import (
	"fmt"
	"reflect"
	"runtime/debug"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type (
	TestFn  func(args ...interface{}) (interface{}, error)
	DiffFn  func(actual interface{}, expected interface{}) string
	EqualFn func(actual interface{}, expected interface{}) bool
)

type Trial struct {
	cases   map[string]Case
	testFn  TestFn
	diffFn  DiffFn
	equalFn EqualFn
}

type Case struct {
	Input    interface{}
	Expected interface{}

	// testing conditions
	ShouldErr   bool  // is an error expected
	ExpectedErr error // the error that was expected (nil is no error expected)
	ShouldPanic bool  // is a panic expected
}

func New(fn TestFn, cases map[string]Case) *Trial {
	if cases == nil {
		cases = make(map[string]Case)
	}
	return &Trial{
		cases:   cases,
		testFn:  fn,
		diffFn:  diffFn,
		equalFn: compareFn,
	}
}

func (t *Trial) EqualFn(fn EqualFn) *Trial {
	t.equalFn = fn
	return t
}

func (t *Trial) DiffFn(fn DiffFn) *Trial {
	t.diffFn = fn
	return t
}

func (trial *Trial) Test(t *testing.T) {
	for msg, test := range trial.cases {
		r := trial.testCase(msg, test)
		if r.Success {
			t.Log(r.Message)
		} else {
			t.Error(r.Message)
		}
	}
}

func (t *Trial) testCase(msg string, test Case) (r result) {
	var finished bool
	defer func() {
		rec := recover()
		if rec == nil && test.ShouldPanic {
			r = fail("FAIL: %q did not panic", msg)
		} else if rec != nil && !test.ShouldPanic {
			r = fail("PANIC: %q %v\n%s", msg, rec, cleanStack())
		} else if !finished {
			r = pass("PASS: %q", msg)
		}
	}()
	var err error
	var result interface{}
	if inputs, ok := test.Input.([]interface{}); ok {
		result, err = t.testFn(inputs...)
	} else {
		result, err = t.testFn(test.Input)
	}

	if (test.ShouldErr && err == nil) || (test.ExpectedErr != nil && err == nil) {
		finished = true
		return fail("FAIL: %q should error", msg)
	} else if !test.ShouldErr && err != nil && test.ExpectedErr == nil {
		finished = true
		return fail("FAIL: %q unexpected error %s", msg, err.Error())
	} else if test.ExpectedErr != nil && !isExpectedError(err, test.ExpectedErr) {
		finished = true
		return fail("FAIL: %q error %q does not match expected %q", msg, err, test.ExpectedErr)
	} else if !test.ShouldErr && test.ExpectedErr == nil && !t.equalFn(result, test.Expected) {
		finished = true
		return fail("FAIL: %q differences %v", msg, t.diffFn(result, test.Expected))
	} else {
		finished = true
		return pass("PASS: %q", msg)
	}
}

// ContainsFn uses the strings.Contain method to compare two interfaces.
// both interfaces need to be strings or implementer the stringer method.
func ContainsFn(actual, expected interface{}) bool {
	// if nothing is expected we have a match
	if expected == nil {
		return true
	}
	s1, ok := actual.(string)
	if !ok {
		s1 = actual.(fmt.Stringer).String()
	}
	s2, ok := expected.(string)
	if !ok {
		s2 = expected.(fmt.Stringer).String()
	}
	return strings.Contains(s1, s2)
}

// compareFn uses the cmp.Equal method to compare two interfaces including unexported fields
func compareFn(actual, expected interface{}) bool {
	t := reflect.TypeOf(actual)
	if reflect.TypeOf(actual) != reflect.TypeOf(expected) {
		return false
	}
	if t != nil && t.Kind() == reflect.Struct {
		return cmp.Equal(actual, expected, cmp.AllowUnexported(actual))
	}
	return cmp.Equal(actual, expected)
}

// diffFn use the cmp.Diff method to display differences between two interfaces
func diffFn(actual, expected interface{}) string {
	var opts []cmp.Option
	if reflect.TypeOf(actual).Kind() == reflect.Struct {
		opts = append(opts, cmp.AllowUnexported(actual))
	}
	if reflect.TypeOf(expected).Kind() == reflect.Struct {
		opts = append(opts, cmp.AllowUnexported(expected))
	}
	return cmp.Diff(actual, expected)
}

// cleanStack removes unhelpful lines from a panic stack track
func cleanStack() (s string) {
	for _, ln := range strings.Split(string(debug.Stack()), "\n") {
		if !strings.Contains(ln, "/go-tools/trial") {
			s += ln + "\n"
		}
	}
	return s
}

func isExpectedError(actual, expected error) bool {
	if err, ok := expected.(errCheck); ok {
		return reflect.TypeOf(actual) == reflect.TypeOf(err.err)
	}
	return strings.Contains(actual.Error(), expected.Error())
}

type errCheck struct {
	err error
}

func (e errCheck) Error() string {
	return e.err.Error()
}

func ErrType(err error) error {
	return errCheck{err}
}

type result struct {
	Success bool
	Message string
}

func pass(format string, args ...interface{}) result {
	return result{
		Success: true,
		Message: fmt.Sprintf(format, args...),
	}
}

func fail(format string, args ...interface{}) result {
	return result{
		Success: false,
		Message: fmt.Sprintf(format, args...),
	}
}
