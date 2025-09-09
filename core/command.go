// Package core
package core

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"runtime"
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

type cmdService struct {
	systemDetails
	commands       []command
	currentCommand *command
}

func NewCMD() cmdService {
	return cmdService{
		systemDetails: DetectHostOS(),
	}
}

func (s *cmdService) AddCommandToList(command *command) {
	s.commands = append(s.commands, *command)
}

func (s *cmdService) ClearCommnads() {
	s.commands = []command{}
}

func (s *cmdService) Execute(cmdStr string) (string, error) {
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
		return systemDetails{OSName: Linux, ScriptLang: "sh"}
	}
}

func (s *cmdService) PrintCommandHistory() {
	fmt.Println("=== Command History ===")
	for i, cmd := range s.commands {
		fmt.Printf("Command #%d:\n%s\n", i+1, cmd.String())
	}
}

func (s systemDetails) String() string {
	return fmt.Sprintf("[OS: %s, Script: %s]", s.OSName, s.ScriptLang)
}

func (c command) String() string {
	return fmt.Sprintf("[\n\tInput: %s\n\tOutput: %s\n\tError: %s\n],", c.Input, c.Output, c.Err)
}
