package sqlh

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"io"
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

func init() {
	d := make(driverPool)
	sql.Register(Mock, d)
}

func NewMock(dsn string) (*sql.DB, *mock) {
	db, _ := sql.Open(Mock, "")
	conn, _ := db.Driver().Open(dsn)
	return db, conn.(*mock)
}

var regPGInput = regexp.MustCompile("([$][0-9]*)")

type driverPool map[string]*mock

func (d driverPool) Open(dsn string) (driver.Conn, error) {
	if m, found := d[dsn]; found {
		return m, nil
	}

	d[dsn] = &mock{}
	return d[dsn], nil
}

type mock struct {
	options string
}

func (m *mock) SetOptions(o string) *mock {
	m.options = o
	return m
}

func (m *mock) Prepare(query string) (driver.Stmt, error) {
	matches := regPGInput.FindAllString(query, -1)
	count := strings.Count(query, "?")
	if len(matches) > count {
		max := 0
		for _, v := range matches {
			i, _ := strconv.Atoi(strings.TrimLeft(v, "$"))
			if i > max {
				max = i
			}
		}
		count = max
	}

	if m.options == "prepare_err" {
		return nil, errors.New(m.options)
	}

	return &stmt{
		options:  m.options,
		numInput: count,
	}, nil
}

func (m *mock) Close() error {
	return nil
}

func (m *mock) Begin() (driver.Tx, error) {
	return tx{}, nil
}

type stmt struct {
	options  string
	numInput int
}

func (s stmt) Close() error { return nil }

func (s stmt) NumInput() int { return s.numInput }

func (s stmt) Exec(args []driver.Value) (driver.Result, error) { return result{}, nil }

func (s stmt) ExecContext(ctx context.Context, args []driver.NamedValue) (driver.Result, error) {
	return result{}, nil
}

func (s stmt) QueryContext(ctx context.Context, args []driver.NamedValue) (driver.Rows, error) {
	if s.options == "query_err" {
		return nil, errors.New(s.options)
	}
	return rows{}, nil
}

func (s stmt) Query(args []driver.Value) (driver.Rows, error) {
	if s.options == "query_err" {
		return nil, errors.New(s.options)
	}
	return rows{}, nil
}

type result struct{}

func (r result) LastInsertId() (int64, error) { return 0, nil }
func (r result) RowsAffected() (int64, error) { return 0, nil }

type rows struct{}

func (r rows) Columns() []string { return []string{} }

func (r rows) Close() error { return nil }

func (r rows) Next(dest []driver.Value) error { return io.EOF }

type tx struct{}

func (tx tx) Commit() error { return nil }

func (tx tx) Rollback() error { return nil }
