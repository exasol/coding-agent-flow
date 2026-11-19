// Package lychee implements BrokenLinkFinder with the lychee link checker command line.
package lychee

import (
	"bytes"
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	exaerror "github.com/exasol/error-reporting-go"

	"github.com/exasol/coding-agent-flow/internal/codevalidation/application/port"
	"github.com/exasol/coding-agent-flow/internal/codevalidation/domain"
	"github.com/exasol/coding-agent-flow/internal/global"
)

// fileExtensions lists the checked document extensions. lychee reads .md as Markdown, .xml as sitemap, and .rst
// and .txt as plain text.
const fileExtensions = "md,rst,txt,xml"

// lychee exit codes, see https://github.com/lycheeverse/lychee#exit-codes
const (
	exitSuccess      = 0
	exitLinkFailures = 2
)

type Adapter struct{}

var _ port.BrokenLinkFinder = Adapter{}

func NewAdapter() Adapter {
	return Adapter{}
}

// FindBrokenLinks runs lychee on the Target Project and reads its JSON report. For LocalLinks, lychee runs offline.
func (Adapter) FindBrokenLinks(ctx context.Context, targetProject global.TargetProject, scope domain.LinkScope) (domain.LinkValidation, error) {
	// [impl -> req~link-validation~1]
	root, err := filepath.Abs(targetProject.Path())
	if err != nil {
		return domain.LinkValidation{}, runError(err, "")
	}
	binary, err := locateBinary(ctx)
	if err != nil {
		return domain.LinkValidation{}, err
	}
	configuration, err := writeConfiguration()
	if err != nil {
		return domain.LinkValidation{}, runError(err, "")
	}
	defer os.Remove(configuration)
	command := exec.CommandContext(ctx, binary, arguments(root, configuration, scope)...)
	command.Dir = root
	command.WaitDelay = waitDelay
	var stderr bytes.Buffer
	command.Stderr = &stderr
	output, err := command.Output()
	if ctxErr := ctx.Err(); ctxErr != nil {
		return domain.LinkValidation{}, ctxErr
	}
	exitCode := exitSuccess
	if err != nil {
		var exitError *exec.ExitError
		if !errors.As(err, &exitError) || exitError.ExitCode() != exitLinkFailures {
			return domain.LinkValidation{}, runError(err, stderr.String())
		}
		exitCode = exitError.ExitCode()
	}
	return readJSONReport(output, root, exitCode == exitLinkFailures)
}

// configurationContent is the CAF lychee configuration. It replaces the lychee.toml that lychee reads from the
// Target Project by default, and keeps the Target Project unchanged.
const configurationContent = "# lychee configuration of CAF: all settings are command line arguments.\ncache = false\n"

// writeConfiguration writes the CAF lychee configuration to a temporary file outside the Target Project.
func writeConfiguration() (string, error) {
	file, err := os.CreateTemp("", "caf-lychee-*.toml")
	if err != nil {
		return "", err
	}
	if _, err := file.WriteString(configurationContent); err != nil {
		file.Close()
		os.Remove(file.Name())
		return "", err
	}
	if err := file.Close(); err != nil {
		os.Remove(file.Name())
		return "", err
	}
	return file.Name(), nil
}

// waitDelay limits how long the adapter waits for the output of a cancelled lychee process.
const waitDelay = 2 * time.Second

// External Link settings: a URL that answers 429 Too Many Requests counts as valid, and an unresponsive host fails
// after about 30 seconds.
const (
	acceptedStatusCodes = "200..=299,429"
	timeoutSeconds      = "10"
	maximumRetries      = "2"
)

func arguments(root, configuration string, scope domain.LinkScope) []string {
	// [impl -> req~external-link-validation~1]
	result := []string{
		"--config", configuration,
		"--format", "json",
		"--no-progress",
		"--include-fragments",
		"--extensions", fileExtensions,
		"--root-dir", root,
	}
	if scope == domain.LocalAndExternalLinks {
		result = append(result, "--accept", acceptedStatusCodes, "--timeout", timeoutSeconds, "--max-retries", maximumRetries)
	} else {
		result = append(result, "--offline")
	}
	return append(result, ".")
}

func runError(err error, stderr string) error {
	return exaerror.New("E-CAF-2").
		Message("Failed to check links with lychee: {{error|uq}}. {{stderr|uq}}").
		Parameter("error", err).
		Parameter("stderr", strings.TrimSpace(stderr)).
		Mitigation("Ensure that the Target Project Path is a readable directory and that lychee runs on this machine.")
}
