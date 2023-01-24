package smb

import (
	"fmt"
	"strings"

	"github.com/seaung/gagako/pkg/utils"
	"github.com/stacktitan/smb/smb"
)

func smbAuthorized(host, port, username, password string) (ok bool, err error) {
	options := smb.Options{
		Host:        host,
		Port:        445,
		User:        username,
		Password:    password,
		Domain:      "",
		Workstation: "",
	}

	session, err := smb.NewSession(options, false)
	if err != nil {
		return false, err
	}

	if session.IsAuthenticated {
		ok = true
	}

	defer session.Close()
	return ok, err
}

func ScanSmb(target, userdict, passdict string) {
	for _, user := range utils.LoadUserDict(userdict) {
		for _, pass := range utils.LoadPasswordDict(passdict) {
			utils.New().Info(fmt.Sprintf("Scan :: %s :: %s :: %s", target, user, pass))
			ok, err := smbAuthorized(target, "445", user, pass)
			if ok && err == nil {
				utils.New().Success(fmt.Sprintf("FOUND Smb :: %s :: %s :: %s", target, user, pass))
				break
			}
		}
	}
}

func ScanSmb2(target, dictfile string) {
	if utils.IsOpenPort(target, 445) {
		for _, dict := range utils.LoadDicts(dictfile) {
			src := strings.Split(dict, ":")
			user := src[0]
			pass := src[1]
			utils.New().Info(fmt.Sprintf("Scan :: %s :: %s :: %s", target, user, pass))
			ok, err := smbAuthorized(target, "445", user, pass)
			if ok && err == nil {
				utils.New().Success(fmt.Sprintf("FOUND Smb :: %s :: %s :: %s", target, user, pass))
				break
			}
		}
	}
}
