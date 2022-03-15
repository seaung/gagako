package cmd

import (
	"github.com/abiosoft/ishell"
)

var encode = &ishell.Cmd{
	Name:     "encode",
	Help:     "编码模块",
	LongHelp: "这个命令提供了一个编码功能的模块,您可以使用这个命令,对输入的字符进行编码",
	Func: func(c *ishell.Context) {
		c.ShowPrompt(false)
		defer c.ShowPrompt(true)

		c.Print("请您提供一个需要被编码的字符串: ")
		isStr := c.ReadLine()

		eType := c.MultiChoice([]string{
			"URL-Encoding",
			"Base64",
			"Hex",
		}, "请您选择一个编码的类型")

		c.Print("您输入的字符串为: ", isStr)
		c.Print("您选择的编码为: ", eType)
		encoding(eType, isStr)
	},
}

func init() {
	rootCmd.AddCmd(encode)
}

func encoding(etype int, estring string) {}
