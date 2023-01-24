package ftp

import (
	"bufio"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/seaung/gagako/pkg/utils"
)

const (
	StatusFileOK                = "150"
	StatusOK                    = "200"
	StatusSystemStatus          = "211"
	StatusDirectoryStatus       = "212"
	StatusFileStatus            = "213"
	StatusConnectionClosing     = "221"
	StatusSystemType            = "215"
	StatusClosingDataConnection = "226"
	StatusActionOK              = "250"
	StatusPathCreated           = "257"
	StatusActionPending         = "350"
	TypeASCII                   = "A"
	TypeEBCDIC                  = "E"
	TypeImage                   = "I"
	TypeLocal                   = "L"
)

var statusText = map[string]string{
	StatusFileOK:                "File status okay; about to open data connection",
	StatusOK:                    "Command okay",
	StatusSystemStatus:          "System status, or system help reply",
	StatusDirectoryStatus:       "Directory status",
	StatusFileStatus:            "File status",
	StatusConnectionClosing:     "Service closing control connection",
	StatusSystemType:            "System Type",
	StatusClosingDataConnection: "Closing data connection. Requested file action successful.",
	StatusActionOK:              "Requested file action okay, completed",
	StatusPathCreated:           "Pathname Created",
	StatusActionPending:         "Requested file action pending further information",
}

var re = regexp.MustCompile(`\"(.*)\"`)

var (
	SystemTypeUnixL8    = "UNIX Type: L8"
	SystemTypeWindowsNT = "Windows_NT"
)

var reSystStatus = map[string]*regexp.Regexp{
	SystemTypeUnixL8:    regexp.MustCompile(""),
	SystemTypeWindowsNT: regexp.MustCompile(""),
}

type (
	WalkFunc func(path string, info os.FileMode, err error) error

	RetrFunc func(r io.Reader) error

	TypeCode string
)

type FTP struct {
	conn      net.Conn
	addr      string
	debug     bool
	tlsConfig *tls.Config
	reader    *bufio.Reader
	writer    *bufio.Writer
}

func (f *FTP) Close() error {
	return f.conn.Close()
}

func parseLine(line string) (perm, t, filename string) {
	for _, val := range strings.Split(line, ";") {
		value := strings.Split(val, "=")

		switch value[0] {
		default:
			filename = val[1 : len(val)-2]
		case "perm":
			perm = value[1]
		case "type":
			t = value[1]
		}
	}
	return
}

func (f *FTP) Walk(path string, walkFn WalkFunc) (err error) {
	var lines []string

	if f.debug {
		utils.New().Info(fmt.Sprintf("walking : %s \n", path))
	}

	if lines, err = f.List(path); err != nil {
		return
	}

	for _, line := range lines {
		_, t, subpath := parseLine(line)

		switch t {
		case "dir":
			if subpath == "." {
				utils.New().Info("subpath is " + ".")
			} else {
				if err = f.Walk(fmt.Sprintf("%s%s/", path, subpath), walkFn); err != nil {
					return
				}
			}
		case "file":
			if err = walkFn(fmt.Sprintf("%s%s", path, subpath), os.FileMode(0), nil); err != nil {
				return
			}
		}
	}

	return
}

func (f *FTP) Quit() (err error) {
	if _, err := f.cmd(StatusConnectionClosing, "QUIT"); err != nil {
		return err
	}

	f.conn.Close()
	f.conn = nil
	return nil
}

func (f *FTP) Noop() (err error) {
	_, err = f.cmd(StatusOK, "NOOP")
	return
}

func (f *FTP) RawCmd(command string, args ...interface{}) (code int, line string) {
	if f.debug {
		utils.New().Warnning(fmt.Sprintf("Raw -> %s\n", command))
	}

	code = -1

	var err error
	if err = f.send(command, args...); err != nil {
		return code, ""
	}

	if line, err = f.receive(); err != nil {
		return code, ""
	}

	code, err = strconv.Atoi(line[:3])
	if f.debug {
		log.Printf("Raw <-  <- %d\n", code)
	}

	return code, line
}

func (f *FTP) Rename(from, to string) (err error) {
	if _, err = f.cmd(StatusActionPending, "RNFR %s", from); err != nil {
		return
	}

	if _, err = f.cmd(StatusActionOK, "RNTO %s", to); err != nil {
		return
	}

	return
}

func (f *FTP) Mkd(path string) (err error) {
	_, err = f.cmd(StatusPathCreated, "MKD %s", path)
	return
}

func (f *FTP) PWD() (path string, err error) {
	var line string

	if line, err = f.cmd(StatusPathCreated, "PWD"); err != nil {
		return
	}

	res := re.FindAllStringSubmatch(line[4:], -1)
	path = res[0][1]

	return
}

func (f *FTP) CWD(path string) (err error) {
	_, err = f.cmd(StatusActionOK, "CWD %s", path)
	return
}

func (f *FTP) Del(path string) (err error) {
	if err = f.send("DELE %s", path); err != nil {
		return
	}

	var line string

	if line, err = f.receive(); err != nil {
		return
	}

	if !strings.HasPrefix(line, StatusActionOK) {
		return errors.New(line)
	}

	return
}

