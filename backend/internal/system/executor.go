package system

import (
	"bytes"
	"fmt"
	"os/exec"
	"os/user"
	"path/filepath"
)

type Executor struct {
	ScriptPath string
}

func NewExecutor(scriptPath string) *Executor {
	return &Executor{
		ScriptPath: scriptPath,
	}
}

// Run executes a bash script by name with the given arguments.
// It will attempt to use sudo -u FrameFlow if the current user is not FrameFlow.
func (e *Executor) Run(scriptName string, args ...string) (string, error) {
	fullPath := filepath.Join(e.ScriptPath, scriptName)

	var cmdArgs []string

	// Check current user
	currentUser, err := user.Current()
	useSudo := false
	if err == nil && currentUser.Username != "FrameFlow" {
		// Check if we are running as root or have sudo privileges.
		// For simplicity, we just prepend sudo -u FrameFlow if we are not FrameFlow.
		useSudo = true
	}

	if useSudo {
		cmdArgs = append([]string{"-u", "FrameFlow", fullPath}, args...)
	} else {
		cmdArgs = append([]string{fullPath}, args...)
	}

	var cmd *exec.Cmd
	if useSudo {
		cmd = exec.Command("sudo", cmdArgs...)
	} else {
		cmd = exec.Command(cmdArgs[0], cmdArgs[1:]...)
	}

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		return "", fmt.Errorf("failed to execute %s: %v, stderr: %s", scriptName, err, stderr.String())
	}

	return out.String(), nil
}
