package service

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/exasol/coding-agent-flow/internal/codevalidation/domain"
	"github.com/exasol/coding-agent-flow/internal/global"
)

// fakeFinder returns fixed results and records the arguments of the last call.
type fakeFinder struct {
	validation domain.LinkValidation
	err        error
	project    *global.TargetProject
	scope      *domain.LinkScope
}

func (finder fakeFinder) FindBrokenLinks(_ context.Context, project global.TargetProject, scope domain.LinkScope) (domain.LinkValidation, error) {
	*finder.project = project
	*finder.scope = scope
	return finder.validation, finder.err
}

// [utest -> req~link-validation~1]
// [utest -> req~external-link-validation~1]
func TestCheckReturnsTheBrokenLinksOfTheTargetProjectInTheScope(t *testing.T) {
	for _, scope := range []domain.LinkScope{domain.LocalLinks, domain.LocalAndExternalLinks} {
		var receivedProject global.TargetProject
		var receivedScope domain.LinkScope
		expected := domain.NewLinkValidation([]domain.BrokenLink{domain.NewBrokenLink(domain.NewLink("README.md", 1, "x.md"), "File not found")})
		project := global.NewTargetProject("project")
		finder := fakeFinder{validation: expected, project: &receivedProject, scope: &receivedScope}

		validation, err := NewLinkChecker(finder, scope).Check(context.Background(), project)

		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(validation, expected) || receivedProject != project || receivedScope != scope {
			t.Fatalf("got %+v for %+v in scope %d, want scope %d", validation, receivedProject, receivedScope, scope)
		}
	}
}

func TestCheckReturnsFinderErrors(t *testing.T) {
	var receivedProject global.TargetProject
	var receivedScope domain.LinkScope
	finder := fakeFinder{err: errors.New("lychee failed"), project: &receivedProject, scope: &receivedScope}

	if _, err := NewLinkChecker(finder, domain.LocalLinks).Check(context.Background(), global.NewTargetProject("project")); err == nil {
		t.Fatal("expected the finder error")
	}
}
