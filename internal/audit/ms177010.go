package audit

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/seaung/gagako/pkg/utils"
)

var (
	negotiateProtocolRequest, _  = hex.DecodeString("00000085ff534d4272000000001853c00000000000000000000000000000fffe00004000006200025043204e4554574f524b2050524f4752414d20312e3000024c414e4d414e312e30000257696e646f777320666f7220576f726b67726f75707320332e316100024c4d312e325830303200024c414e4d414e322e3100024e54204c4d20302e313200")
	sessionSetupRequest, _       = hex.DecodeString("00000088ff534d4273000000001807c00000000000000000000000000000fffe000040000dff00880004110a000000000000000100000000000000d40000004b000000000000570069006e0064006f007700730020003200300030003000200032003100390035000000570069006e0064006f007700730020003200300030003000200035002e0030000000")
	treeConnectRequest, _        = hex.DecodeString("00000060ff534d4275000000001807c00000000000000000000000000000fffe0008400004ff006000080001003500005c005c003100390032002e003100360038002e003100370035002e003100320038005c00490050004300240000003f3f3f3f3f00")
	transNamedPipeRequest, _     = hex.DecodeString("0000004aff534d42250000000018012800000000000000000000000000088ea3010852981000000000ffffffff0000000000000000000000004a0000004a0002002300000007005c504950455c00")
	trans2SessionSetupRequest, _ = hex.DecodeString("0000004eff534d4232000000001807c00000000000000000000000000008fffe000841000f0c0000000100000000000000a6d9a40000000c00420000004e0001000e000d0000000000000000000000000000")
)

func AuditMS17010(host string, timeout time.Duration) {
	var osInfo string

	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:445", host), time.Second*timeout)
	if err != nil {
		utils.New().Warnning(fmt.Sprintf("connect 445 port timeout :: error :: %v\n", err))
		return
	}

	defer conn.Close()

	conn.SetDeadline(time.Now().Add(time.Second * timeout))
	conn.Write(negotiateProtocolRequest)
	replay := make([]byte, 1024)

	if nbyte, err := conn.Read(replay); err != nil || nbyte < 36 {
		utils.New().Warnning(fmt.Sprintf("Read Replay Data from connect error :: %v\n", err))
		return
	}

	if binary.LittleEndian.Uint32(replay[9:13]) != 0 {
		return
	}

	conn.Write(sessionSetupRequest)

	n, err := conn.Read(replay)
	if err != nil || n < 36 {
		return
	}

	if binary.LittleEndian.Uint32(replay[9:13]) != 0 {
		return
	}

	sessionSetupResponse := replay[36:n]
	if count := sessionSetupResponse[0]; count != 0 {
		byteCount := binary.LittleEndian.Uint16(sessionSetupResponse[7:9])
		if n != int(byteCount)+45 {
			utils.New().Warnning("Invalid session setup Andx response")
		} else {
			for i := 10; i < len(sessionSetupResponse)-1; i++ {
				if sessionSetupResponse[i] == 0 && sessionSetupResponse[i+1] == 0 {
					osInfo = string(sessionSetupResponse[10:i])
					osInfo = strings.Replace(osInfo, string([]byte{0x00}), "", -1)
					break
				}
			}
		}
	}

	userID := replay[32:34]
	treeConnectRequest[32] = userID[0]
	treeConnectRequest[33] = userID[1]

	conn.Write(treeConnectRequest)

	if n, err := conn.Read(replay); err != nil || n < 36 {
		return
	}

	treeID := replay[28:30]
	transNamedPipeRequest[28] = treeID[0]
	transNamedPipeRequest[29] = treeID[1]
	transNamedPipeRequest[32] = userID[0]
	transNamedPipeRequest[33] = userID[1]

	conn.Write(transNamedPipeRequest)
	if n, err := conn.Read(replay); err != nil || n < 36 {
		return
	}

	if replay[9] == 0x05 && replay[10] == 0x02 && replay[11] == 0x00 && replay[12] == 0xc0 {
		trans2SessionSetupRequest[28] = treeID[0]
		trans2SessionSetupRequest[29] = treeID[1]
		trans2SessionSetupRequest[32] = userID[0]
		trans2SessionSetupRequest[33] = userID[1]

		conn.Write(trans2SessionSetupRequest)

		if n, err := conn.Read(replay); err != nil || n < 36 {
			return
		}

		if replay[34] == 0x51 {
			utils.New().Success(fmt.Sprintf("DOUBLEPULSAR SMB IMPLANT in %s\n", host))
		}

	} else {
		utils.New().Info(fmt.Sprintf("%s\t        \t(%s)\n", host, osInfo))
	}

}
