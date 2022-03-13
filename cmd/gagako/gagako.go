package main

import (
	"github.com/seaung/gagako/pkg/cmd"
	"github.com/seaung/gagako/pkg/utils"
)

func main() {
	utils.InitConsole()

	cmd.Execute()
}
