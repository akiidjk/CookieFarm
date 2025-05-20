package main

import (
	_ "embed"
	"fmt"

	"github.com/ByteTheCookies/cookieclient/cmd"
)

//go:embed banner.txt
var banner string

func main() {
	fmt.Println(banner)
	cmd.Execute()
}
