// Package core
package core

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"runtime"
)

type OperationSystem int

const (
	Windows = iota
	Linux
	MacOS
)

type Command struct {
	Input  string `json:"input"`
	Output string `json:"output"`
	Err    string `json:"error,omitempty"`
}

type System struct {
	OS         OperationSystem
	ScriptLang string
}

type CMDService struct {
	System
	Commands       []Command
	CurrentCommand *Command
}

func (s *CMDService) AddCommandToList(command *Command) {
	s.Commands = append(s.Commands, *command)
}

func (s *CMDService) ClearCommnads() {
	s.Commands = []Command{}
}

func (s *CMDService) Execute(cmdStr string) (string, error) {
	var cmd *exec.Cmd

	switch s.OS {
	case Windows:
		cmd = exec.Command("powershell", "-Command", cmdStr)
	case Linux, MacOS:
		cmd = exec.Command("bash", "-c", cmdStr)
	default:
		return "", errors.New("unsupported OS")
	}

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()

	s.CurrentCommand = &Command{
		Input:  cmdStr,
		Output: out.String(),
	}

	if err != nil {
		s.CurrentCommand.Err = err.Error()
		s.AddCommandToList(s.CurrentCommand)
		return s.CurrentCommand.Output, err
	}

	s.AddCommandToList(s.CurrentCommand)
	return s.CurrentCommand.Output, nil
}

func DetectHostOS() System {
	switch runtime.GOOS {
	case "windows":
		return System{OS: Windows, ScriptLang: "powershell"}
	case "linux":
		return System{OS: Linux, ScriptLang: "bash"}
	case "darwin":
		return System{OS: MacOS, ScriptLang: "bash"}
	default:
		return System{OS: Linux, ScriptLang: "sh"}
	}
}

func (s *CMDService) PrintCommandHistory() {
	fmt.Println("=== Command History ===")
	for i, cmd := range s.Commands {
		fmt.Printf("Command #%d:\n%s\n", i+1, cmd.String())
	}
}

func (c *Command) String() string {
	return fmt.Sprintf("[\n\tInput: %s\n\tOutput: %s\n\tError: %s\n],", c.Input, c.Output, c.Err)
}
