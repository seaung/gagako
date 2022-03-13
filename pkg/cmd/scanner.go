package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/Ullaakut/nmap"
	"github.com/abiosoft/ishell"
	"github.com/seaung/gagako/pkg/utils"
)

var scannerCmd = &ishell.Cmd{
	Name: "scan",
	Help: "端口扫描",
	Func: func(c *ishell.Context) {
		c.ShowPrompt(false)
		defer c.ShowPrompt(true)

		c.Print("请您提供一个目标IP地址: ")
		ipaddr := c.ReadLine()
		c.Print("请您提供一个端口或端口的范围,以逗号或-号分割: ")
		ports := c.ReadLine()
		choices := c.MultiChoice([]string{
			"TCP-FULL",
			"TCP-SYN",
			"TCP-ACK",
			"UDP",
		}, "请您选择扫描的类型!")

		scanPort(ipaddr, ports, choices)
	},
}

func init() {
	rootCmd.AddCmd(scannerCmd)
}

func isScanType(code int) string {
	switch code {
	default:
		return "-sT"
	case 0:
		return "-sS"
	case 1:
		return "-sA"
	case 2:
		return "-sU"
	}
}

func scanPort(ipaddr, ports string, stype int) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	scanType := isScanType(stype)

	scanner, err := nmap.NewScanner(
		nmap.WithTargets(ipaddr),
		nmap.WithPorts(ports),
		nmap.WithContext(ctx),
		nmap.WithCustomArguments(scanType),
	)

	if err != nil {
		utils.New().LoggerError(fmt.Sprintf("%v\n", err))
	}

	result, warnnings, err := scanner.Run()
	if err != nil {
		utils.New().LoggerError(fmt.Sprintf("%v\n", err))
	}

	if warnnings != nil {
		utils.New().Warnning(fmt.Sprintf("%v\n", warnnings))
	}

	for _, host := range result.Hosts {
		utils.New().Info(fmt.Sprintf("Host %q:\n", host.Addresses[0]))

		for _, port := range host.Ports {
			utils.New().Info(fmt.Sprintf("\tPort %d/%s %s %s\n", port.ID, port.Protocol, port.State, port.Service.Name))
		}
	}

	utils.New().Info(fmt.Sprintf("scan done: %d hosts up scanned in %3f seconds\n", len(result.Hosts), result.Stats.Finished.Elapsed))
}
