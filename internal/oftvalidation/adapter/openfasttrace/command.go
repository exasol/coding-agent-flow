package openfasttrace

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"time"
)

// Limits for the output that CAF keeps from OpenFastTrace.
const (
	// maximumReportSize limits the ASpec report on standard output. A larger report stops OpenFastTrace.
	maximumReportSize = 64 << 20
	// maximumDiagnosticsSize limits the tool output that CAF keeps for error messages. The first lines hold the cause.
	maximumDiagnosticsSize = 8 << 10
	// waitDelay limits how long CAF waits for the output of a stopped process.
	waitDelay = 2 * time.Second
)

// limitedBuffer keeps the first limit bytes written to it and drops the rest. It calls onExceed once, when the
// written data first exceeds the limit.
type limitedBuffer struct {
	buffer   bytes.Buffer
	limit    int
	exceeded bool
	onExceed func()
}

func (limited *limitedBuffer) Write(data []byte) (int, error) {
	if limited.exceeded {
		return len(data), nil
	}
	if room := limited.limit - limited.buffer.Len(); len(data) > room {
		limited.buffer.Write(data[:room])
		limited.exceeded = true
		if limited.onExceed != nil {
			limited.onExceed()
		}
		return len(data), nil
	}
	return limited.buffer.Write(data)
}

// runBounded runs the command and returns its standard output, its standard error, and the error of the run. Standard
// output larger than outputLimit stops the command and returns an error; standard error keeps its first
// maximumDiagnosticsSize bytes.
func runBounded(ctx context.Context, outputLimit int, name string, arguments ...string) ([]byte, string, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	stdout := &limitedBuffer{limit: outputLimit, onExceed: cancel}
	stderr := &limitedBuffer{limit: maximumDiagnosticsSize}
	command := exec.CommandContext(ctx, name, arguments...)
	command.Stdout = stdout
	command.Stderr = stderr
	command.WaitDelay = waitDelay
	err := command.Run()
	if stdout.exceeded {
		return nil, stderr.buffer.String(), fmt.Errorf("%s output exceeds %d bytes", name, outputLimit)
	}
	return stdout.buffer.Bytes(), stderr.buffer.String(), err
}

// runWithDiagnostics runs the command in the current working directory and returns the first maximumDiagnosticsSize
// bytes of its combined standard output and standard error, and the error of the run. The command runs to its end.
func runWithDiagnostics(ctx context.Context, name string, arguments ...string) (string, error) {
	return runWithDiagnosticsIn(ctx, "", name, arguments...)
}

// runWithDiagnosticsIn works like runWithDiagnostics, but runs the command in directory.
func runWithDiagnosticsIn(ctx context.Context, directory string, name string, arguments ...string) (string, error) {
	output := &limitedBuffer{limit: maximumDiagnosticsSize}
	command := exec.CommandContext(ctx, name, arguments...)
	command.Dir = directory
	command.Stdout = output
	command.Stderr = output
	command.WaitDelay = waitDelay
	err := command.Run()
	return output.buffer.String(), err
}
