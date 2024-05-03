package notedb_test

import (
	"context"
	"database/sql"
	"errors"
)

type stubSQLDB struct{}

func (s *stubSQLDB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return nil, errors.New("DBError")
}

func (s *stubSQLDB) QueryRowContext(ctx context.Context, query string, args ...any) (row *sql.Row) {
	return
}

func (s *stubSQLDB) Exec(query string, args ...any) (res sql.Result, err error) {
	return nil, errors.New("DBError")
}
