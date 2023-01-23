package icmp

import (
	"fmt"
	"net"
	"time"

	"github.com/seaung/gagako/pkg/utils"
)

func ICMPStatus(host string) {
}

func isOnline(host string) {
	if isOK(host) {
		utils.New().Info(fmt.Sprintf("the target %s host is UP", host))
		return
	}

	utils.New().Info(fmt.Sprintf("the target %s host is Down", host))
}

func isOK(host string) bool {
	var size int
	var timeout int64
	var seq int16 = 1

	const REQUEST_HEAD_LEN = 8

	size = 32
	timeout = 1000

	startime := time.Now()
	connect, err := net.DialTimeout("ipv4:icmp", host, time.Duration(timeout*1000*1000))
	if err != nil {
		return false
	}

	defer connect.Close()

	id0, id1 := genidentifier(host)

	var message []byte = make([]byte, size+REQUEST_HEAD_LEN)
	message[0] = 8
	message[1] = 0
	message[2] = 0
	message[3] = 0
	message[4], message[5] = id0, id1
	message[6], message[7] = gensequence(seq)

	length := size + REQUEST_HEAD_LEN

	sum := checkSum(message[0:length])
	message[2] = byte(sum >> 8)
	message[3] = byte(sum & 255)

	connect.SetDeadline(startime.Add(time.Duration(timeout * 1000 * 1000)))
	_, _ = connect.Write(message[0:length])

	const REPLY_HEAD_LEN = 20

	var receive []byte = make([]byte, REPLY_HEAD_LEN+length)
	_, err = connect.Read(receive)

	var endtime int = int(int64(time.Since(startime)) / (1000 * 1000))

	if err != nil || receive[REPLY_HEAD_LEN+4] != message[4] || receive[REPLY_HEAD_LEN+5] != message[5] || receive[REPLY_HEAD_LEN+6] != message[6] || receive[REPLY_HEAD_LEN+7] != message[7] || endtime >= int(timeout) || receive[REPLY_HEAD_LEN] == 11 {
		utils.New().Info(fmt.Sprintf("Not Found ICMP Status"))
		return false
	}

	return true
}

func checkSum(message []byte) uint16 {
	sum := 0
	length := len(message)

	for i := 0; i < length-1; i += 2 {
		sum += int(message[i])*256 + int(message[i+1])
	}

	if length%2 == 1 {
		sum += int(message[length-1]) * 256
	}

	sum = (sum >> 16) + (sum & 0xffff)
	sum += (sum >> 16)

	an := uint16(^sum)
	return an
}

func gensequence(val int16) (byte, byte) {
	val1 := byte(val >> 8)
	val2 := byte(val & 255)
	return val1, val2
}

func genidentifier(host string) (byte, byte) {
	return host[0], host[1]
}
