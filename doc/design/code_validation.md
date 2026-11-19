# Code Validation Design

## Link Checker

`dsn~link-checker~1`

The Link Checker delegates link checking to [lychee](https://github.com/lycheeverse/lychee), like OFT Validation delegates tracing to OpenFastTrace. The service `LinkChecker` calls the port `BrokenLinkFinder`, and the adapter `internal/codevalidation/adapter/lychee` implements it.

Covers:

- `req~link-validation~1`

### lychee Command

The adapter runs lychee in the Target Project:

```shell
lychee --config <CAF configuration> --format json --no-progress --include-fragments --extensions md,rst,txt,xml --root-dir <target project> <scope arguments> .
```

The scope arguments come from `dsn~external-link-checker~1`.

The CAF configuration is a temporary file outside the Target Project. It replaces the `lychee.toml` that lychee reads from the Target Project by default, so the Target Project stays unchanged.

### lychee Report

lychee exits with `0` when all Links resolve and with `2` for Broken Links. For these exit codes, the adapter reads the Broken Links from the JSON report and converts `file://` URLs in the Target Project to project-relative paths. It reports other exit codes, and a report that differs from the exit code, as errors.

### lychee Binary

The adapter runs the binary from `LYCHEE_BINARY`, else the supported lychee version from the user cache directory, which it downloads on first use with a SHA-256 checksum check. Before each check, it requires `lychee --version` to report the supported version, because the adapter reads the JSON report of this version. The supported version and the checksums are in `internal/codevalidation/adapter/lychee/binary.go`.

## External Link Checker

`dsn~external-link-checker~1`

The Link Scope selects the lychee arguments. `NewLinkChecker()` checks Local Links, and `NewLinkChecker(WithExternalLinks())` adds External Links.

Covers:

- `req~external-link-validation~1`

| Link Scope | Public API | Scope arguments |
|---|---|---|
| Local Links | `NewLinkChecker()` | `--offline` |
| Local and External Links | `NewLinkChecker(WithExternalLinks())` | `--accept 200..=299,429 --timeout 10 --max-retries 2` |
