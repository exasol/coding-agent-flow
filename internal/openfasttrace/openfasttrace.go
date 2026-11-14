// Package openfasttrace provides the OpenFastTrace command-line adapter.
package openfasttrace

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/exasol/coding-agent-flow/internal/oftvalidation"
)

// OpenFastTraceVersion is the OpenFastTrace version used for validation.
const OpenFastTraceVersion = "4.11.0"

const (
	openFastTraceArtifact        = "org.itsallcode.openfasttrace:openfasttrace:"
	mavenDependencyPluginVersion = "3.11.0"
	mavenDependencyPluginGet     = "org.apache.maven.plugins:maven-dependency-plugin:" + mavenDependencyPluginVersion + ":get"
	mavenRepositoryDirectory     = ".m2/repository"
)

// Adapter runs OpenFastTrace as a validation engine.
type Adapter struct{}

// NewAdapter creates an OpenFastTrace validation engine.
func NewAdapter() Adapter {
	return Adapter{}
}

// Validate runs OpenFastTrace against a Target Project.
func (Adapter) Validate(ctx context.Context, targetProject oftvalidation.TargetProject) (oftvalidation.OFTValidation, error) {
	jarFile, err := jarFile(ctx)
	if err != nil {
		return oftvalidation.OFTValidation{}, err
	}

	command := exec.CommandContext(ctx, "java", "-jar", jarFile, "trace", "--output-format", "aspec", targetProject.Path())
	output, err := command.Output()
	if err != nil {
		var exitError *exec.ExitError
		if !errors.As(err, &exitError) {
			return oftvalidation.OFTValidation{}, fmt.Errorf("run OpenFastTrace: %w", err)
		}
	}
	var document traceDocument
	if err := xml.Unmarshal(output, &document); err != nil {
		return oftvalidation.OFTValidation{}, fmt.Errorf("read OFT trace: %w", err)
	}
	return oftvalidation.OFTValidation{Valid: err == nil, OFTElementIDs: document.elementIDs()}, nil
}

type traceDocument struct {
	Collections []traceCollection `xml:"specobjects"`
}

type traceCollection struct {
	Type  oftvalidation.OFTElementType `xml:"doctype,attr"`
	Items []traceItem                  `xml:"specobject"`
}

type traceItem struct {
	ID      string `xml:"id"`
	Version string `xml:"version"`
}

func (document traceDocument) elementIDs() map[oftvalidation.OFTElementType][]oftvalidation.OFTElementID {
	ids := make(map[oftvalidation.OFTElementType][]oftvalidation.OFTElementID)
	for _, collection := range document.Collections {
		for _, item := range collection.Items {
			ids[collection.Type] = append(ids[collection.Type], oftElementID(collection.Type, item))
		}
	}
	return ids
}

func oftElementID(elementType oftvalidation.OFTElementType, item traceItem) oftvalidation.OFTElementID {
	return oftvalidation.OFTElementID(fmt.Sprintf("%s~%s~%s", elementType, item.ID, item.Version))
}

func jarFile(ctx context.Context) (string, error) {
	homeDirectory, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("determine home directory for the Maven repository: %w", err)
	}
	jarFile := filepath.Join(
		homeDirectory,
		mavenRepositoryDirectory,
		"org",
		"itsallcode",
		"openfasttrace",
		"openfasttrace",
		OpenFastTraceVersion,
		"openfasttrace-"+OpenFastTraceVersion+".jar",
	)
	if _, err := os.Stat(jarFile); err == nil {
		return jarFile, nil
	} else if !errors.Is(err, os.ErrNotExist) {
		return "", fmt.Errorf("inspect OpenFastTrace JAR: %w", err)
	}

	command := exec.CommandContext(
		ctx,
		"mvn",
		"--batch-mode",
		mavenDependencyPluginGet,
		"-Dartifact="+openFastTraceArtifact+OpenFastTraceVersion,
		"-Dtransitive=false",
	)
	if output, err := command.CombinedOutput(); err != nil {
		return "", fmt.Errorf("download OpenFastTrace %s: %w: %s", OpenFastTraceVersion, err, output)
	}
	if _, err := os.Stat(jarFile); err != nil {
		return "", fmt.Errorf("find downloaded OpenFastTrace JAR: %w", err)
	}
	return jarFile, nil
}
