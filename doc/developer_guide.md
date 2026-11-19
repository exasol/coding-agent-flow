# Developer Guide

## Requirements Document Detection

CAF recognizes a direct regular file in `doc/requirements/` as a requirements
document when it uses one of these extensions:

- `.md` for Markdown;
- `.rst` for reStructuredText.

The check does not search subdirectories or validate document content.

## Link Checker Setup

The Link Checker runs the lychee binary from `LYCHEE_BINARY`. Without it, CAF downloads the supported lychee version into the user cache directory on first use, for example `~/.cache/coding-agent-flow/lychee/`. CAF downloads lychee for Linux and macOS on amd64 and arm64; on other platforms, set `LYCHEE_BINARY`.

Run the acceptance tests of the Link Checker:

```shell
go test ./test/acceptance/code-validation/...
```

They check URLs against a local web server on `127.0.0.1`.

## Continuous Integration

The GitHub Actions workflow in `.github/workflows/ci-build.yml` runs on every pull request and on pushes to `main`. It verifies:

* all packages build;
* unit and acceptance tests pass; and
* the OFT trace of the repository passes.

## Local Setup

On Linux, let Go switch to the project's Go version:

```shell
go env -w GOTOOLCHAIN=auto
```

### OpenFastTrace JAR

The tests need Java on the `PATH`. OFT validation runs the OpenFastTrace JAR of the supported version, `openFastTraceVersion` in `internal/oftvalidation/adapter/openfasttrace/jar.go` (currently 4.11.0). CAF finds the JAR in this order:

1. `OPENFASTTRACE_JAR`: the path to a JAR. CAF requires `java -jar <JAR> help` to report the supported version. A missing file, another version, or an unrelated JAR fails with an error that names the variable.
2. The local Maven repository: `~/.m2/repository/org/itsallcode/openfasttrace/openfasttrace/<version>/openfasttrace-<version>.jar`.
3. Automatic download: when the JAR is missing in the local Maven repository, CAF runs Maven with `-Dmaven.repo.local=~/.m2/repository` to download it there, also when `~/.m2/settings.xml` sets another local repository. This needs Maven on the `PATH` and network access to Maven Central or to the repository mirror in `~/.m2/settings.xml`.

When the download fails, for example offline or with Maven missing, OFT validation fails with `download OpenFastTrace <version>: …` and the beginning of Maven's output. To work offline, download the JAR once while online, or set `OPENFASTTRACE_JAR`.

The automatic download runs this command, which also downloads the JAR ahead of time:

```shell
mvn --batch-mode org.apache.maven.plugins:maven-dependency-plugin:3.11.0:get -Dmaven.repo.local=$HOME/.m2/repository \
    -Dartifact=org.itsallcode.openfasttrace:openfasttrace:4.11.0 -Dtransitive=false
```

### Running the Tests

Run all tests and the OFT trace:

```shell
go test ./...
java -jar ~/.m2/repository/org/itsallcode/openfasttrace/openfasttrace/4.11.0/openfasttrace-4.11.0.jar trace doc internal test oft_validation.go
```

## Code Structure

Each domain has its own folder under `internal/`, for example `internal/oftvalidation/` for OFT Validation. A domain has four layers:

| Layer | Folder | Content | Example in `internal/oftvalidation/` |
|---|---|---|---|
| Domain | `domain/` | Domain classes and their rules | `OFTValidation`, `OFTSpecificationItemID` |
| Port | `application/port/` | Interfaces that the use cases call | `OFTTracer` |
| Service | `application/service/` | One service per use case | `OFTTracingValidator` |
| Adapter | `adapter/<technology>/` | Implementations of the ports | `adapter/openfasttrace/`: `OFTTracer` with OpenFastTrace |

Domain classes that all domains share go to `internal/global/`.

The root package `codingagentflow` is the public API, with one file per domain, for example `oft_validation.go`. Each constructor passes the adapters to the service of its use case:

```go
func NewOFTTracingValidator() OFTTracingValidator {
	return service.NewOFTTracingValidator(openfasttrace.NewAdapter())
}
```

