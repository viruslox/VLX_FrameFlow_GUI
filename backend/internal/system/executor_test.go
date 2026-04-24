package system

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func fakeExecCommand(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = append(os.Environ(), "GO_WANT_HELPER_PROCESS=1")
	return cmd
}

func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}

	// os.Args is like: [binary, -test.run=TestHelperProcess, --, command, arg1, arg2, ...]
	args := os.Args
	for len(args) > 0 {
		if args[0] == "--" {
			args = args[1:]
			break
		}
		args = args[1:]
	}

	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "No command provided to helper process\n")
		os.Exit(1)
	}

	// This expects the environment variables to control behavior.
	if os.Getenv("TEST_FAIL") == "1" {
		fmt.Fprintf(os.Stderr, "simulated error")
		os.Exit(1)
	}

	// Output exactly what was passed.
	fmt.Printf("%s\n", strings.Join(args, " "))
	os.Exit(0)
}

func TestExecutor_Run_Success(t *testing.T) {
	// Mock execCommand
	originalExecCommand := execCommand
	execCommand = fakeExecCommand
	defer func() { execCommand = originalExecCommand }()

	executor := NewExecutor()

	scriptName := "my_script"
	args := []string{"arg1", "arg2"}

	output, err := executor.Run(scriptName, args...)

	assert.NoError(t, err)

	// Expected output: "bash -ic my_script "$@" -- arg1 arg2"
	// Let's reconstruct exactly what the helper will print.
	// Command: "bash"
	// Args: "-ic" "my_script \"$@\"" "--" "arg1" "arg2"
	expectedOutput := "bash -ic my_script \"$@\" -- arg1 arg2\n"
	assert.Equal(t, expectedOutput, output)
}

func TestExecutor_Run_Failure(t *testing.T) {
	// Mock execCommand
	originalExecCommand := execCommand
	execCommand = fakeExecCommand
	defer func() { execCommand = originalExecCommand }()

	// Force failure
	os.Setenv("TEST_FAIL", "1")
	defer os.Unsetenv("TEST_FAIL")

	executor := NewExecutor()

	scriptName := "my_script"

	output, err := executor.Run(scriptName)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "simulated error")
	assert.Contains(t, err.Error(), "exit status 1")
	assert.Empty(t, output)
}

func TestExecutor_Run_NoArgs(t *testing.T) {
	// Mock execCommand
	originalExecCommand := execCommand
	execCommand = fakeExecCommand
	defer func() { execCommand = originalExecCommand }()

	executor := NewExecutor()

	scriptName := "my_script"

	output, err := executor.Run(scriptName)

	assert.NoError(t, err)

	// Command: "bash"
	// Args: "-ic" "my_script \"$@\"" "--"
	expectedOutput := "bash -ic my_script \"$@\" --\n"
	assert.Equal(t, expectedOutput, output)
}
