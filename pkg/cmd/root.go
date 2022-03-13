package cmd

import (
	"github.com/abiosoft/ishell"
)

var rootCmd = ishell.New()

func Execute() {
	rootCmd.Run()
}
