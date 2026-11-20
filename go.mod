module github.com/exasol/coding-agent-flow

go 1.27.4

require (
	github.com/exasol/error-reporting-go v0.2.1
	// Test-only dependency: internal/architecture_test.go uses archunit.
	github.com/kcmvp/archunit v0.2.0-alpha1
	github.com/stretchr/testify v1.11.1
	gopkg.in/yaml.v3 v3.0.1
)

// Indirect test dependencies.
require (
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/fatih/color v1.17.0 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/samber/lo v1.51.0 // indirect
	golang.org/x/mod v0.43.0 // indirect
	golang.org/x/sync v0.25.0 // indirect
	golang.org/x/sys v0.50.0 // indirect
	golang.org/x/text v0.22.0 // indirect
	// archunit requires golang.org/x/tools v0.22.0, which fails to compile with the current Go version.
	golang.org/x/tools v0.53.0 // indirect
)
