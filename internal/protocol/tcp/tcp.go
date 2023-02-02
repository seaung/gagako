package tcp

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/seaung/gagako/pkg/utils"
)

func isOpen(host string, port int) bool {
	connect, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", host, port), time.Second*5)
	if err != nil {
		utils.New().LoggerError(fmt.Sprintf("the port state is closed : %d", port))
		return false
	}

	utils.New().Info(fmt.Sprintf("the port state is open : %d", port))

	defer connect.Close()

	return true
}

func isSSHService(address string) (string, error) {
	connect, err := net.DialTimeout("tcp", address, time.Second*10)
	if err != nil {
		return "", err
	}

	defer connect.Close()

	tcpConn := connect.(*net.TCPConn)
	tcpConn.SetReadDeadline(time.Now().Add(time.Second * 5))
	reader := bufio.NewReader(connect)

	return reader.ReadString('\n')
}

func splitHttpHeader(data []byte, at bool) (int, []byte, error) {
	end := bytes.Index(data, []byte("\r\n\r\n"))
	if end == -1 {
		return 0, []byte(""), fmt.Errorf("split http header error : %d", end)
	}

	return end + 4, data[:end+4], nil
}

func isHTTPService(address string) (string, error) {
	connect, err := net.DialTimeout("tcp", address, time.Second*10)
	if err != nil {
		return "", err
	}

	defer connect.Close()

	tcpConn := connect.(*net.TCPConn)
	tcpConn.SetWriteDeadline(time.Now().Add(time.Second * 5))
	scanner := bufio.NewScanner(connect)
	scanner.Split(splitHttpHeader)

	if scanner.Scan() {
		return scanner.Text(), nil
	}

	err = scanner.Err()
	if err == nil {
		err = io.EOF
	}

	return "", err
}

func GetBanner(address string, port int) {
	if isOpen(address, port) {
		var wg sync.WaitGroup
		result := make(chan string, 2)
		done := make(chan int, 1)

		wg.Add(2)

		go func() {
			if r, err := isSSHService(fmt.Sprintf("%s:%d", address, port)); err == nil {
				result <- fmt.Sprintf("SSH Banner : %s", r)
			}

			wg.Done()
		}()

		go func() {
			if r, err := isHTTPService(fmt.Sprintf("%s:%d", address, port)); err == nil {
				result <- fmt.Sprintf("HTTP Banner : %s", r)
			}

			wg.Done()
		}()

		go func() {
			wg.Wait()
			done <- 1
		}()

		select {
		case <-done:
		case r := <-result:
			fmt.Sprintf("%s:%d\t%s", address, port, r)
		}
	}
}

func GetTCPBanner(host string, port int) {
	var banner string
	buufer := make([]byte, 1024)

	connect, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", host, port), time.Second*5)
	if err != nil {
		utils.New().Warnning(fmt.Sprintf("connet to %s:%d timeout : error info : %v\n", host, port, err))
		return
	}

	defer connect.Close()

	fmt.Fprintf(connect, "\r\n\r\n")
	connect.SetReadDeadline(time.Now().Add(time.Second * 8))

	n, _ := connect.Read(buufer)

	banner = string(buufer[:n])

	if banner == "" {
		fmt.Fprintf(connect, "GET / HTTP/1.1\r\n\r\n")
		connect.SetReadDeadline(time.Now().Add(time.Second * 8))
		rbuffer := make([]byte, 1024)
		n, _ := connect.Read(rbuffer)

		banner = string(rbuffer[:n])
	}

	banner = strings.Replace(banner, "\r\n", " ", -1)

	if strings.Contains(banner, "SSH-") {
		utils.New().Success(fmt.Sprintf("%s:%d is open \tbanner: %s\n", host, port, banner))
	} else if strings.Contains(banner, "HTTP/1") {
		utils.New().Success(fmt.Sprintf("%s:%d is open \tbanner: %s\n", host, port, banner))
	} else if strings.Contains(banner, "FTP") {
		utils.New().Success(fmt.Sprintf("%s:%d is open \tbanner: %s\n", host, port, banner))
	} else {
		utils.New().Info(fmt.Sprintf("%s:%d is open \tbanner: %s\n", host, port, banner))
	}
}
