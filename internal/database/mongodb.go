package database

import (
	"fmt"
	"time"

	"github.com/seaung/gagako/pkg/utils"
	"gopkg.in/mgo.v2"
)

func mongodbUnauthorized(host, port string) (bool, error) {
	session, err := mgo.Dial(fmt.Sprintf("%s:%s", host, port))
	if err != nil && session.Run("serverStatus", nil) != nil {
		return false, err
	}

	return true, nil
}

func mongodbUnauthorizedWithUserAndPasswd(host, port, username, password string) (bool, error) {
	session, err := mgo.DialWithTimeout(fmt.Sprintf("mongodb://%s:%s@%s:%s/admin", username, password, host, port), time.Second*3)
	defer session.Close()

	if err != nil && session.Ping() != nil {
		return false, err
	}

	if session.Run("serverStatus", nil) != nil {
		return false, err
	}

	return true, nil
}

func ScanMogondb(target, userdict, passdict string) {
	if utils.IsOpenPort(target, 27017) {
	Loop:
		for _, user := range utils.LoadUserDict(userdict) {
			for _, pass := range utils.LoadPasswordDict(passdict) {
				ok, err := mongodbUnauthorizedWithUserAndPasswd(target, "27017", user, pass)
				if ok && err == nil {
					utils.New().Success(fmt.Sprintf("Target :: %s username :: %s password :: %s", target, user, pass))
					break Loop
				}
			}
		}
	}
}