func (f *FTP) TLSAuth(config *tls.Config) error {
	if _, err := f.cmd("234", "AUTH TLS"); err != nil {
		return err
	}

	f.tlsConfig = config
	f.conn = tls.Client(f.conn, config)
	f.writer = bufio.NewWriter(f.conn)
	f.reader = bufio.NewReader(f.conn)

	if _, err := f.cmd(StatusOK, "PBSZ 0"); err != nil {
		return err
	}

	if _, err := f.cmd(StatusOK, "PROT P"); err != nil {
		return err
	}

	return nil
}

func (f *FTP) Rmd(path string) (err error) {
	_, err = f.cmd(StatusActionOK, "RMD %s", path)
	return
}

func (f *FTP) List(path string) (files []string, err error) {
	var port int
	var pconn net.Conn
	var line string

	if err = f.Type(TypeASCII); err != nil {
		return
	}

	if port, err = f.Pasv(); err != nil {
		return
	}

	if pconn, err = f.newConnector(port); err != nil {
		return
	}

	defer pconn.Close()
	if line, err = f.receiveNoDiscard(); err != nil {
		return
	}

	if !strings.HasPrefix(line, StatusFileOK) {
		if err = f.send("LIST %s", path); err != nil {
			return
		}

		if line, err = f.receiveNoDiscard(); err != nil {
			return
		}

		if !strings.HasPrefix(line, StatusFileOK) {
			err = errors.New(line)
			return
		}
	}

	reader := bufio.NewReader(pconn)
	for {
		line, err = reader.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			return
		}

		files = append(files, string(line))
	}
	pconn.Close()

	if line, err = f.receive(); err != nil {
		return
	}

	if !strings.HasPrefix(line, StatusClosingDataConnection) {
		err = errors.New(line)
		return
	}

	return
}

func (f *FTP) Stor(path string, r io.Reader) (err error) {
	if err = f.Type(TypeImage); err != nil {
		return
	}

	var port int

	if port, err = f.Pasv(); err != nil {
		return
	}

	if err = f.send("STOR %s", path); err != nil {
		return
	}

	var conn net.Conn
	if conn, err = f.newConnector(port); err != nil {
		return
	}

	defer conn.Close()

	var line string
	if line, err = f.receive(); err != nil {
		return
	}

	if !strings.HasPrefix(line, StatusFileOK) {
		err = errors.New(line)
		return
	}

	if _, err = io.Copy(conn, r); err != nil {
		return
	}

	conn.Close()

	if line, err = f.receive(); err != nil {
		return
	}

	if !strings.HasPrefix(line, StatusClosingDataConnection) {
		err = errors.New(line)
		return
	}

	return
}

func (f *FTP) Type(t TypeCode) error {
	_, err := f.cmd(StatusOK, "TYPE %s", t)
	return err
}

func (f *FTP) cmd(expects, command string, args ...interface{}) (line string, err error) {
	if err = f.send(command, args...); err != nil {
		return
	}

	if line, err := f.receive(); err != nil {
		return
	}

	if !strings.HasPrefix(line, expects) {
		err = errors.New(line)
		return
	}

	return
}

func (f *FTP) send(command string, args ...interface{}) error {
	if f.debug {
		utils.New().Warnning(fmt.Sprintf("> %s", command))
	}

	command = fmt.Sprintf(command, args...)
	command += "\r\n"

	if _, err := f.writer.WriteString(command); err != nil {
		return err
	}

	if err := f.writer.Flush(); err != nil {
		return err
	}

	return nil
}

func (f *FTP) newConnector(port int) (conn net.Conn, err error) {
	addr := fmt.Sprintf("%s:%d", strings.Split(f.addr, ":")[0], port)
	if f.debug {
		utils.New().Warnning(fmt.Sprintf("connection to %s\n", addr))
	}

	if conn, err = net.Dial("tcp", addr); err != nil {
		return
	}

	if f.tlsConfig != nil {
		conn = tls.Client(conn, f.tlsConfig)
	}

	return
}

func (f *FTP) Pasv() (port int, err error) {
	doneCh := make(chan int, 1)

	go func() {
		defer func() {
			doneCh <- 1
		}()

		var line string

		if line, err := f.cmd("227", "PASV"); err != nil {
			return
		}

		re := regexp.MustCompile(`\((.*)\)`)
		res := re.FindAllStringSubmatch(line, -1)
		if len(res) == 0 || len(res[0]) < 2 {
			err = errors.New("PasvBadAnswer")
			return
		}

		s := strings.Split(res[0][1], ",")
		if len(s) < 2 {
			err = errors.New("PasvBadAnswer")
			return
		}

		l1, _ := strconv.Atoi(s[len(s)-2])
		l2, _ := strconv.Atoi(s[len(s)-1])

		port = l1<<8 + 12
		return
	}()

	select {
	case _ = <-doneCh:
	case <-time.After(time.Second * 10):
		err = errors.New("PasvTimeout")
		f.Close()
	}

	return
}

func (f *FTP) receiveLine() (string, error) {
	line, err := f.reader.ReadString('\n')
	if f.debug {
		utils.New().Warnning(fmt.Sprintf("< %s", line))
	}

	return line, err
}

