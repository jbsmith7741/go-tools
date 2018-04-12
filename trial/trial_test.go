package trial

import (
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"
)

func TestTrial_TestCase(t *testing.T) {

	divideFn := func(args ...interface{}) (interface{}, error) {
		return func(a, b int) (int, error) {
			if b == 0 {
				return 0, errors.New("divide by zero")
			}
			return a / b, nil
		}(args[0].(int), args[1].(int))
	}

	panicFn := func(args ...interface{}) (interface{}, error) {
		return func(s string) string {
			t, err := time.Parse(time.RFC3339, s)
			if err != nil {
				panic(err)
			}
			return t.Format("2006-01-02")
		}(args[0].(string)), nil
	}

	cases := map[string]struct {
		trial     *Trial
		Case      Case
		expResult result
	}{
		"1/1 - pass case": {
			trial: New(divideFn, nil),
			Case: Case{
				Input:    []interface{}{1, 1},
				Expected: 1,
			},
			expResult: result{true, `PASS: "1/1 - pass case"`},
		},
		"1/0 - error check": {
			trial: New(divideFn, nil),
			Case: Case{
				Input:     []interface{}{1, 0},
				ShouldErr: true,
			},
			expResult: result{true, `PASS: "1/0 - error check"`},
		},
		"1/0 - unexpected error": {
			trial: New(divideFn, nil),
			Case: Case{
				Input: []interface{}{1, 0},
			},
			expResult: result{false, `FAIL: "1/0 - unexpected error" unexpected error divide by zero`},
		},
		"10/2 - unexpected result": {
			trial: New(divideFn, nil),
			Case: Case{
				Input:    []interface{}{10, 2},
				Expected: 10,
			},
			expResult: result{false, "FAIL: \"10/2 - unexpected result\" differences {int}:\n\t-: 5\n\t+: 10\n"},
		},
		"parse time": {
			trial: New(panicFn, nil),
			Case: Case{
				Input:    "2018-01-02T00:00:00Z",
				Expected: "2018-01-02",
			},
			expResult: result{true, `PASS: "parse time"`},
		},
		"parse time with panic": {
			trial: New(panicFn, nil),
			Case: Case{
				Input:       "invalid",
				ShouldPanic: true,
			},
			expResult: result{true, `PASS: "parse time with panic"`},
		},
		"parse time with unexpected panic": {
			trial: New(panicFn, nil),
			Case: Case{
				Input: "invalid",
			},
			expResult: result{false, `PANIC: "parse time with unexpected panic" parsing time "invalid" as "2006-01-02T15:04:05Z07:00": cannot parse "invalid" as "2006"`},
		},
		"expected panic did not occur": {
			trial: New(func(args ...interface{}) (interface{}, error) {
				return nil, nil
			}, nil),
			Case: Case{
				ShouldPanic: true,
			},
			expResult: result{false, `FAIL: "expected panic did not occur" did not panic`},
		},
		"test should error but no error occurred": {
			trial: New(func(args ...interface{}) (interface{}, error) {
				return nil, nil
			}, nil),
			Case: Case{
				ShouldErr: true,
			},
			expResult: result{false, `FAIL: "test should error but no error occurred" should error`},
		},
		"expected error string match": {
			trial: New(func(args ...interface{}) (interface{}, error) {
				return nil, errors.New("test error")
			}, nil),
			Case: Case{
				ExpectedErr: errors.New("test error"),
			},
			expResult: result{true, `PASS: "expected error string match"`},
		},
		"expected error string does not match": {
			trial: New(divideFn, nil),
			Case: Case{
				Input:       Interfaces(10, 0),
				ExpectedErr: errors.New("test error"),
			},
			expResult: result{false, `FAIL: "expected error string does not match" error "divide by zero" does not match expected "test error"`},
		},
	}
	for msg, test := range cases {
		r := test.trial.testCase(msg, test.Case)
		if r.Success != test.expResult.Success || !strings.Contains(r.Message, test.expResult.Message) {
			t.Errorf("FAIL: %q %v", msg, cmp.Diff(r, test.expResult))
		} else {
			t.Logf("PASS: %q", msg)
		}
	}
}