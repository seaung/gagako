package database

import (
	"fmt"
	"time"

	"gopkg.in/mgo.v2"
)

func MongodbUnauthorized(host, port string) (bool, error) {
	session, err := mgo.Dial(fmt.Sprintf("%s:%s", host, port))
	if err != nil && session.Run("serverStatus", nil) != nil {
		return false, err
	}

	return true, nil
}

func MongodbUnauthorizedWithUserAndPasswd(host, port, username, password string) (bool, error) {
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
