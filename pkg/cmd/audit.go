package cmd

import (
	"fmt"

	"github.com/abiosoft/ishell"
	"github.com/seaung/gagako/internal/audit"
	"github.com/seaung/gagako/pkg/utils"
)

var auditCmd = &ishell.Cmd{
	Name:     "audit",
	Help:     "审计漏洞模块",
	LongHelp: "这个模块只负责漏洞的审计工作，并不提供漏洞利用的功能",
	Func: func(c *ishell.Context) {
		c.ShowPrompt(false)
		defer c.ShowPrompt(true)

		c.Println("请您键入一个目标:")
		target := c.ReadLine()
		isType := c.MultiChoice([]string{
			"CVE-2021-21972",
			"CVE-2021-26855",
			"MS17010",
		}, "请您选择一个漏洞类型")

		c.Println()
		utils.New().Info(fmt.Sprintf("您提供的目标为 : %s 漏洞类型为 : %d", target, isType))
		c.Println()

		auditVuln(target, isType)
	},
}

func init() {
	rootCmd.AddCmd(auditCmd)
}

func auditVuln(target string, isType int) {
	switch isType {
	default:
		utils.New().Warnning("未知类型！")
	case 0:
		audit.CVE202121972(target)
	case 1:
		audit.CVE202126855(target)
	case 2:
		audit.AuditMS17010(target, 3)
	}
}
