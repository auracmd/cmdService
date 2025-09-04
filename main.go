package main

import (
	"fmt"

	"github.com/auracmd/cmdService/core"
)

func main() {
	service := core.CMDService{
		System: core.DetectHostOS(),
	}

	// test linux commnads
	if service.OS == core.Linux || service.OS == core.MacOS {
		service.Execute("echo hello")
		service.Execute("ls -la")
		service.Execute("a") // command not found.
	} else { // TODO: test windows commnads
		fmt.Println("Not Linux")
	}

	service.PrintCommandHistory()
}
