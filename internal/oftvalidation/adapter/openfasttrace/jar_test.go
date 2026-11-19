package openfasttrace

import (
	"context"
	"path/filepath"
	"strings"
	"testing"
)

func TestLocateJarReportsAMissingConfiguredJar(t *testing.T) {
	t.Setenv(jarEnvironmentVariable, filepath.Join(t.TempDir(), "openfasttrace.jar"))

	_, err := locateJar(context.Background())

	if err == nil || !strings.Contains(err.Error(), jarEnvironmentVariable) {
		t.Fatalf("expected an error naming %s for a missing JAR, got %v", jarEnvironmentVariable, err)
	}
}
