package postgres_test

import (
	"database/sql"
	"errors"
)

type SQLDB struct{}

func (s *SQLDB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return nil, errors.New("DBError")
}

func (s *SQLDB) QueryRow(query string, args ...any) (row *sql.Row) {
	return
}

func (s *SQLDB) Exec(query string, args ...any) (res sql.Result, err error) {
	return
}
