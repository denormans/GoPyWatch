package gopywatch

import (
	"fmt"
	"gopkg.in/fsnotify.v1"
	"os"
)

func ListenForPythonFileEvents(pythonEnv *PythonEnvironment) {
	fmt.Println("Watching file:", pythonEnv.PythonFilePath)

	fsWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic(fmt.Sprint("Can't setup python file watching via fsnotify:", err))
	}

	err = fsWatcher.Add(pythonEnv.PythonFilePath)
	if err != nil {
		panic(fmt.Sprintf("Error watching python file path %s: %s", pythonEnv.PythonFilePath, err))
	}

	for {
		select {
		case fsEvent := <-fsWatcher.Events:
			if fsEvent.Op != fsnotify.Chmod {
				fmt.Println("Python file changed. Reloading Python environment:", fsEvent)

				if fsEvent.Op == fsnotify.Rename {
					err = fsWatcher.Add(pythonEnv.PythonFilePath)
					if err != nil {
						panic(fmt.Sprintf("Error watching python file path %s: %s", pythonEnv.PythonFilePath, err))
					}
				}

				pythonEnv.SendEvent(Restart)
			}

		case err := <-fsWatcher.Errors:
			fmt.Fprintln(os.Stderr, "Python file watching error:", err)
		}
	}
}

func ListenForExtraDirEvents(pythonEnv *PythonEnvironment, extraWatchDirPath string) {
	fmt.Println("Watching directory:", extraWatchDirPath)

	fsWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic(fmt.Sprint("Can't setup file watching via fsnotify:", err))
	}

	err = fsWatcher.Add(extraWatchDirPath)
	if err != nil {
		panic(fmt.Sprintf("Error watching extra directory path %s: %s", extraWatchDirPath, err))
	}

	for {
		select {
		case fsEvent := <-fsWatcher.Events:
			if fsEvent.Op != fsnotify.Chmod {
				fmt.Println("Something changed in the extra directory. Reloading Python environment:", fsEvent)
				pythonEnv.SendEvent(Restart)
			}

		case err := <-fsWatcher.Errors:
			fmt.Fprintln(os.Stderr, "Extra directory watching error:", err)
		}
	}
}
