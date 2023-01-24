package http

import (
	"fmt"
	"net/http"

	"github.com/seaung/gagako/pkg/utils"
)

func ISBasicAuth(target string) bool {
	client := &http.Client{}
	request, _ := http.NewRequest(http.MethodGet, target, nil)
	request.SetBasicAuth("", "")

	resp, err := client.Do(request)
	if err != nil {
		return false
	}

	if resp.StatusCode == 401 {
		return true
	}

	defer resp.Body.Close()

	return false
}

func BasicAuth(url, username, password string) (bool, error) {
	client := &http.Client{}
	request, _ := http.NewRequest(http.MethodGet, url, nil)
	request.SetBasicAuth(username, password)

	resp, err := client.Do(request)
	if err != nil {
		return false, err
	}

	if resp.StatusCode == 401 {
		return true, nil
	}

	defer resp.Body.Close()

	return false, err
}

func ScanBasicAuth(userdict, passdict, target string) {
	if ISBasicAuth(target) {
		for _, user := range utils.LoadUserDict(userdict) {
			for _, pass := range utils.LoadPasswordDict(passdict) {
				res, err := BasicAuth(target, user, pass)
				if res == true && err == nil {
					utils.New().Success(fmt.Sprintf("FOUND Basic UnAuthorized : host - %s user - %s pass - %s", target, user, pass))
					break
				}
			}
		}
	}
}
