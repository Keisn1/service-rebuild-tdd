package userdb_test

import (
	"database/sql"
	"fmt"
	"os"
	"testing"
)

func run(m *testing.M) int {
	var (
		dropDB   = fmt.Sprintf(`DROP DATABASE IF EXISTS %s;`, os.Getenv("DB_NAME"))
		createDB = fmt.Sprintf(`CREATE DATABASE %s;`, os.Getenv("DB_NAME"))
	)

	dsn := fmt.Sprintf(
		"host=localhost port=5432 user=%s password=%s sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
	)
	postgresDB, err := sql.Open("pgx", dsn)
	if err != nil {
		panic(err)
	}
	defer postgresDB.Close()

	_, err = postgresDB.Exec(dropDB)
	if err != nil {
		panic(err)
	}

	_, err = postgresDB.Exec(createDB)
	if err != nil {
		panic(err)
	}

	defer func() {
		_, err = postgresDB.Exec(dropDB)
		if err != nil {
			panic(fmt.Errorf("postgresDB.Exec() err = %s", err))
		}
	}()

	return m.Run()
}
