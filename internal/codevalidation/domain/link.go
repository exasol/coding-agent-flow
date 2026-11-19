// Package domain holds the Code Validation domain classes and their rules.
package domain

// Link is a reference in a Target Project document to a file, an anchor, or a URL.
type Link struct {
	source string
	line   int
	target string
}

// NewLink creates a Link in the document at the slash-separated, project-relative path source. target is a
// project-relative path with an optional #anchor, or a URL.
func NewLink(source string, line int, target string) Link {
	return Link{source: source, line: line, target: target}
}

// Source is the project-relative path of the document that contains the Link.
func (link Link) Source() string {
	return link.source
}

func (link Link) Line() int {
	return link.line
}

func (link Link) Target() string {
	return link.target
}
