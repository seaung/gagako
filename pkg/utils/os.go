package utils

import (
	"bytes"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

func PrintCurrentOS() {
	categorie := runtime.GOOS

	switch categorie {
	default:
		New().Info("Unkown system !")
	case "linux":
		var out bytes.Buffer
		cmd := exec.Command("/bin/sh", "-c", "uname -a")
		cmd.Stdout = &out
		cmd.Run()
		New().Info(out.String())
	case "windows":
		var out bytes.Buffer
		cmd := exec.Command("cmd", "/c", "ver")
		cmd.Stdout = &out
		cmd.Run()
		New().Info(out.String())

	}
}

func ISPingOK(host string) bool {
	categorie := runtime.GOOS

	var res bool
	switch categorie {
	case "linux":
		var out bytes.Buffer
		cmd := exec.Command("/bin/sh", "-c", fmt.Sprintf("ping -c 1 %s", host))
		cmd.Stdout = &out
		cmd.Run()
		if strings.Contains(out.String(), "ttl=") {
			res = true
		}
	case "windows":
		var out bytes.Buffer
		cmd := exec.Command("cmd", "/c", fmt.Sprintf("ping -n 1 %s", host))
		cmd.Stdout = &out
		cmd.Run()
		if strings.Contains(out.String(), "TTL=") {
			res = true
		}
	}

	return res
}
