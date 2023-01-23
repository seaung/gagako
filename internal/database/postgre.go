package database

import (
	"database/sql"

	"fmt"

	_ "github.com/lib/pq"
)

func ConnectPostgresql(host, port, username, password string) (bool, error) {
	db, err := sql.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s:%s/postgres?sslmode=verify-full", username, password, host, port))
	defer db.Close()
	if err != nil {
		return false, err
	}

	if err := db.Ping(); err != nil {
		return false, err
	}

	return true, nil
}
