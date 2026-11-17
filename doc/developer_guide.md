# Developer Guide

## Continuous Integration

The GitHub Actions workflow in `.github/workflows/ci-build.yml` runs on every pull request and on pushes to `main`. It uses Go 1.27.4 and verifies:

* all packages build;
* unit and acceptance tests pass; and
* the OFT trace of the repository passes.

## Local Setup

On Linux, let Go switch to the project's Go version:

```shell
go env -w GOTOOLCHAIN=auto
```

OFT validation runs the OpenFastTrace JAR from `OPENFASTTRACE_JAR` if set, else from `~/.m2/repository`. Download it once with:

```shell
mvn --batch-mode org.apache.maven.plugins:maven-dependency-plugin:3.11.0:get \
    -Dartifact=org.itsallcode.openfasttrace:openfasttrace:4.11.0 -Dtransitive=false
```

Run the tests and the OFT trace:

```shell
go test ./...
java -jar ~/.m2/repository/org/itsallcode/openfasttrace/openfasttrace/4.11.0/openfasttrace-4.11.0.jar trace doc internal test oft_validation.go
```

## OFT Validation Code Structure

The OFT validation code lives in `internal/oftvalidation/`:

| Package | Content | Example |
|---|---|---|
| `domain` | Domain classes and their rules | `OFTValidation`, `OFTElementID` |
| `application/port` | Interfaces the use cases call | `OFTTracer` |
| `application/service` | One service per use case | `OFTTracingValidator` |
| `adapter/openfasttrace` | Implementation of `OFTTracer` with OpenFastTrace | `Adapter` |

The root package `codingagentflow` is the public API. Its constructor passes the adapter to the use case:

```go
func NewOFTTracingValidator() OFTTracingValidator {
	return service.NewOFTTracingValidator(openfasttrace.NewAdapter())
}
```

`internal/oftvalidation/architecture_test.go` checks the imports of each layer.

## Acceptance Tests

An acceptance test runs a use case against the public API with the data of its `at_data.toml`.

The Gherkin scenarios in `doc/requirements/scenarios/` document the use cases. The skill `reconcile-acceptance-tests` writes the acceptance tests from them, and OFT traces each scenario `@id` to its `[at -> ...]` marker.

### Running Acceptance Tests

Run all acceptance tests:

```shell
go test ./test/acceptance/...
```

### Example

This section's example is the acceptance test of the use case OFT Requirements Detector in `test/acceptance/oft-validation/oft-requirements-detection/oft-requirements-detector/`:

| File | Content |
|---|---|
| `at_data.toml` | Input files and expected values per run |
| `at_test.go` | The public call, checks, and test function with one `[at -> ...]` OFT trace marker per scenario |

#### at_data.toml

Each `[[scenario]]` block is one run. Its `id` and `description` name the run in the test output, for example `TestOFTRequirementsDetector/scn~oft-requirements-detection~1:_Detect_OFT_Requirements_by_file_location_and_extension:_requirements.md`. Each block has two parts:

| Part | Content | Used by |
|---|---|---|
| `[scenario.given]` | Input values, and `[[scenario.given.file]]` entries with the files of the target project | The runner |
| `[scenario.then]` | Expected values | The checks |

```toml
[variables]
requirement_document = '''
...
{{req_id}}
'''

[[scenario]]
id = "scn~oft-requirements-detection~1"
description = "Detect OFT Requirements by file location and extension: doc/requirements/requirements.rst"

[scenario.given]
location = "doc/requirements/requirements.rst"
req_id = "req~valid-doc-requirements-rst~1"

[[scenario.given.file]]
path = "{{location}}"
content = "${requirement_document}"

[scenario.then]
found = true

[[scenario]]
...
```

For this block, the runner:

1. replaces `${requirement_document}` with the text from `[variables]`;
2. replaces `{{location}}` and `{{req_id}}` in the file with the values from `[scenario.given]`;
3. writes the file `doc/requirements/requirements.rst` containing `req~valid-doc-requirements-rst~1`.

#### at_test.go

```go
var scenarios = acceptance.NewSuite(
	oftvalidation.When(
		"the caller validates OFT tracing for the target project path",          // the When step
		"codingagentflow.NewOFTTracingValidator().Validate(ctx, targetProject)", // the public method it calls
	),
	checkRequirementFound, // compares the response with the scenario values in `at_data.toml`
)
```

