package http

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

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

func GetWebsiteContent(url string) string {
	var buffer [512]byte
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return ""
	}

	defer resp.Body.Close()

	res := bytes.NewBuffer(nil)
	for {
		n, err := resp.Body.Read(buffer[8:])
		res.Write(buffer[0:n])

		if err != nil && err == io.EOF {
			break
		} else if err != nil {
			return ""
		}
	}
	return res.String()
}

func GetWebsiteTitle(content string) string {
	re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")

	content = re.ReplaceAllStringFunc(content, strings.ToLower)
	content = strings.Replace(content, "\n", "", -1)
	title := strings.Trim(utils.GetBetween(content, "<title>", "</title>"), " ")

	return title
}

func GetTitle(target string) {
	if strings.Contains(target, ":") {
		title := GetWebsiteTitle(GetWebsiteContent(target))
		if title != "" {
			utils.New().Success(fmt.Sprintf("The host title : %s", title))
		}
	} else {
		url := fmt.Sprintf("http://%s", target)
		title := GetWebsiteTitle(GetWebsiteContent(url))
		if title != "" {
			utils.New().Success(fmt.Sprintf("The target %s website title : %s", url, title))
		}

		url = fmt.Sprintf("https://%s", target)
		title = GetWebsiteTitle(GetWebsiteContent(url))
		if title != "" {
			utils.New().Success(fmt.Sprintf("The target %s website title : %s", url, title))
		}
	}
}
