package ofttracingvalidator

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	codingagentflow "github.com/exasol/coding-agent-flow"
	"github.com/exasol/coding-agent-flow/test/acceptance"
	oftvalidation "github.com/exasol/coding-agent-flow/test/acceptance/oft-validation"
)

const featureTemplate = `
# OFT Feature Document

%s

Needs: req
`

const requirementTemplate = `
# OFT Requirement Document

%s

Covers:

- %s
`

const implementingRequirementTemplate = `
# OFT Requirement Document

%s

Covers:

- %s

Needs: impl
`

const implementationTemplate = "// [impl->%s]\n"

func featureFile(featureID string) acceptance.File {
	return oftvalidation.File("features.md", featureTemplate, featureID)
}

func requirementFile(requirementID, featureID string) acceptance.File {
	return oftvalidation.File("requirements.md", requirementTemplate, requirementID, featureID)
}

func expectValid(valid bool) func(codingagentflow.OFTValidation) error {
	return func(response codingagentflow.OFTValidation) error {
		return acceptance.Expect("valid", valid, response.Valid())
	}
}

// [at -> scn~oft-tracing-validator_valid~1]
func TestOFTTracingValidator_valid(t *testing.T) {
	oftvalidation.Validate(t,
		[]acceptance.File{
			featureFile("feat~valid-fixture-feature~1"),
			requirementFile("req~valid-fixture-requirement~1", "feat~valid-fixture-feature~1"),
		},
		expectValid(true),
	)
}

// [at -> scn~oft-tracing-validator_invalid~1]
func TestOFTTracingValidator_invalid(t *testing.T) {
	oftvalidation.Validate(t,
		[]acceptance.File{featureFile("feat~invalid-fixture-feature~1")},
		expectValid(false),
	)
}

// [at -> scn~oft-tracing-validator_requirement-covers-other-feature~1]
func TestOFTTracingValidator_requirementCoversOtherFeature(t *testing.T) {
	oftvalidation.Validate(t,
		[]acceptance.File{
			featureFile("feat~needed-fixture-feature~1"),
			requirementFile("req~other-feature-fixture-requirement~1", "feat~other-fixture-feature~1"),
		},
		expectValid(false),
	)
}

func expectValidWith(itemID string) func(codingagentflow.OFTValidation) error {
	return func(response codingagentflow.OFTValidation) error {
		id, err := codingagentflow.ParseOFTSpecificationItemID(itemID)
		if err != nil {
			return err
		}
		if !response.Valid() {
			return fmt.Errorf("valid: got false, want true")
		}
		return acceptance.Expect(itemID+" found", true, response.Contains(id))
	}
}

func expectValidWithoutIDs(response codingagentflow.OFTValidation) error {
	output, err := json.Marshal(response)
	if err != nil {
		return err
	}
	if want := `{"valid":true,"oftSpecificationItemIDs":{}}`; string(output) != want {
		return fmt.Errorf("response: got %s, want %s", output, want)
	}
	return nil
}

// [at -> scn~oft-tracing-validator_requirement-covers-other-revision~1]
func TestOFTTracingValidator_requirementCoversOtherRevision(t *testing.T) {
	for feature, covered := range map[string]string{
		"feat~revised-fixture-feature~2": "feat~revised-fixture-feature~1",
		"feat~revised-fixture-feature~1": "feat~revised-fixture-feature~2",
	} {
		oftvalidation.Validate(t,
			[]acceptance.File{featureFile(feature), requirementFile("req~revision-fixture-requirement~1", covered)},
			expectValid(false),
		)
	}
}

// [at -> scn~oft-tracing-validator_unwanted-coverage~1]
func TestOFTTracingValidator_unwantedCoverage(t *testing.T) {
	oftvalidation.Validate(t,
		[]acceptance.File{
			featureFile("feat~unwanted-fixture-feature~1"),
			requirementFile("req~unwanted-fixture-requirement~1", "feat~unwanted-fixture-feature~1"),
			oftvalidation.File("impl.go", implementationTemplate, "feat~unwanted-fixture-feature~1"),
		},
		expectValid(false),
	)
}

// [at -> scn~oft-tracing-validator_duplicate-item~1]
func TestOFTTracingValidator_duplicateItem(t *testing.T) {
	oftvalidation.Validate(t,
		[]acceptance.File{
			featureFile("feat~duplicate-fixture-feature~1"),
			oftvalidation.File("more_features.md", featureTemplate, "feat~duplicate-fixture-feature~1"),
			requirementFile("req~duplicate-fixture-requirement~1", "feat~duplicate-fixture-feature~1"),
		},
		expectValid(false),
	)
}

