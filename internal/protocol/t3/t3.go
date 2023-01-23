package t3

import (
	"bufio"
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/seaung/gagako/pkg/utils"
)

const (
	number = 255
	load   = 255
)

var (
	wg       sync.WaitGroup
	payload  []byte
	t3Header = []byte{0x74, 0x33, 0x20, 0x31, 0x32, 0x2e, 0x32, 0x2e, 0x31, 0x0a, 0x41, 0x53, 0x3a, 0x32, 0x35, 0x35, 0x0a, 0x48, 0x4c, 0x3a, 0x31, 0x39, 0x0a, 0x4d, 0x53, 0x3a, 0x31, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x0a, 0x50, 0x55, 0x3a, 0x74, 0x33, 0x3a, 0x2f, 0x2f, 0x75, 0x73, 0x2d, 0x6c, 0x2d, 0x62, 0x72, 0x65, 0x65, 0x6e, 0x73, 0x3a, 0x37, 0x30, 0x30, 0x31, 0x0a, 0x0a}
)

type WorkerT3 struct {
	Host string
}

func ConnectWithT3Protocol(host string) (bool, error) {
	conn, err := net.DialTimeout("tcp", host, 5*time.Second)
	if err != nil {
		return false, err
	}

	_, err = conn.Write(t3Header)
	if err != nil {
		return false, err
	}

	ReadAllFromT3(conn, host)

	return true, nil
}

func ReadAllFromT3(connect net.Conn, host string) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(2*time.Second))
	ch := make(chan struct{})
	buffer := make([]byte, 1024)
	reader := bufio.NewReader(connect)

	go func(ch chan struct{}) {
		for {
			nbytes, err := reader.Read(buffer)
			if err != nil {
				break
			}

			res := string(buffer[:nbytes])
			verbose := GetString(res, ":", ".false")
			if strings.Contains(res, ".false") {
				utils.New().Info(fmt.Sprintf("%s \tWeblogic%s\n", host, verbose))
			}

		}

		ch <- struct{}{}
	}(ch)

	select {
	case <-ch:
	case <-ctx.Done():
		break
	}

	cancel()
}

func GetString(src, start, end string) string {
	nbytes := strings.Index(src, start)
	if nbytes == -1 {
		nbytes = 0
	}

	src = string([]byte(src)[nbytes:])
	mstr := strings.Index(src, end)
	if mstr == -1 {
		mstr = len(src)
	}

	src = string([]byte(src)[:mstr])
	return src
}

func PrintT3ProtoclVersion(host string) {
	ConnectWithT3Protocol(fmt.Sprintf("%s:7100", host))
	ConnectWithT3Protocol(fmt.Sprintf("%s:7002", host))
	ConnectWithT3Protocol(fmt.Sprintf("%s", host))
	ConnectWithT3Protocol(fmt.Sprintf("%s:8080", host))
}

func Int8ToByte(val int) []byte {
	var x uint16
	x = uint16(val)

	buffer := bytes.NewBuffer([]byte{})
	binary.Write(buffer, binary.BigEndian, x)
	return buffer.Bytes()
}

func Int32ToByte(val int) []byte {
	var x uint32

	x = uint32(val)
	buffer := bytes.NewBuffer([]byte{})
	binary.Write(buffer, binary.BigEndian, x)
	return buffer.Bytes()
}

func AddByte(buffer *[]byte, buffs ...[]byte) {
	for _, nbytes := range buffs {
		*buffer = append(*buffer, nbytes...)
	}
}

func RunTaskWithT3Protocol(host string) {
	tasks := make(chan WorkerT3, load)
	wg.Add(number)

	for num := 1; num <= number; num++ {
		go workerWithT3Protocol(tasks)
	}

	for i := 0; i < 256; i++ {
		target := fmt.Sprintf("%s.%d", host, i)
		task := WorkerT3{
			Host: target,
		}
		tasks <- task
	}

	close(tasks)
	wg.Wait()
}

func workerWithT3Protocol(tasks chan WorkerT3) {
	defer wg.Done()

	task, ok := <-tasks
	if !ok {
		return
	}

	host := task.Host

	PrintT3ProtoclVersion(host)
}
