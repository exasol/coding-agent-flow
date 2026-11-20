package acceptance

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestExpectUnchanged(t *testing.T) {
	files := []File{{Path: "README.md", Content: "readme"}, {Path: "doc/a.md", Content: "a"}}
	for name, test := range map[string]struct {
		change func(root string) error
		want   string
	}{
		"unchanged": {func(string) error { return nil }, ""},
		"added":     {func(root string) error { return os.WriteFile(filepath.Join(root, ".lycheecache"), nil, 0o600) }, "added .lycheecache"},
		"modified":  {func(root string) error { return os.WriteFile(filepath.Join(root, "README.md"), []byte("x"), 0o600) }, "modified README.md"},
		"deleted":   {func(root string) error { return os.Remove(filepath.Join(root, "doc", "a.md")) }, "deleted doc/a.md"},
	} {
		targetProject := TargetProject(t, files)
		if err := test.change(targetProject.Path()); err != nil {
			t.Fatal(err)
		}

		err := ExpectUnchanged(targetProject, files)

		if test.want == "" && err != nil {
			t.Errorf("%s: expected no error, got %v", name, err)
		}
		if test.want != "" && (err == nil || !strings.Contains(err.Error(), test.want)) {
			t.Errorf("%s: expected an error with %q, got %v", name, test.want, err)
		}
	}
}
