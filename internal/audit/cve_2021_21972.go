package audit

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"github.com/seaung/gagako/pkg/utils"
)

func CVE202121972(target string) {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: false,
		},
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   time.Duration(3 * time.Second),
	}

	url := fmt.Sprintf("%s%s", target, "/ui/vropspluginui/rest/services/uploadova")

	request, _ := http.NewRequest(http.MethodGet, url, nil)
	request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/534.50 (KHTML, like Gecko) Version/5.1 Safari/534.50")
	resp, _ := client.Do(request)

	if resp.StatusCode == 405 {
		utils.New().Success(fmt.Sprintf("Found CVE-2021-21972 : target %s", target))
	}
}
