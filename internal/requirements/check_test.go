package requirements

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// [utest~requirements-document-detection~1 -> req~requirements-document-detection~1]
func TestHasDocumentReturnsTrueForSupportedDirectFile(t *testing.T) {
	for _, extension := range []string{".md", ".rst"} {
		t.Run(extension, func(t *testing.T) {
			projectRoot := t.TempDir()
			requirementsDirectory := filepath.Join(projectRoot, "doc", "requirements")
			require.NoError(t, os.MkdirAll(requirementsDirectory, 0o755))

			requirementFile := filepath.Join(requirementsDirectory, "example"+extension)
			require.NoError(t, os.WriteFile(requirementFile, nil, 0o600))

			found, err := HasDocument(projectRoot)
			require.NoError(t, err)
			assert.True(t, found, "expected to find a requirements document for extension %s", extension)
		})
	}
}

func TestHasDocumentReturnsFalseWhenDirectoryIsMissing(t *testing.T) {
	projectRoot := t.TempDir()

	found, err := HasDocument(projectRoot)

	require.NoError(t, err)
	assert.False(t, found)
}

func TestHasDocumentReturnsFalseWhenDirectoryIsEmpty(t *testing.T) {
	projectRoot := t.TempDir()
	requirementsDirectory := filepath.Join(projectRoot, "doc", "requirements")
	require.NoError(t, os.MkdirAll(requirementsDirectory, 0o755))

	found, err := HasDocument(projectRoot)

	require.NoError(t, err)
	assert.False(t, found)
}

func TestHasDocumentReturnsFalseForUnsupportedExtension(t *testing.T) {
	projectRoot := t.TempDir()
	requirementsDirectory := filepath.Join(projectRoot, "doc", "requirements")
	require.NoError(t, os.MkdirAll(requirementsDirectory, 0o755))
	require.NoError(t, os.WriteFile(
		filepath.Join(requirementsDirectory, "example.txt"),
		nil,
		0o600,
	))

	found, err := HasDocument(projectRoot)

	require.NoError(t, err)
	assert.False(t, found)
}
