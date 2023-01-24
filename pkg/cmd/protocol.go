package cmd

import "github.com/abiosoft/ishell"

var protoCmd = &ishell.Cmd{
	Name:     "protocol",
	Help:     "根据协议探测主机是否存活",
	LongHelp: "根据用户提供的协议类型来探测内网主机是否存活",
	Func: func(c *ishell.Context) {
		c.ShowPrompt(false)
		defer c.ShowPrompt(true)

		c.Println("请您提供一个检测目标 : ")
		target := c.ReadLine()

		categories := c.MultiChoice([]string{
			"HTTP",
			"ICMP",
			"SNMP",
			"T3",
			"SMB",
			"TCP",
		}, "请您选择一个网络协议")

		c.Printf("您提供的检测目标为 : %s 选择的网络协议为 : %s", target, categories)
	},
}

func init() {
	rootCmd.AddCmd(protoCmd)
}

func detectHostFromProtocol(target string, code int) {}
