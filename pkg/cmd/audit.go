package cmd

import "github.com/abiosoft/ishell"

var auditCmd = &ishell.Cmd{
	Name:     "audit",
	Help:     "审计漏洞模块",
	LongHelp: "这个模块只负责漏洞的审计工作，并不提供漏洞利用的功能",
}

func init() {
	rootCmd.AddCmd(auditCmd)
}
