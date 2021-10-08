package main

import (
	"errors"
	"log"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	var err error
	for key, value := range env {
		if value.NeedRemove {
			err = os.Unsetenv(key)
			if err != nil {
				log.Println(err)
				return 1
			}
		}
		err = os.Setenv(key, value.Value)
		if err != nil {
			log.Println(err)
			return 1
		}
	}

	currentCommand := cmd[0]
	command := exec.Command(currentCommand, cmd[1:]...)
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	err = command.Run()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return exitErr.ExitCode()
		}
		log.Println(err)
		return 1
	}

	return
}
