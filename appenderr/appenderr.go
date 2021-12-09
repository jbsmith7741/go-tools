package appenderr

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
)

// AppendErr a helper error struct that counts and removes duplicates errors.
// errors are considered identical if their error signature is the same.
type AppendErr struct {
	mu   sync.RWMutex
	errs map[string]errData
}

type errData struct {
	err   error
	count int
	t     time.Time
}

type timeErrs []errData

func (a timeErrs) Len() int      { return len(a) }
func (a timeErrs) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a timeErrs) Less(i, j int) bool {
	return a[i].t.Before(a[j].t)
}

// New AppendErr
func New() *AppendErr {
	e := &AppendErr{
		errs: make(map[string]errData, 0),
	}
	return e
}

// Is uses errors.Is for each recorded error
func (e *AppendErr) Is(err error) bool {
	for _, v := range e.errs {
		if errors.Is(v.err, e) {
			return true
		}
	}
	return false
}

// As uses errors.As for each recorded error
func (e *AppendErr) As(err interface{}) bool {
	for _, v := range e.errs {
		if errors.As(v.err, e) {
			return true
		}
	}
	return false
}

// Add an error to stack, appendErr will ignore nil errors
func (e *AppendErr) Add(err error) {
	if err == nil {
		return
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	if appendErr, ok := err.(*AppendErr); ok {
		for key, aErr := range appendErr.errs {
			if d, found := e.errs[key]; found {
				d.count += aErr.count
				e.errs[key] = d
			} else {
				e.errs[key] = aErr
			}
		}
		return
	}
	key := err.Error()
	if d, found := e.errs[key]; found {
		d.count++
		e.errs[key] = d
	} else {
		e.errs[key] = errData{
			err:   err,
			count: 1,
			t:     time.Now(),
		}
	}
	return
}

// Addf creates an error using fmt.Errorf and adds it to the error list
func (e *AppendErr) Addf(format string, args ...interface{}) {
	err := fmt.Errorf(format, args...)
	e.Add(err)
}

// Wrap an error message and add it to the error list.
// see errors.Wrap
func (e *AppendErr) Wrap(err error, message string) {
	if err != nil {
		e.Add(errors.Wrap(err, message))
	}
}

// Wrapf wrap an error message with a formated message
// see errors.Wrapf
func (e *AppendErr) Wrapf(err error, format string, args ...interface{}) {
	if err != nil {
		e.Add(errors.Wrapf(err, format, args...))
	}
}

func (e *AppendErr) Error() (s string) {

	errs := make([]errData, 0, len(e.errs))
	e.mu.RLock()
	for _, err := range e.errs {
		errs = append(errs, err)
	}
	e.mu.RUnlock()

	sort.Sort(timeErrs(errs))

	for _, err := range errs {
		if err.count > 1 {
			s += "[" + strconv.Itoa(err.count) + "] "
		}
		s += err.err.Error() + "\n"
	}
	return strings.TrimRight(s, "\n")
}

// ErrOrNil returns itself if errors are found or nil
func (e *AppendErr) ErrOrNil() error {
	if len(e.errs) == 0 {
		return nil
	}
	return e
}
