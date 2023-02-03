package database

import (
	"database/sql"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/seaung/gagako/pkg/utils"

	"fmt"
)

func mssqlUnauthorized(host, port, username, password string) (bool, error) {
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

func MssalUnauthorizedBureforce(target, userdict, passdict string) {
	if utils.IsOpenPort(target, 1433) {
	Loop:
		for _, user := range utils.LoadUserDict(userdict) {
			for _, pass := range utils.LoadPasswordDict(passdict) {
				ok, err := mssqlUnauthorized(target, "1433", user, pass)
				if ok && err == nil {
					utils.New().Success(fmt.Sprintf("Target :: %s :: username :: %s password :: %s\n", target, user, pass))
					break Loop
				}
			}
		}
	}
}
