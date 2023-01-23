package cmd

import (
	"github.com/abiosoft/ishell"
)

type burteforce struct{}

type options func(*burteforce)

var bruteforceCmd = &ishell.Cmd{
	Name:     "bruteforce",
	Help:     "密码破解模块",
	LongHelp: "这个命令提供了一些常见的密码破解模块",
	Func: func(c *ishell.Context) {
		c.ShowPrompt(false)
		defer c.ShowPrompt(true)

		c.Print("请您输入一个需要破解的密钥: ")
		cryptoStr := c.ReadLine()

		c.Print("请您提供一个字典的路径: ")
		filepath := c.ReadLine()

		ctype := c.MultiChoice([]string{
			"MD5",
			"HASH256",
		}, "请您选择需要被破解的类型")

		burte(ctype, cryptoStr, filepath)
	},
}

func burte(ctype int, cryptoStr, filepath string) {}

func burteMD5(password string) {}

func burteHASH256(password string) {}
