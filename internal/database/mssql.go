package database

import (
	"database/sql"

	_ "github.com/denisenkom/go-mssqldb"

	"fmt"
)

func MssqlUnauthorized(host, port, username, password string) (bool, error) {
	db, err := sql.Open("mssql", fmt.Sprintf("server=%s;user id=%s;password=%s;port=%s;encrypt=disable", host, username, password, port))

	defer db.Close()

	if err != nil {
		return false, err
	}

	if err := db.Ping(); err != nil {
		return false, err
	}

	return true, nil
}
