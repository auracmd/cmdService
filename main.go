package main

import (
	"fmt"

	"github.com/auracmd/cmdService/core"
)

func main() {
	service := core.NewCMD()
	// test linux commnads
	if service.OSName == core.Linux || service.OSName == core.MacOS {
		service.Execute("echo hello")
		service.Execute("ls -la")
		service.Execute("a") // command not found.
	} else { // TODO: test windows commnads
		fmt.Println("Not Linux")
	}

	service.PrintCommandHistory()
}
