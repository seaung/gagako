package database

import (
	"database/sql"

	"fmt"

	_ "github.com/lib/pq"
	"github.com/seaung/gagako/pkg/utils"
)

func connectPostgresql(host, port, username, password string) (bool, error) {
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

func ScanPostgresql(target, userdict, passdict string) {
	if utils.IsOpenPort(target, 5432) {
	Loop:
		for _, user := range utils.LoadUserDict(userdict) {
			for _, pass := range utils.LoadPasswordDict(passdict) {
				ok, err := connectPostgresql(target, "5432", user, pass)
				if ok && err == nil {
					utils.New().Success(fmt.Sprintf("Target :: %s :: username :: %s password :: %s", target, user, pass))
					break Loop
				}
			}
		}
	}
}
