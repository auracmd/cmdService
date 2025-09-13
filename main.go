package main

import (
	"fmt"

	"github.com/auracmd/cmdService/core"
)

func main() {
	service := core.NewCMD()
	// test linux commnads
	if service.IsLinux() {
		service.Execute("echo hello")
		service.Execute("ls -la")
		service.Execute("a") // command not found.
	} else { // TODO: test windows commnads
		fmt.Println("Not Linux")
	}

	service.PrintCommandHistory()
}