// [at -> scn~oft-tracing-validator_implementation-coverage~1]
func TestOFTTracingValidator_implementationCoverage(t *testing.T) {
	chain := []acceptance.File{
		featureFile("feat~chain-fixture-feature~1"),
		oftvalidation.File("requirements.md", implementingRequirementTemplate, "req~chain-fixture-requirement~1", "feat~chain-fixture-feature~1"),
	}
	for _, row := range []struct {
		source []acceptance.File
		valid  bool
	}{
		{[]acceptance.File{oftvalidation.File("src/main.go", implementationTemplate, "req~chain-fixture-requirement~1")}, true},
		{nil, false},
		{[]acceptance.File{oftvalidation.File("src/main.go", implementationTemplate, "req~missing-fixture-requirement~1")}, false},
	} {
		oftvalidation.Validate(t, append(chain[:len(chain):len(chain)], row.source...), expectValid(row.valid))
	}
}

// [at -> scn~oft-tracing-validator_no-specification-items~1]
func TestOFTTracingValidator_noSpecificationItems(t *testing.T) {
	oftvalidation.Validate(t, nil, expectValidWithoutIDs)
	oftvalidation.Validate(t, []acceptance.File{oftvalidation.File("README.md", "# Project\n\nNo OFT Specification Items.\n")}, expectValidWithoutIDs)
}

// [at -> scn~oft-tracing-validator_nested-files~1]
func TestOFTTracingValidator_nestedFiles(t *testing.T) {
	oftvalidation.Validate(t,
		[]acceptance.File{
			oftvalidation.File("doc/deep/er/features.md", featureTemplate, "feat~nested-fixture-feature~1"),
			requirementFile("req~nested-fixture-requirement~1", "feat~nested-fixture-feature~1"),
		},
		expectValid(true),
	)
}

// [at -> scn~oft-tracing-validator_target-project-path-no-directory~1]
func TestOFTTracingValidator_targetProjectPathNoDirectory(t *testing.T) {
	oftvalidation.ValidateFails(t, oftvalidation.MissingPath, nil, acceptance.ExpectError("Target Project"))
	oftvalidation.ValidateFails(t, oftvalidation.RegularFile, nil, acceptance.ExpectError("Target Project"))
}

// [at -> scn~oft-tracing-validator_symlinked-target-project~1]
func TestOFTTracingValidator_symlinkedTargetProject(t *testing.T) {
	oftvalidation.ValidateIn(t, oftvalidation.SymbolicLink,
		[]acceptance.File{
			featureFile("feat~linked-fixture-feature~1"),
			requirementFile("req~linked-fixture-requirement~1", "feat~linked-fixture-feature~1"),
		},
		expectValidWith("feat~linked-fixture-feature~1"),
	)
}

// [at -> scn~oft-tracing-validator_unreadable-file~1]
func TestOFTTracingValidator_unreadableFile(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("root reads every file")
	}
	unreadable := func(t *testing.T, files []acceptance.File) codingagentflow.TargetProject {
		targetProject := acceptance.TargetProject(t, files)
		requirements := filepath.Join(targetProject.Path(), "requirements.md")
		if err := os.Chmod(requirements, 0); err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() { _ = os.Chmod(requirements, 0o600) })
		return targetProject
	}
	oftvalidation.ValidateFails(t, unreadable,
		[]acceptance.File{
			featureFile("feat~unreadable-fixture-feature~1"),
			requirementFile("req~unreadable-fixture-requirement~1", "feat~unreadable-fixture-feature~1"),
		},
		acceptance.ExpectError("Error reading"),
	)
}

// [at -> scn~oft-tracing-validator_unusable-jar~1]
func TestOFTTracingValidator_unusableJar(t *testing.T) {
	noJar := filepath.Join(t.TempDir(), "openfasttrace.jar")
	if err := os.WriteFile(noJar, []byte("no JAR content"), 0o600); err != nil {
		t.Fatal(err)
	}
	for _, jar := range []string{filepath.Join(t.TempDir(), "missing.jar"), noJar} {
		t.Setenv("OPENFASTTRACE_JAR", jar)
		oftvalidation.ValidateFails(t, acceptance.TargetProject,
			[]acceptance.File{featureFile("feat~jar-fixture-feature~1")},
			acceptance.ExpectError("OPENFASTTRACE_JAR"),
		)
	}
}

// [at -> scn~oft-tracing-validator_read-only~1]
func TestOFTTracingValidator_readOnly(t *testing.T) {
	files := []acceptance.File{
		featureFile("feat~read-only-fixture-feature~1"),
		requirementFile("req~read-only-fixture-requirement~1", "feat~read-only-fixture-feature~1"),
	}
	oftvalidation.ValidateTargetProject(t, files, func(_ codingagentflow.OFTValidation, targetProject codingagentflow.TargetProject) error {
		return acceptance.ExpectUnchanged(targetProject, files)
	})
}