### Architecture Test

`internal/architecture_test.go` checks every folder under `internal/` that has a `domain/` folder, so a new domain gets the checks automatically:

* archunit checks that each layer has at least one package and that the layers depend inward: the domain refers only to itself, `internal/global/`, and the standard library; the port avoids the service and the adapters; the service avoids the adapters; the adapters avoid the service.
* `go/parser` checks that the domain, the port, and the service avoid `os` and `net`. archunit applies its rules to module packages only.

### Adding a Domain

1. Create `internal/<domain>/` with the four layers.
2. Add the public API of the domain to the root package, in `<domain>.go`.
3. Add the acceptance-test fixture of the domain in `test/acceptance/<domain>/`.

## Acceptance Tests

Acceptance tests are Go integration tests. Each test writes a Target Project to a temporary directory, calls a use case through the public API, and checks the response. The Gherkin scenarios in `doc/requirements/scenarios/` describe the use cases, and OFT traces each scenario `@id` to the `[at -> ...]` marker of its test.

Run the acceptance tests with `go test ./test/acceptance/...`. CI runs every acceptance test and fails for a build constraint in `test/acceptance/`.

### Structure

| File | Content |
|---|---|
| `test/acceptance/fixture.go` | Helpers for all domains |
| `test/acceptance/<domain>/<domain>_fixture.go` | Helpers of one domain, for example `test/acceptance/oft-validation/oft_validation_fixture.go` |
| `test/acceptance/<domain>/<requirement>/<use-case>/*_test.go` | The test files of one use case: file templates, checks, and one test function per scenario |

| Helper | Purpose |
|---|---|
| `acceptance.Run(t, files, given, when, then)` | Builds the input from the files with `given`, calls the use case with `when`, and checks the response with `then` |
| `acceptance.Expect(name, expected, actual)` | Returns an error that names the mismatch, for example `valid: got true, want false` |
| `oftvalidation.File(path, template, parameters...)` | Creates a Target Project file from a template filled with `fmt.Sprintf` |
| `oftvalidation.Validate(t, files, then)` | Runs the OFT Tracing Validator on a Target Project with the files |

### Example

The acceptance test of the use case OFT Requirements Detector is `test/acceptance/oft-validation/oft-requirements-detection/oft-requirements-detector/oft_requirements_detector_test.go`:

```go
package oftrequirementsdetector

const requirementTemplate = `
OFT Requirement document
========================

%s
`

func requirementFiles(path, requirementID string) []acceptance.File {
	return []acceptance.File{oftvalidation.File(path, requirementTemplate, requirementID)}
}

func expectRequirementFound(requirementID string, found bool) func(codingagentflow.OFTValidation) error {
	return func(response codingagentflow.OFTValidation) error {
		id, err := codingagentflow.ParseOFTSpecificationItemID(requirementID)
		if err != nil {
			return err
		}
		return acceptance.Expect(requirementID+" found", found, response.Contains(id))
	}
}

// [at -> scn~oft-requirements-detector~1]
func TestOFTRequirementsDetector(t *testing.T) {
	oftvalidation.Validate(t,
		requirementFiles("doc/requirements/requirements.rst", "req~valid-doc-requirements-rst~1"),
		expectRequirementFound("req~valid-doc-requirements-rst~1", true),
	)

	oftvalidation.Validate(t,
		requirementFiles("requirements.txt", "req~invalid-root-text~1"),
		expectRequirementFound("req~invalid-root-text~1", false),
	)
	// one Validate call per row of the Examples table
}
```

### Adding an Acceptance Test

1. Add the scenario with its `@id` to the `scenarios.feature` of the use case in `doc/requirements/scenarios/`.
2. Add a test function with the marker `// [at -> <scenario id>]` to a test file in the folder of the use case.
3. Build the Target Project files with the helpers of the domain fixture, call the use case, and check the response with an `expect...` function.

### Failing Test

A use-case error stops the test with `When:`. A check mismatch reports `Then:`, for example:

```
Then: req~valid-root-rst~1 found: got false, want true
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
