package cmd

import "github.com/abiosoft/ishell"

var protoCmd = &ishell.Cmd{
	Name:     "protocol",
	Help:     "根据协议探测主机是否存活",
	LongHelp: "根据用户提供的协议类型来探测内网主机是否存活",
}

func init() {
	rootCmd.AddCmd(protoCmd)
}
