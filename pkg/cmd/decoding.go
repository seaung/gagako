package cmd

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/url"

	"github.com/abiosoft/ishell"
	"github.com/seaung/gagako/pkg/utils"
)

var decode = &ishell.Cmd{
	Name:     "decoding",
	Help:     "解码模块",
	LongHelp: "这个命令提供了一个解码的接口,需要您提供一个被编码过后的字符串",
	Func: func(c *ishell.Context) {
		c.ShowPrompt(false)
		defer c.ShowPrompt(true)

		c.Println("请您提供一个被编码过后的字符串: ")
		dstr := c.ReadLine()

		dtype := c.MultiChoice([]string{
			"URL-decode",
			"De-Base64",
			"De-Hex",
		}, "请您选择一个需要被解码的类型.")

		decoding(dtype, dstr)
	},
}

func init() {
	rootCmd.AddCmd(decode)
}

func decoding(dtype int, dstr string) {
	switch dtype {
	default:
		utils.New().LoggerError("未知编码!")
	case 0:
		utils.New().Success(fmt.Sprintf("URL解码结果: %s\n", urlDecoding(dstr)))
	case 1:
		utils.New().Success(fmt.Sprintf("Base64 解码结果: %s\n", base64decoding(dstr)))
	case 2:
		utils.New().Success(fmt.Sprintf("HEX 解码结果: %s\n", hexDecoding(dstr)))
	}
}

func base64decoding(str string) string {
	d, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		utils.New().LoggerError(fmt.Sprintf("Base64 decoding error: %v\n", err))
		return "is invalid base64 string"
	}
	return string(d)
}

func hexDecoding(str string) string {
	d, err := hex.DecodeString(str)
	if err != nil {
		utils.New().LoggerError(fmt.Sprintf("Hex decoding error: %v\n", err))
		return "is invalid hex string"
	}
	return string(d)
}

func urlDecoding(src string) string {
	e, err := url.QueryUnescape(src)
	if err != nil {
		utils.New().LoggerError(fmt.Sprintf("URL decoding error: %v\n", err))
		return "is invalid url string"
	}
	return e
}
