package main

import (
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	if len(cmd) == 0 {
		return 1
	}

	bin := cmd[0]
	args := cmd[1:]

	eCmd := exec.Command(bin, args...)

	for k, v := range env {
		if !v.NeedRemove {
			os.Setenv(k, v.Value)
		} else {
			os.Unsetenv(k)
		}
	}

	eCmd.Stdin = os.Stdin
	eCmd.Stdout = os.Stdout
	eCmd.Stderr = os.Stderr

	if err := eCmd.Run(); err != nil {
		if exitErr, ok := interface{}(err).(*exec.ExitError); ok {
			return exitErr.ExitCode()
		}
		return 127
	}

	return 0
}
