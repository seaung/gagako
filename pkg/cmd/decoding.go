package cmd

import (
	"github.com/abiosoft/ishell"
)

var decode = &ishell.Cmd{
	Name:     "decoding",
	Help:     "解码模块",
	LongHelp: "这个命令提供了一个解码的接口,需要您提供一个被编码过后的字符串",
	Func: func(c *ishell.Context) {
		c.ShowPrompt(false)
		defer c.ShowPrompt(true)

		c.Print("请您提供一个被编码过后的字符串: ")
		dstr := c.ReadLine()

		dtype := c.MultiChoice([]string{
			"URL-decode",
			"De-Base64",
			"De-Hex",
		}, "请您选择一个需要被解码的类型.")

		decoding(dtype, dstr)
	},
}

func decoding(dtype int, dstr string) {}

func init() {}
