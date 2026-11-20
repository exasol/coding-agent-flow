package acceptance

import (
	"context"
	"fmt"
	"strings"
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

// RunFails runs a scenario whose When step returns an error and checks the error with then.
func RunFails[Input, Response any](t *testing.T, files []File, given func(*testing.T, []File) Input, when func(context.Context, Input) (Response, error), then func(error) error) {
	t.Helper()
	_, err := when(t.Context(), given(t, files))
	if err == nil {
		t.Fatal("When: expected an error")
	}
	if err := then(err); err != nil {
		t.Errorf("Then: %v", err)
	}
}

// ExpectError reports an error if err lacks the text.
func ExpectError(text string) func(error) error {
	return func(err error) error {
		if !strings.Contains(err.Error(), text) {
			return fmt.Errorf("error: got %q, want it to contain %q", err, text)
		}
		return nil
	}
}
