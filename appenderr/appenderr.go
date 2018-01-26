package appenderr

import (
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
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
