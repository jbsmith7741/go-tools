package appenderr

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestAdd(t *testing.T) {
	cases := []struct {
		msg      string
		errs     []error
		expected string
		isNil    bool
	}{
		{
			msg:   "No Errors",
			isNil: true,
		},
		{
			msg:      "3 identical errors",
			errs:     makeErrs("err", "err", "err"),
			expected: "[3] err",
		},
		{
			msg:      "3 unique errors",
			errs:     makeErrs("err3", "err2", "err1"),
			expected: "err3\nerr2\nerr1",
		},
		{
			msg: "add append errors",
			errs: append(makeErrs("err1", "err2"), &AppendErr{
				errs: map[string]errData{
					"err1": errData{
						err:   fmt.Errorf("err1"),
						count: 1,
						t:     time.Now().Add(time.Second),
					},
					"err3": errData{
						err:   fmt.Errorf("err3"),
						count: 2,
						t:     time.Now().Add(time.Second),
					},
				},
			}),
			expected: "[2] err1\nerr2\n[2] err3",
		},
		{
			msg:      "errors with one nil",
			errs:     []error{(error)(nil), fmt.Errorf("err1")},
			expected: "err1",
		},
		{
			msg:   "nil error",
			errs:  []error{(error)(nil)},
			isNil: true,
		},
	}
	for _, test := range cases {
		aErr := New()
		for _, err := range test.errs {
			aErr.Add(err)
		}
		if test.isNil != (aErr.ErrOrNil() == nil) {
			t.Errorf("FAIL: %s %v", test.msg, aErr.ErrOrNil())
		} else if !cmp.Equal(aErr.Error(), test.expected) {
			t.Errorf("FAIL: %s %s", test.msg, cmp.Diff((aErr.Error()), test.expected))
		} else {
			t.Logf("PASS: %s", test.msg)
		}
	}
}

func makeErrs(data ...string) (errs []error) {
	for _, v := range data {
		errs = append(errs, fmt.Errorf(v))
	}
	return errs
}
