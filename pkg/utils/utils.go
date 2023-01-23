package utils

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func IsOpenPort(ipaddr string, port int) bool {
	addr := net.TCPAddr{
		IP:   net.ParseIP(ipaddr),
		Port: port,
	}

	connect, err := net.DialTCP("tcp", nil, &addr)
	if err != nil {
		New().Info(fmt.Sprintf("the connect is close, error : %v", err))
		return false
	}

	if connect != nil {
		connect.Close()
		New().Info(fmt.Sprintf("the %d port is Open", port))
		return true
	}

	return false
}

func LoadDicts(filename string) []string {
	dicts := []string{}
	file, err := os.Open(filename)
	if err != nil {
		New().Warnning(fmt.Sprintf("Open %s error %v", filename, err))
		return dicts
	}

	fd, _ := os.Stat(filename)
	if fd.Size() == 0 {
		New().Info(fmt.Sprintf("The file %s is NULL ", filename))
		return dicts
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			dicts = append(dicts, line)
		}
	}

	return dicts
}

func LoadUserDict(filename string) []string {
	dicts := []string{}
	file, err := os.Open(filename)
	if err != nil {
		New().Warnning(fmt.Sprintf("Open %s error %v", filename, err))
		return dicts
	}

	fd, _ := os.Stat(filename)

	if fd.Size() == 0 {
		New().Info(fmt.Sprintf("The file %s is NULL ", filename))
		return dicts
	}

	defer file.Close()

	reader := bufio.NewScanner(file)
	reader.Split(bufio.ScanLines)

	for reader.Scan() {
		user := strings.TrimSpace(reader.Text())
		if user != "" {
			dicts = append(dicts, user)
		}
	}

	return dicts
}

func LoadPasswordDict(filename string) []string {
	dicts := []string{}
	file, err := os.Open(filename)
	if err != nil {
		New().LoggerError(fmt.Sprintf("Open %s error : %v", filename, err))
		return dicts
	}

	fd, _ := os.Stat(filename)
	if fd.Size() == 0 {
		New().Info(fmt.Sprintf("The file %s is NULL ", filename))
		return dicts
	}

	defer file.Close()

	reader := bufio.NewScanner(file)
	reader.Split(bufio.ScanLines)

	for reader.Scan() {
		pwd := strings.TrimSpace(reader.Text())
		if pwd != "" {
			dicts = append(dicts, pwd)
		}
	}

	return dicts
}
