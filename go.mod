module github.com/exasol/coding-agent-flow

go 1.27.4

// Test-only dependencies: internal/architecture_test.go uses archunit. The CAF code itself has no third-party dependency.
require github.com/kcmvp/archunit v0.2.0-alpha1

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
)
