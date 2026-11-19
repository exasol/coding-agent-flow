# Handoff: Integration after #43

> [!IMPORTANT]
> This is a temporary document. Remove it after #43 (PR #44, OFT Tracing Validator) is merged to `main` and the steps below are done.

This branch (`skyglass/issue49`, Validation Runner) uses the OFT Tracing Validator of #43. It compiles after #43 is on `main` and the steps below are done. The steps also cover the other open branches that overlap with #43.

## Branches

| Branch | Content | State |
| --- | --- | --- |
| `main` | OFT Validation docs and acceptance tests (#38, #42) | base |
| `skyglass/issue43` (PR #44) | OFT Tracing Validator, archunit architecture test for all domains, CI with Java and the OFT JAR | merge first |
| `skyglass/issue49` (this branch) | Validation Runner with the `oft-tracing-validator` validation | compiles after #43 |
| `skyglass/issue45` (PR #46) | Link Checker and External Link Checker | optional for the MVP |
| `feat/machine-validation-cli-12` (PR #35) | `caf validate` CLI | rework through #47 |

## Steps on this branch after #43 is on `main`

### 1. Merge `main`

Merge `origin/main` into `skyglass/issue49` and resolve the conflicts:

* `go.mod`: take the `go` directive and the archunit requirements of `main`, keep `gopkg.in/yaml.v3` and `github.com/exasol/error-reporting-go`, then run `go mod tidy`.
* `doc/requirements/features/features.md`, `doc/requirements/features/requirements.md`, `doc/requirements/domain/index.md`, `doc/changes/changes_0.1.0.md`: keep the entries of both sides.
* `error_code_config.yml`: see step 5.

### 2. Unify the Target Project

`TargetProject` lives in `internal/global/target_project.go`, with the alias in the root `target_project.go`. #43 has its own `TargetProject` in `internal/oftvalidation/domain/target_project.go`, with an alias in `oft_validation.go`.

* Delete `internal/oftvalidation/domain/target_project.go`.
* Remove `TargetProject` and `NewTargetProject` from `oft_validation.go`.
* Use `global.TargetProject` and `global.NewTargetProject` in all OFT Validation packages, including `internal/oftvalidation/adapter/openfasttrace/adapter_test.go`.

The `oft-tracing-validator` adapter passes `global.TargetProject` to the OFT Tracing Validator, so it compiles after this step.

### 3. Translate the architecture tests to archunit

#43 adds `internal/architecture_test.go`, which checks every domain under `internal/` with archunit. This branch has a `go/parser` test in `internal/validationentrypoint/architecture_test.go`.

* Delete `internal/validationentrypoint/architecture_test.go`. The shared test covers `validationentrypoint` with its layers `domain`, `application/port`, `application/service`, and `adapter`.
* Move the rule of `TestOtherDomainsThroughServiceAndDomain` into the shared test: a domain uses another domain only through its `application/service` and `domain`, and the root package wires the adapters.

```go
// A domain uses another domain through its application service and domain classes only.
for _, other := range domains(t) {
	if other == domain {
		continue
	}
	rules["refers to "+other+" through its service and domain only"] = archunit.Layers(
		layer(domain, "Domain"), layer(domain, "Port"), layer(domain, "Service"), layer(domain, "Adapter"),
	).ShouldNotRefer(archunit.Layers(layer(other, "Port"), layer(other, "Adapter")))
}
```

* Check the new rule with an injected violation, for example an import of `internal/oftvalidation/adapter/openfasttrace` in `internal/validationentrypoint/adapter/ofttracingvalidator`. The test fails with this import and passes after you remove it.

### 4. Use one set of acceptance-test helpers

* `test/acceptance/fixture.go` of `main` provides `File`, `Run`, `WriteFiles`, and `Expect`.
* `test/acceptance/validation-entry-point/validation_entry_point_fixture.go` has its own `TargetProject` helper. Keep one shared `TargetProject` helper in `test/acceptance/` when #45 brings `acceptance.TargetProject` and `acceptance.NewFile`, and use it in all domain fixtures.
* The acceptance tests of this branch run OpenFastTrace, so they need Java and the OFT JAR. The CI workflow of #43 provides both.
* #43 removes the temporary `oftvalidation` build tag and adds a CI check against build constraints in `test/acceptance/`. Keep the acceptance tests of this branch free of build constraints.

### 5. Renumber the error codes

Each branch starts its error codes after the `highest-index` of its base, so the numbers collide:

| Branch | Codes |
| --- | --- |
| `main` | `E-CAF-1` (requirements document check) |
| this branch | `E-CAF-2`, `E-CAF-3` (`.caf.yaml` reader in `adapter/yamlconfig`) |
| `skyglass/issue45` | `E-CAF-2`, `E-CAF-3`, `E-CAF-4` (lychee adapter) |
| PR #35 | `E-CAF-2`, `E-CAF-3` (configuration discovery) |

Give the codes of each branch the next free numbers when it merges, and update `highest-index` in `error_code_config.yml`.

### 6. Remove obsolete code

* #43 replaces the requirements document check (`internal/requirements`, `E-CAF-1`, `req~requirements-document-detection~1`). Check that `main` is free of references to it after the merge, and remove leftovers.
* The unit test of the Validation Runner uses the text `E-CAF-9: broken` as fake error message. Replace it with a plain message, such as `broken`, so `E-CAF` numbers appear only as real error codes.

### 7. Verify

Run the full checks and fix each finding:

```shell
go build ./...
go vet ./...
go test ./...
java -jar ~/.m2/repository/org/itsallcode/openfasttrace/openfasttrace/4.11.0/openfasttrace-4.11.0.jar trace doc internal test *.go
```

* The OFT trace of the whole repository is clean.
* The CI workflow passes, including the architecture test and the check against build constraints.

## Other branches after #43 is on `main`

### `skyglass/issue45` (PR #46, Link Checker)

* `internal/global/target_project.go` and the root `target_project.go` are identical to this branch, so they merge cleanly.
* Delete `internal/codevalidation/architecture_test.go`; the shared archunit test covers `codevalidation`.
* Keep `test/acceptance/global_fixture.go` and `test/acceptance/target_project_fixture.go`, and delete `test/acceptance/fixture.go` of `main`; they declare the same helpers.
* Renumber the lychee error codes (step 5).
* Follow-up #54 moves the link validity decision into the domain, like #52 for OFT Validation.

### PR #35 (`caf validate` CLI)

* Rework it through #47: the migration table in #47 maps each PR #35 package to its target in `internal/validationentrypoint/`.
* Move its requirements from `doc/requirements/system_requirements.md` into the `features`, `domain`, and `scenarios` structure.
* Renumber its error codes (step 5).

## Next tickets

After this branch is merged: #48 (Agent Turn Validator), #50 (Manual Validator), then #47 (CLI Integration). The follow-ups #52 and #53 extend the OFT Tracing Validator, and #55 adds checks on changed files.

## Open points to watch

* `adapter/yamlconfig` reads `.caf.yaml` in the Target Project root. #47 adds the upward search from the working directory.
* The Go version differs between branches (`go` directive and CI `go-version`). Use the version of `main` after the merge.
* The OFT acceptance tests download the OFT JAR through Maven on first use, so the first local run needs network access.
* Check the numbering in `error_code_config.yml` by hand; Exasol's error code crawlers support Java and C#, and Go is outside their scope.

## Remove this document

Delete `doc/handoff.md` after the steps above are done and the branch is merged.
