package cmd

import "github.com/abiosoft/ishell"

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
			"postgre",
			"mariadb",
			"oracle",
			"mssql",
		}, "请您选择一个数据库类型")

		c.Print("请您提供一个数据库用户名字典文件: ")
		user := c.ReadLine()
		c.Print("请您提供一个数据库用户密码字典文件: ")
		pass := c.ReadLine()

		c.Print("您提供的目标为 : ", target)
		c.Print("您选择的数据类型为 : ", isType)
		c.Print("您提供的用户名字典为 ： ", user)
		c.Print("您提供的密码字典为 : ", pass)
	},
}

func init() {
	rootCmd.AddCmd(databaseCmd)
}
