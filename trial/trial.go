package trial

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type (
	TestFn  func(args ...interface{}) (interface{}, error)
	DiffFn  func(interface{}, interface{}) string
	EqualFn func(interface{}, interface{}) bool
)

type Trial struct {
	cases   map[string]Case
	testFn  TestFn
	diffFn  DiffFn
	equalFn EqualFn
}

func New(fn TestFn, cases map[string]Case) *Trial {
	if cases == nil {
		cases = make(map[string]Case)
	}
	return &Trial{
		cases:  cases,
		testFn: fn,
		diffFn: func(i1, i2 interface{}) string {
			return cmp.Diff(i1, i2)
		},
		equalFn: func(i1, i2 interface{}) bool {
			return cmp.Equal(i1, i2)
		},
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

type Case struct {
	Input    interface{}
	Expected interface{}

	// testing conditions
	ShouldErr   bool  // is an error expected
	ExpectedErr error // the error that was expected (nil is no error expected)
	ShouldPanic bool  // is a panic expected
}

/*func (t *Trial) Add(msg string, c Case) *Trial {
	if _, found := t.cases[msg]; found {
		// todo: considering changing to t.Fatalf (t *testing.T)
		log.Fatalf("test case %q already exists", msg)
	}
	t.cases[msg] = c
	return t
}*/

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
			r = fail("PANIC: %q %v", msg, rec)
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

	if test.ShouldErr && err == nil {
		finished = true
		return fail("FAIL: %q should error", msg)
	} else if !test.ShouldErr && err != nil {
		finished = true
		return fail("FAIL: %q unexpected error %s", msg, err.Error())
	} else if !test.ShouldErr && !t.equalFn(result, test.Expected) {
		finished = true
		return fail("FAIL: %q differences %v", msg, t.diffFn(result, test.Expected))
	} else {
		finished = true
		return pass("PASS: %q", msg)
	}
	// todo add error type check test
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

func Interfaces(args ...interface{}) interface{} {
	return args
}