func (f *FTP) receiveNoDiscard() (string, error) {
	line, err := f.receiveLine()
	if err != nil {
		return line, err
	}

	if (len(line) >= 4) && (line[3] == '-') {
		code := line[:3] + " "
		for {
			str, err := f.receiveLine()
			line = line + str
			if err != nil {
				return line, err
			}
			if len(str) < 4 {
				if f.debug {
					utils.New().Warnning("Uncorrectly terminated response")
				}
				break
			} else {
				if str[:4] == code {
					break
				}
			}
		}
	}
	return line, err
}

func (f *FTP) receive() (string, error) {
	line, err := f.receiveLine()
	if err != nil {
		return line, err
	}

	if (len(line) >= 4) && (line[3] == '-') {
		code := line[:3] + " "
		for {
			str, err := f.receiveLine()
			line = line + str
			if err != nil {
				return line, err
			}

			if len(str) < 4 {
				if f.debug {
					utils.New().Warnning("Uncorrectly terminated response")
				}
				break
			} else {
				if str[:4] == code {
					break
				}
			}
		}
	}
	f.ReadAndDiscard()
	return line, nil
}

func (f *FTP) ReadAndDiscard() (int, error) {
	var count int

	bufferSize := f.reader.Buffered()
	for count = 0; count < bufferSize; count++ {
		if _, err := f.reader.ReadByte(); err != nil {
			return count, err
		}
	}
	return count, nil
}

func (f *FTP) Syst() (line string, err error) {
	if err = f.send("SYST"); err != nil {
		return
	}

	if line, err = f.receive(); err != nil {
		return
	}

	if !strings.HasPrefix(line, StatusSystemType) {
		err = errors.New(line)
		return
	}

	return strings.SplitN(strings.TrimSpace(line), " ", 2)[1], nil
}

func (f *FTP) Stat(path string) ([]string, error) {
	if err := f.send("STAT %s", path); err != nil {
		return nil, err
	}

	stat, err := f.receive()
	if err != nil {
		return nil, err
	}

	if !strings.HasPrefix(stat, StatusFileStatus) &&
		!strings.HasPrefix(stat, StatusDirectoryStatus) &&
		!strings.HasPrefix(stat, StatusSystemStatus) {
		return nil, errors.New(stat)
	}

	if strings.HasPrefix(stat, StatusSystemStatus) {
		return strings.Split(stat, "\n"), nil
	}

	lines := []string{}

	for _, line := range strings.Split(stat, "\n") {
		if strings.HasPrefix(line, StatusFileStatus) {
			continue
		}
		lines = append(lines, strings.TrimSpace(line))

	}

	return lines, nil
}

func (f *FTP) Retr(path string, retrFn RetrFunc) (s string, err error) {
	if err = f.Type(TypeImage); err != nil {
		return
	}

	var port int
	if port, err = f.Pasv(); err != nil {
		return
	}

	if err = f.send("RETR %s", path); err != nil {
		return
	}

	var pconn net.Conn
	if pconn, err = f.newConnector(port); err != nil {
		return
	}
	defer pconn.Close()

	var line string
	if line, err = f.receiveNoDiscard(); err != nil {
		return
	}

	if !strings.HasPrefix(line, StatusFileOK) {
		err = errors.New(line)
		return
	}

	if err = retrFn(pconn); err != nil {
		return
	}

	pconn.Close()

	if line, err = f.receive(); err != nil {
		return
	}

	if !strings.HasPrefix(line, StatusClosingDataConnection) {
		err = errors.New(line)
		return
	}

	return
}

func (f *FTP) Login(username, password string) (err error) {
	if _, err = f.cmd("331", "USER %s", username); err != nil {
		if strings.HasPrefix(err.Error(), "230") {
			err = nil
		} else {
			return
		}
	}

	if _, err = f.cmd("230", "PASS %s", password); err != nil {
		return
	}

	return
}

func Connect(addr string) (*FTP, error) {
	var err error
	var conn net.Conn

	if conn, err = net.Dial("tcp", addr); err != nil {
		return nil, err
	}

	writer := bufio.NewWriter(conn)
	reader := bufio.NewReader(conn)

	object := &FTP{conn: conn, addr: addr, reader: reader, writer: writer, debug: false}
	object.receive()

	return object, nil
}

func ConnectDbg(addr string) (*FTP, error) {
	var err error
	var conn net.Conn

	if conn, err = net.Dial("tcp", addr); err != nil {
		return nil, err
	}

	writer := bufio.NewWriter(conn)
	reader := bufio.NewReader(conn)

	var line string

	object := &FTP{conn: conn, addr: addr, reader: reader, writer: writer, debug: true}
	line, _ = object.receive()

	log.Print(line)

	return object, nil
}

func (f *FTP) Size(path string) (size int, err error) {
	line, err := f.cmd("213", "SIZE %s", path)

	if err != nil {
		return 0, err
	}

	return strconv.Atoi(line[4 : len(line)-2])
}

func GetFTPStatusText(code string) string {
	return statusText[code]
}
