package acceptance

import (
	"context"
	"fmt"
	"testing"
)

func Run[Input, Response any](t *testing.T, files []File, given func(*testing.T, []File) Input, when func(context.Context, Input) (Response, error), then func(Response) error) {
	t.Helper()
	response, err := when(t.Context(), given(t, files))
	if err != nil {
		t.Fatalf("When: %v", err)
	}
	if err := then(response); err != nil {
		t.Errorf("Then: %v", err)
	}
}

func Expect(name string, expected, actual bool) error {
	if expected != actual {
		return fmt.Errorf("%s: got %t, want %t", name, actual, expected)
	}
	return nil
}
