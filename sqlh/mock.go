package sqlh

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"io"

	"github.com/pkg/errors"
)

func init() {
	d := driverMock{}
	sql.Register(Mock, d)
}

type driverMock struct{}

func (d driverMock) Open(dsn string) (driver.Conn, error) {
	m := &mock{}
	if dsn == "open_err" {
		return m, errors.New(dsn)
	}
	return m, nil
}

type mock struct{}

func (m *mock) Prepare(query string) (driver.Stmt, error) {
	return &stmt{}, nil
}

func (m *mock) Close() error {
	return nil
}

func (m *mock) Begin() (driver.Tx, error) {
	return tx{}, nil
}

type stmt struct{}

func (s stmt) Close() error { return nil }

func (s stmt) NumInput() int { return 0 }

func (s stmt) Exec(args []driver.Value) (driver.Result, error) { return result{}, nil }

func (s stmt) ExecContext(ctx context.Context, args []driver.NamedValue) (driver.Result, error) {
	return result{}, nil
}

func (s stmt) QueryContext(ctx context.Context, args []driver.NamedValue) (driver.Rows, error) {
	return rows{}, nil
}

func (s stmt) Query(args []driver.Value) (driver.Rows, error) {
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
