package main

import (
	"denormans/gopywatch"
	"flag"
	"fmt"
	"os"
)

type ErrorCode int

const (
	UsageError ErrorCode = 1 << iota
)

func main() {
	var err error

	var isInteractive bool
	var pythonFilePath string
	var extraWatchDirPath string

	flag.BoolVar(&isInteractive, "interactive", false, "Whether or not to use interactive mode")
	flag.StringVar(&pythonFilePath, "file", "", "The python file to watch (required)")
	flag.StringVar(&extraWatchDirPath, "extraWatchDir", "", "Another directory to watch for file changes (optional)")

	flag.Parse()

	if len(pythonFilePath) == 0 {
		fmt.Fprintln(os.Stderr, "The python file to watch is required")
		flag.PrintDefaults()
		os.Exit(int(UsageError))
	}

	pythonFileInfo, err := os.Stat(pythonFilePath)
	if err != nil {
		ExitWithError(err, UsageError, "Couldn't get info on the file to watch:", pythonFilePath)
	}

	if pythonFileInfo.IsDir() {
		ExitWithError(nil, UsageError, "File to watch is a directory:", pythonFilePath)
	}

	pythonEnv := gopywatch.NewPythonEnvironment(pythonFilePath, isInteractive)

	go gopywatch.ListenForPythonFileEvents(pythonEnv)

	if len(extraWatchDirPath) > 0 {
		extraWatchDirInfo, err := os.Stat(extraWatchDirPath)
		if err != nil {
			ExitWithError(err, UsageError, "Couldn't get info on the extra directory to watch:", extraWatchDirPath)
		}

		if !extraWatchDirInfo.IsDir() {
			ExitWithError(nil, UsageError, "Extra directory to watch is not a directory:", pythonFilePath)
		}

		go gopywatch.ListenForExtraDirEvents(pythonEnv, extraWatchDirPath)
	}

	interactiveStarted := false

	// run the program for the first time
	go pythonEnv.Run()

	// process events until we're done
	for {
		event := <-pythonEnv.Events

		switch event.Type {
		case gopywatch.ProgramDone:
			if isInteractive {
				// interactive mode should never start until the first python run is done
				if !interactiveStarted {
					go pythonEnv.ProcessInteractive()
					interactiveStarted = true
				}
			} else {
				// if the program is done, and we're not in interactive mode, quit
				return
			}

		case gopywatch.Restart:
			go pythonEnv.Restart()

		case gopywatch.Quit:
			return
		}
	}
}

func ExitWithError(err error, exitCode ErrorCode, message ...interface{}) {
	fmt.Fprintln(os.Stderr, message...)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	os.Exit(int(exitCode))
}
