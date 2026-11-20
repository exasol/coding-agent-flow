module github.com/exasol/coding-agent-flow

go 1.27.4

require (
	github.com/exasol/error-reporting-go v0.2.1
	// Test-only dependency: internal/architecture_test.go uses archunit.
	github.com/kcmvp/archunit v0.2.0-alpha1
	github.com/stretchr/testify v1.11.1
	gopkg.in/yaml.v3 v3.0.1
)

// Test-only dependencies of archunit.
require (
	github.com/fatih/color v1.17.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/samber/lo v1.51.0 // indirect
	golang.org/x/mod v0.43.0 // indirect
	golang.org/x/sync v0.25.0 // indirect
	golang.org/x/sys v0.50.0 // indirect
	golang.org/x/text v0.22.0 // indirect
	// archunit requires golang.org/x/tools v0.22.0, which fails to compile with the current Go version.
	golang.org/x/tools v0.53.0 // indirect
	// Below dependencies are required indirectly by testify/assert:
	github.com/davecgh/go-spew v1.1.1 // indirect; indirect - formats values in assertion failures.
	github.com/pmezard/go-difflib v1.0.0 // indirect; indirect - displays value differences.
)
