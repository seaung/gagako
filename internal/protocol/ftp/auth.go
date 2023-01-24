package ftp

import "fmt"

func AuthorizedFTP(host, port, username, password string) (res bool, err error) {
	var ftp *FTP

	if ftp, err = Connect(fmt.Sprintf("%s:%s", host, port)); err != nil {
		return
	}

	defer ftp.Close()

	if err = ftp.Login(username, password); err == nil {
		res = true
	}

	return res, err
}
