package main

import (
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(args []string, env Environment) (returnCode int) {
	if len(args) == 0 {
		return -1
	}

	for envVarName, envInfo := range env {
		if err := os.Unsetenv(envVarName); err != nil {
			return -1
		}

		if !envInfo.NeedRemove {
			if err := os.Setenv(envVarName, envInfo.Value); err != nil {
				return -1
			}
		}
	}

	name := args[0]
	command := args[1:]

	cmd := exec.Command(name, command...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()

	err := cmd.Run()
	if err != nil {
		return -1
	}
	return 0
}
