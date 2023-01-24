package routeros

import (
	"fmt"
	"strings"

	"github.com/go-routeros/routeros"
	"github.com/seaung/gagako/pkg/utils"
)

func routerosAuthorized(host, port, username, password string) (bool, error) {
	_, err := routeros.Dial(fmt.Sprintf("%s:%s", host, port), username, password)
	if err != nil {
		return false, err
	}

	return true, nil
}

func ScanRouteros(target, userdict, passdict string) {
	for _, user := range utils.LoadUserDict(userdict) {
		for _, pass := range utils.LoadPasswordDict(passdict) {
			utils.New().Info(fmt.Sprintf("Scan %s :  %s : %s :", target, user, pass))
			ok, err := routerosAuthorized(target, "8728", user, pass)
			if ok && err == nil {
				utils.New().Success(fmt.Sprintf("FOUND routeros unauthorized : %s :: %s :: %s\n", target, user, pass))
				break
			}
		}
	}
}

func ScanRouteros2(target, dictfile string) {
	if utils.IsOpenPort(target, 8728) {
		for _, dict := range utils.LoadDicts(dictfile) {
			src := strings.Split(dict, ":")
			user := src[0]
			pass := src[1]
			utils.New().Info(fmt.Sprintf("Scan  %s : %s : %s", target, user, pass))
			ok, err := routerosAuthorized(target, "8728", user, pass)
			if ok && err == nil {
				utils.New().Success(fmt.Sprintf("FOUND routeros unauthorized : %s :: %s :: %s", target, user, pass))
				break
			}
		}
	}
}
