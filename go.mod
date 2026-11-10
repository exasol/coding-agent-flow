module github.com/exasol/coding-agent-flow

go 1.26.8

require (
	github.com/exasol/error-reporting-go v0.2.1
	github.com/stretchr/testify v1.11.1
)

require (
	//below 3 are required indirectly by testify/assert:
	github.com/davecgh/go-spew v1.1.1 // indirect; indirect - formats values in assertion failures.
	github.com/pmezard/go-difflib v1.0.0 // indirect; indirect - displays value differences.
	gopkg.in/yaml.v3 v3.0.1 // indirect; indirect -supports Testify’s YAML assertions.
)
