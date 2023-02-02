package cmd

import (
	"crypto/md5"
	"crypto/sha256"
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
			"HTML-Encoding",
			"MD5-32",
			"MD5-16",
			"HASH256",
		}, "请您选择一个编码的类型")

		c.Println("您输入的字符串为: ", isStr)
		c.Println("您选择的编码为: ", eType)
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
		utils.New().Success(fmt.Sprintf("Base64 编码结果: %s\n", base64encoding(estring)))
	case 0:
		utils.New().Success(fmt.Sprintf("URL 编码结果: %s\n", urlEncoding(estring)))
	case 2:
		utils.New().Success(fmt.Sprintf("HEX 编码结果: %s\n", hexEncoding(estring)))
	case 3:
		utils.New().Success(fmt.Sprintf("HTML 编码结果: %s\n", htmlEncoding(estring)))
	case 4:
		utils.New().Success(fmt.Sprintf("MD5 32编码结果: %s\n", md5Encoding32(estring)))
	case 5:
		utils.New().Success(fmt.Sprintf("MD5 16编码结果: %s\n", md5Encoding16(estring)))
	case 6:
		utils.New().Success(fmt.Sprintf("HASH256编码结果: %s\n", sha256Encoding(estring)))
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

func htmlEncoding(str string) string {
	var r string

	arr := []rune(str)
	num := len(arr)
	for i := 0; i < num; i++ {
		r += fmt.Sprint("&#", arr[i]) + ";"
	}
	return r
}

func md5Encoding32(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func md5Encoding16(str string) string {
	return md5Encoding32(str)[8:24]
}

func sha256Encoding(str string) string {
	h := sha256.New()
	return hex.EncodeToString(h.Sum(nil))
}
