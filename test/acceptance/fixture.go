package acceptance

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type File struct{ Path, Content string }

func Run[Parameter, Input, Response any](t *testing.T, parameters []Parameter, given func(*testing.T, Parameter) Input, when func(context.Context, Input) (Response, error), then func(Parameter, Response) error) {
	t.Helper()
	for index, parameter := range parameters {
		t.Run(fmt.Sprintf("scenario-%d", index+1), func(t *testing.T) {
			response, err := when(t.Context(), given(t, parameter))
			if err != nil {
				t.Fatalf("When: %v", err)
			}
			if err := then(parameter, response); err != nil {
				t.Errorf("Then: %v", err)
			}
		})
	}
}

func WriteFiles(t *testing.T, files []File) string {
	t.Helper()
	directory := t.TempDir()
	for _, file := range files {
		path := filepath.Join(directory, filepath.FromSlash(file.Path))
		if !strings.HasPrefix(path, filepath.Clean(directory)+string(filepath.Separator)) {
			t.Fatalf("file %q leaves target project", file.Path)
		}
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(file.Content), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	return directory
}

func Expect(name string, expected, actual bool) error {
	if expected != actual {
		return fmt.Errorf("%s: got %t, want %t", name, actual, expected)
	}
	return nil
}
