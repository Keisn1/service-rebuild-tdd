package userdb_test

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func TestMain(m *testing.M) {
	exitCode := run(m)
	os.Exit(exitCode)
}

func Test_Create(t *testing.T) {
	dsn := fmt.Sprintf(
		"host=localhost port=5432 user=%s password=%s sslmode=disable dbname=%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)
	testDB, err := sql.Open("pgx", dsn)
	if err != nil {
		panic(err)
	}
	defer testDB.Close()

	err = testDB.Ping()
	if err != nil {
		panic(err)
	}
}
