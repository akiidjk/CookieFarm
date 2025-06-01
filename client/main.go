package main

import (
	_ "embed"
	"fmt"
	"os"
	"strings"

	"github.com/ByteTheCookies/cookieclient/cmd"
)

func isCompletionCommand() bool {
	for _, arg := range os.Args {
		if strings.Contains(arg, "__complete") || strings.Contains(arg, "completion") {
			return true
		}
	}
	return false
}

//go:embed banner.txt
var banner string

func main() {
	if !isCompletionCommand() {
		fmt.Println(banner)
	}
	cmd.Execute()
}
