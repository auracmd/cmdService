// Package core
package core

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"runtime"
	"sync"
)

type OS string

const (
	Windows = "windows"
	Linux   = "linux"
	MacOS   = "macos"
)

type command struct {
	Input  string `json:"input"`
	Output string `json:"output"`
	Err    string `json:"error,omitempty"`
}

type systemDetails struct {
	OSName     OS     `json:"os"`
	ScriptLang string `json:"script-lang"`
}

type CmdService struct {
	systemDetails
	commands       []command
	currentCommand *command
	mu             sync.Mutex // Mutex for inconcurrent access
}

func NewCMD() CmdService {
	return CmdService{
		systemDetails: DetectHostOS(),
	}
}

func (s *CmdService) AddCommandToList(command *command) {
	s.mu.Lock()         // Lock to ensure thread safety
	defer s.mu.Unlock() // Unlock once the operation is complete
	s.commands = append(s.commands, *command)
}

func (s *CmdService) ClearCommands() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.commands = []command{}
}

func (s *CmdService) Execute(cmdStr string) (string, error) {
	var cmd *exec.Cmd

	switch s.OSName {
	case Windows:
		cmd = exec.Command(s.ScriptLang, "-Command", cmdStr)
	case Linux, MacOS:
		cmd = exec.Command(s.ScriptLang, "-c", cmdStr)
	default:
		return "", errors.New("unsupported OS")
	}

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()

	s.currentCommand = &command{
		Input:  cmdStr,
		Output: out.String(),
	}

	if err != nil {
		s.currentCommand.Err = err.Error()
		s.AddCommandToList(s.currentCommand)
		return s.currentCommand.Output, err
	}

	s.AddCommandToList(s.currentCommand)
	return s.currentCommand.Output, nil
}

func DetectHostOS() systemDetails {
	switch runtime.GOOS {
	case "windows":
		return systemDetails{OSName: Windows, ScriptLang: "powershell"}
	case "linux":
		return systemDetails{OSName: Linux, ScriptLang: "bash"}
	case "darwin":
		return systemDetails{OSName: MacOS, ScriptLang: "bash"}
	default:
		fmt.Printf("unsupported OS detected (%s). Defaulting to Linux.\n", runtime.GOOS)
		return systemDetails{OSName: Linux, ScriptLang: "sh"}
	}
}

func (s *CmdService) PrintCommandHistory() {
	fmt.Println("=== Command History ===")
	for i, cmd := range s.commands {
		fmt.Printf("Command #%d:\n%s\n", i+1, cmd.PrettyString())
	}
}

// 	if service.OSName == core.Linux || service.OSName == core.MacOS {

func (s *CmdService) IsLinux() bool {
	switch s.OSName {
	case Linux, MacOS:
		return true
	default:
		return false
	}
}

func (c command) PrettyString() string {
	return fmt.Sprintf("[Input: %s]\n[Output: %s]\n[Error: %s]", c.Input, c.Output, c.Err)
}
