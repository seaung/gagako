package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/Ullaakut/nmap/v2"
	"github.com/abiosoft/ishell"
	"github.com/seaung/gagako/pkg/utils"
)

var detect = &ishell.Cmd{
	Name:     "detect",
	Help:     "用于发现一些有价值的信息",
	LongHelp: "这个命令提供了一些发现目标服务和特征的功能",
	Func: func(c *ishell.Context) {
		c.ShowPrompt(false)
		defer c.ShowPrompt(true)

		utils.New().Success("请提供一个目标 e.g 192.168.12/24")
		hosts := c.ReadLine()
		utils.New().Success("请提供一个端口范围 e.g 80-9000")
		ports := c.ReadLine()

		istype := c.MultiChoice([]string{
			"Detect-Services",
			"List-Interface",
		}, "请您选择一个扫描的类型")

		run(istype, hosts, ports)
	},
}

func init() {
	rootCmd.AddCmd(detect)
}

func run(istype int, hosts, ports string) {
	switch istype {
	default:
		utils.New().Warnning("未知的扫描类型!")
	case 0:
		detectServices(hosts, ports)
	case 1:
		listInterfaces()
	}
}

func detectServices(hosts, ports string) {
	scanner, err := nmap.NewScanner(
		nmap.WithTargets(hosts),
		nmap.WithPorts(ports),
		nmap.WithServiceInfo(),
		nmap.WithTimingTemplate(nmap.TimingAggressive),
	)

	if err != nil {
		utils.New().LoggerError(fmt.Sprintf("Unable to create nmap scanner: %v\n", err))
	}

	result, _, err := scanner.Run()
	if err != nil {
		utils.New().LoggerError(fmt.Sprintf("Nmap scan failed: %v\n", err))
	}

	for _, host := range result.Hosts {
		utils.New().Success(fmt.Sprintf("Host %s\n", host.Addresses[0]))
		for _, port := range host.Ports {
			utils.New().Success(fmt.Sprintf("\tPort %d open \n", port.ID))
			utils.New().Success(fmt.Sprintf("\tTarget Serivice info: %s\n", port.Service.ServiceFP))
		}
	}
}

func listInterfaces() {
	scanner, err := nmap.NewScanner()
	if err != nil {
		utils.New().LoggerError(fmt.Sprintf("Unable to create nmap scanner: %v\n", err))
	}

	interfaceLists, err := scanner.GetInterfaceList()
	if err != nil {
		utils.New().LoggerError(fmt.Sprintf("Could not get interface list: %v\n", err))
	}

	content, err := json.MarshalIndent(interfaceLists, "", "\t")
	if err != nil {
		utils.New().LoggerError(fmt.Sprintf("Unable to marshal: %v\n", err))
	}

	utils.New().Success(string(content))
}
