package cmd

import (
	"github.com/abiosoft/ishell"
	"github.com/seaung/gagako/internal/database"
	"github.com/seaung/gagako/pkg/utils"
)

var databaseCmd = &ishell.Cmd{
	Name:     "dbconnect",
	Help:     "数据爆破模块",
	LongHelp: "这个命令用于爆破常用数据的密码",
	Func: func(c *ishell.Context) {
		c.ShowPrompt(false)
		defer c.ShowPrompt(true)

		c.Print("请您提供一个爆破目标: ")
		target := c.ReadLine()

		isType := c.MultiChoice([]string{
			"mysql",
			"mariadb",
			"postgre",
			"oracle",
			"mssql",
			"redis",
			"mongodb",
			"memcached",
		}, "请您选择一个数据库类型")

		c.Println("请您提供一个数据库用户名字典文件路径  : ")
		user := c.ReadLine()
		c.Println("请您提供一个数据库用户密码字典文件路径: ")
		pass := c.ReadLine()

		c.Println("您提供的目标为       : ", target)
		c.Println("您选择的数据类型为   : ", isType)
		c.Println("您提供的用户名字典为 : ", user)
		c.Println("您提供的密码字典为   : ", pass)

		bureforceDatabase(target, user, pass, isType)
	},
}

func init() {
	rootCmd.AddCmd(databaseCmd)
}

func bureforceDatabase(target, userdict, passdict string, isType int) {
	switch isType {
	default:
		utils.New().Warnning("未知选项！\n")
	case 0:
	case 1:
		database.BureforceMysql(target, userdict, passdict)
	case 2:
		database.ScanPostgresql(target, userdict, passdict)
	case 3:
		database.OracleBureforce(target, userdict, passdict)
	case 4:
		database.MssalUnauthorizedBureforce(target, userdict, passdict)
	case 5:
		database.RedisUnauthorized2(target)
	case 6:
		database.ScanMogondb(target, userdict, passdict)
	case 7:
		database.MssalUnauthorizedBureforce(target, userdict, passdict)
	}
}
