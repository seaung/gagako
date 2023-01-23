package utils

import (
	"bytes"
	"crypto/md5"
	"io"
)

type User struct {
	Username string
	Password string
}

func GetUserAndDecryptPassword(data []byte) ([]User, error) {
	var res []User

	dt := bytes.Split(data, []byte("M2"))[1:]

	for _, val := range dt {
		name, pass, err := extractorUserAndPasswd(val)
		if err != nil {
			continue
		}

		password, err := decryptPassword(name, pass)
		if err != nil {
			return nil, err
		}

		res = append(res, User{
			Username: string(name),
			Password: string(password),
		})
	}
	return res, nil
}

func decryptPassword(name, encrypt []byte) ([]byte, error) {
	var pass []byte
	key := []byte("283i4jfkai3389")
	data := md5.New()
	if _, err := io.WriteString(data, string(name)+string(key)); err != nil {
		return nil, err
	}

	digitKey := data.Sum(nil)
	for i := range encrypt {
		pass = append(pass, encrypt[i]^digitKey[i%len(digitKey)])
	}

	return pass, nil
}

func extractorUserAndPasswd(data []byte) (name, pass []byte, err error) {
	username := bytes.Split(data, []byte("\x01\x00\x00\x21"))
	userpass := bytes.Split(data, []byte("\x11\x00\x00\x21"))
	if len(username) != 1 && len(userpass) != 1 {
		usernameLen := username[1][0]
		userpassLen := userpass[1][0]
		name = username[1][1 : 1+int(usernameLen)]
		pass = userpass[1][1 : 1+int(userpassLen)]
		return name, pass, nil
	}
	return nil, nil, err
}