`targetProject` is the target project with the block's `[[scenario.given.file]]` entries.

A check is a Go function that receives the response and the current block's Given and Then values. It returns an error for a mismatch. For the block above, the check receives:

```go
given = acceptance.Values{"location": "doc/requirements/requirements.rst", "req_id": "req~valid-doc-requirements-rst~1"}
then = acceptance.Values{"found": "true"}
```

The check reads the requirement ID from Given and the expected result from Then:

```go
func checkRequirementFound(response codingagentflow.OFTValidation, given acceptance.Values, then acceptance.Values) error {
	reqID := given.String("req_id") // "req~valid-doc-requirements-rst~1"
	expectedID, err := codingagentflow.ParseOFTElementID(reqID)
	if err != nil {
		return err
	}
	// then.Bool("found") is true: the response must contain expectedID
	return acceptance.Expect(reqID+" found", then.Bool("found"), response.Contains(expectedID))
}
```

`acceptance.Expect` names the mismatch, for example `req~valid-root-rst~1 found: expected false, got true`.

The OFT Tracing Validator uses a check with one value, from `valid = true` or `valid = false`:

```go
func checkValid(response codingagentflow.OFTValidation, _ acceptance.Values, then acceptance.Values) error {
	return acceptance.Expect("valid", then.Bool("valid"), response.Valid())
}
```

The shared `When` helper lives in `test/acceptance/oft-validation/oft_validation_fixture.go`.

The same `at_test.go` file ends with the test function:

```go
// [at -> scn~oft-requirements-detection~1]
func TestOFTRequirementsDetector(t *testing.T) {
	scenarios.Run(t)
}
```

#### Using Given values in the public call

Use `given.String("name")` for a value from the current block's `[scenario.given]` table. With the shared helper:

```go
var scenarios = acceptance.NewSuite(
	oftvalidation.When(
		"the caller validates the requested requirement",
		`api.Validate(ctx, targetProject, given.String("requirement_id"))`,
	),
	checkResponse,
)
```

### What scenarios.Run Does

1. Builds one small program that calls the public API. `at_test.go` stores the call as text, so the acceptance-test package compiles without compiling that call. The runner builds the program during the test. A fixture compile error fails the test.
2. Runs that program once per block, with the blocks in parallel: writes the Given files, makes the call, and checks the result against the Then values.

### Red Result

A fixture compile error is reported as a test error:

```
public call codingagentflow.DoesNotExist(ctx, targetProject) does not compile:
undefined: codingagentflow.DoesNotExist
```

A behavior difference produces a check error:

```
Then: item found: expected true, got false
```

## Managing Release Notes

Keep one Markdown file per released version under `doc/changes/`, named `changes_<version>.md`, for example `changes_0.2.0.md`. Add the newest version as the first entry in `doc/changes/changelog.md`.

Prepare the next release-note file in the release preparation pull request. Its heading contains the project name, semantic version, and release date. Add a short summary followed by relevant sections such as `Features`, `Bugfixes`, `Security`, `Documentation`, and `Dependency Updates`. Reference related GitHub issues or pull requests in the entries. Keep notes for released versions unchanged. Publish corrections as a new release when they affect users.

Use semantic versioning for the public Go module API:

* patch: backward-compatible bug fixes, dependency or security updates;
* minor: backward-compatible features; and
* major: incompatible exported-API changes. Go modules require a `/vN` module path suffix for major versions 2 and later.

## Releasing

Create GitHub releases with [Release Droid](https://github.com/exasol/release-droid) and the generic language configuration. Go users consume the module by its Git tag, so there is no binary artifact to upload unless the project later introduces a CLI distribution.

Before each release:

1. Choose the semantic version and update the release-note file and changelog.
2. Update every project-owned version source to that version. When a runtime version is introduced, keep it in `internal/version/version.go`; configure the same version in `.project-keeper.yml`.
3. Ensure the CI workflow is green and run the release validation locally:

   ```shell
   java -jar .../release-droid-*.jar -n coding-agent-flow --local . --goal validate
   ```

4. Merge the prepared release into `main`, then create the GitHub release:

   ```shell
   java -jar .../release-droid-*.jar -n coding-agent-flow --goal release --skipvalidation
   ```

5. Verify that the created tag is the intended Go module version (for example, `v0.1.0`) and that the GitHub release text links to the matching change file.

`release_config.yml` limits Release Droid to GitHub and the generic release mode. It intentionally does not publish assets.
