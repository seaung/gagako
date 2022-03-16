package cmd

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/url"

	"github.com/abiosoft/ishell"
	"github.com/seaung/gagako/pkg/utils"
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

func encoding(etype int, estring string) {
	switch etype {
	default:
		utils.New().LoggerError("未知编码!")
	case 1:
		utils.New().Success(fmt.Sprintf("Base64 编码结果: %s\n", base64decoding(estring)))
	case 0:
		utils.New().Success(fmt.Sprintf("URL 编码结果: %s\n", urlEncoding(estring)))
	case 2:
		utils.New().Success(fmt.Sprintf("HEX 编码结果: %s\n", hexEncoding(estring)))
	}
}

func base64encoding(str string) string {
	src := []byte(str)
	e := base64.StdEncoding.EncodeToString(src)
	return e
}

func hexEncoding(str string) string {
	src := []byte(str)
	hexstr := hex.EncodeToString(src)
	return hexstr
}

func urlEncoding(src string) string {
	return url.QueryEscape(src)
}
