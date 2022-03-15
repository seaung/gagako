package utils

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
)

var (
	author  string = "seaung"
	version string = "1.0.0"
)

func checkSudo() {
	if os.Geteuid() != 0 {
		New().LoggerError("[!] This program need to have root permission to execute nmap for now.")
		os.Exit(1)
	}
}

func showBanner() {
	name := fmt.Sprintf("gagako (v.%s)", version)
	banner := `
                            __       
   ____ _____ _____ _____ _/ /______ 
  / __ '/ __ '/ __ '/ __ '/ //_/ __ \
 / /_/ / /_/ / /_/ / /_/ / ,< / /_/ /
 \__, /\__,_/\__, /\__,_/_/|_|\____/ 
/____/      /____/                   

	`

	lines := strings.Split(banner, "\n")
	width := len(lines[1])

	fmt.Println(banner)
	color.Green(fmt.Sprintf("%[1]*s", -width, fmt.Sprintf("%[1]*s", (width+len(name))/2, name)))
	color.Blue(fmt.Sprintf("%[1]*s", -width, fmt.Sprintf("%[1]*s", (width+len(author))/2, author)))
	fmt.Println()
}

func InitConsole() {
	checkSudo()
	showBanner()
}
