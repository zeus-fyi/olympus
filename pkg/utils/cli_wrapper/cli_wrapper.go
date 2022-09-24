package cli_wrapper

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/zeus-fyi/olympus/pkg/utils/printer"
)

type TaskCmd struct {
	Command string
	Args    []string
	Shell   bool
	Env     []string
	Dir     string

	PrintPath     string
	PrintFilename string
	Print         bool
	Environment   string

	StdIO
}

type StdIO struct {
	// Stdin connect a reader to stdin for the command
	// being executed.
	Stdin io.Reader

	// prints stdout and stderr directly to os.Stdout/err as
	// the command runs.
	StreamStdoutOff bool

	Stdout   string
	Stderr   string
	ExitCode int
}

func (t TaskCmd) shellCmd() *exec.Cmd {
	var args []string
	if len(t.Args) == 0 {
		startArgs := strings.Split(t.Command, " ")
		script := strings.Join(startArgs, " ")
		args = append([]string{"-c"}, fmt.Sprintf("%s", script))

	} else {
		script := strings.Join(t.Args, " ")
		args = append([]string{"-c"}, fmt.Sprintf("%s %s", t.Command, script))
	}
	return exec.Command("/bin/bash", args...)
}

func (t TaskCmd) execCmd() *exec.Cmd {
	if strings.Index(t.Command, " ") > 0 {
		parts := strings.Split(t.Command, " ")
		command := parts[0]
		args := parts[1:]
		return exec.Command(command, args...)

	} else {
		return exec.Command(t.Command, t.Args...)
	}
}

func (t *TaskCmd) SetPrintOptions(print bool, printPath, printFilename, env string) {
	t.Print = print
	t.PrintFilename = printFilename
	t.PrintPath = printPath
	t.Environment = env
}

func (t TaskCmd) ExecuteCmd() (string, string, error) {
	var cmd *exec.Cmd
	if t.Shell {
		cmd = t.shellCmd()
	} else {
		cmd = t.execCmd()
	}

	cmd.Dir = t.Dir
	if len(t.Env) > 0 {
		overrides := map[string]bool{}
		for _, env := range t.Env {
			key := strings.Split(env, "=")[0]
			overrides[key] = true
			cmd.Env = append(cmd.Env, env)
		}

		for _, env := range os.Environ() {
			key := strings.Split(env, "=")[0]
			if _, ok := overrides[key]; !ok {
				cmd.Env = append(cmd.Env, env)
			}
		}
	}

	if t.Stdin != nil {
		cmd.Stdin = t.Stdin
	}

	stdoutBuff := bytes.Buffer{}
	stderrBuff := bytes.Buffer{}

	var stdoutWriters io.Writer
	var stderrWriters io.Writer

	if !t.StreamStdoutOff {
		stdoutWriters = io.MultiWriter(os.Stdout, &stdoutBuff)
		stderrWriters = io.MultiWriter(os.Stderr, &stderrBuff)
	} else {
		stdoutWriters = &stdoutBuff
		stderrWriters = &stderrBuff
	}

	cmd.Stdout = stdoutWriters
	cmd.Stderr = stderrWriters

	err := cmd.Start()
	if err != nil {
		return "", "", err
	}

	t.ExitCode = 0
	execErr := cmd.Wait()
	if execErr != nil {
		if exitError, ok := execErr.(*exec.ExitError); ok {
			t.ExitCode = exitError.ExitCode()
		}
	}

	if t.Print && t.PrintFilename != "" && t.PrintPath != "" && t.Environment != "" {
		printer.Printer(t.PrintPath, t.PrintFilename, t.Environment, stdoutBuff.Bytes())
	}
	stdOut, stdErr := string(stdoutBuff.Bytes()), string(stderrBuff.Bytes())
	return stdOut, stdErr, nil
}
