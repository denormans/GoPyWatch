package gopywatch

import (
	"fmt"
	"github.com/sbinet/go-python"
	"os"
	"strings"
	"time"
)

type PythonEnvironment struct {
	PythonFilePath string
	IsInteractive  bool
	Events         chan *Event
}

func NewPythonEnvironment(pythonFilePath string, isInteractive bool) *PythonEnvironment {
	return &PythonEnvironment{
		PythonFilePath: pythonFilePath,
		IsInteractive:  isInteractive,
		Events:         make(chan *Event),
	}
}

func (pythonEnv *PythonEnvironment) SendEvent(eventType EventType) {
	pythonEnv.Events <- NewEvent(eventType)
}

func (pythonEnv *PythonEnvironment) Restart() {
	pythonEnv.Stop()
	pythonEnv.Run()
}

func (pythonEnv *PythonEnvironment) Run() {
	// python interactions
	err := python.Initialize()
	if err != nil {
		panic(fmt.Sprintf("Couldn't initialize Python environment %s: %s", pythonEnv.PythonFilePath, err))
	}

	var args []string
	args = append(args, pythonEnv.PythonFilePath)

	python.PySys_SetArgv(args)
	err = python.PyRun_SimpleFile(pythonEnv.PythonFilePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error executing python file:", err)
		return
	}

	pythonEnv.SendEvent(ProgramDone)
}

func (pythonEnv *PythonEnvironment) Stop() {
	python.PyErr_SetInterrupt()
	time.Sleep(100 * time.Millisecond)

	err := python.Finalize()
	if err != nil {
		//		fmt.Fprint(os.Stderr, "Error exiting Python environment:", err)
		panic(fmt.Sprint("Error exiting Python environment:", err))
	}
}

func (pythonEnv *PythonEnvironment) ProcessInteractive() {
	fmt.Println("Note: Python interactive mode will reset every time a file changes")
	fmt.Println()
	fmt.Println("Type \"help\", \"copyright\", \"python_credits\" or \"python_license\" for more information")

	var err error
	for line, err := GetNextLine(); err == nil; line, err = GetNextLine() {
		trimmedLine := strings.TrimSpace(line)

		if trimmedLine == "exit" || trimmedLine == "exit()" || trimmedLine == "quit" || trimmedLine == "quit()" {
			fmt.Println("Use Ctrl-D (i.e. EOF) to exit")
			continue
		}

		if trimmedLine == "help" {
			fmt.Println("Use help() for interactive help, or help(object) for help about object.")
			continue
		}

		if trimmedLine == "copyright" || trimmedLine == "copyright()" {
			fmt.Println("Copyright (c) 2015 deNormans.com")
			fmt.Println("All Rights Reserved.")
			fmt.Println()

			line = "copyright()"
		}

		if trimmedLine == "python_credits" || trimmedLine == "python_credits()" {
			line = "credits()"
		}

		if trimmedLine == "python_license" || trimmedLine == "python_license()" {
			line = "license()"
		}

		python.PyRun_SimpleString(line)
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
	}

	// we're always done after the interactive mode is done
	pythonEnv.SendEvent(Quit)
}
