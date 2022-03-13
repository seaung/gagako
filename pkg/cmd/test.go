package cmd

import (
	"github.com/abiosoft/ishell"
)

var tCmd = &ishell.Cmd{
	Name: "test",
	Help: "test help",
	Func: func(c *ishell.Context) {
		c.Println("test")
	},
}

func init() {
	rootCmd.AddCmd(tCmd)
}
