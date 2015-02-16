package gopywatch

import (
	"fmt"
	"os"
	"os/exec"
)

type PythonEnvironment struct {
	PythonFilePath string
	IsInteractive  bool
	Events         chan *Event
	Process        *os.Process
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
	var args []string
	if pythonEnv.IsInteractive {
		args = append(args, "-i")
	}
	args = append(args, pythonEnv.PythonFilePath)

	cmd := exec.Command("python", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Start()
	if err != nil {
		panic(fmt.Sprintf("Couldn't start Python process %s: %s", pythonEnv.PythonFilePath, err))
	}

	pythonEnv.SendEvent(ProgramStarted)

	pythonEnv.Process = cmd.Process

	err = cmd.Wait()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Python process ended with error %s: %s", pythonEnv.PythonFilePath, err)
		return
	}

	pythonEnv.SendEvent(ProgramDone)
	pythonEnv.SendEvent(Quit)
}

func (pythonEnv *PythonEnvironment) Stop() {
	err := pythonEnv.Process.Kill()
	if err != nil {
		//		fmt.Fprint(os.Stderr, "Error exiting Python environment:", err)
		panic(fmt.Sprint("Error killing Python process:", err))
	}
}
