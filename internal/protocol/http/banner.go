package http

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/seaung/gagako/pkg/utils"
)

func ISURL(url string) string {
	if !strings.Contains(url, "http") {
		return fmt.Sprintf("http://%s", url)
	}

	return url
}

func GetHTTPBanner(url string) (bool, error) {
	target := ISURL(url)

	response, err := http.Head(target)
	if err != nil {
		return false, err
	}

	for key, val := range response.Header {
		if key == "Server" {
			utils.New().Info(fmt.Sprintf("Target : %s - Header value : %s", target, val))
		}
	}

	return true, nil
}
