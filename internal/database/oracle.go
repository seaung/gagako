package database

import (
	"database/sql"
	"fmt"
)

func OracleUnauthorized(host, port, username, password string) (bool, error) {
	db, err := sql.Open("godror", fmt.Sprintf("%s/%s@%s:%s/orcl", username, password, host, port))
	defer db.Close()

	if err != nil {
		return false, err
	}

	if err := db.Ping(); err != nil {
		return false, err
	}

	return true, nil
}
