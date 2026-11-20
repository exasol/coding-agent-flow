package domain

// LinkScope selects which Links the Link Checker checks.
type LinkScope int

const (
	// LocalLinks are Links to files and Anchors in the Target Project. lychee checks them offline.
	LocalLinks LinkScope = iota
	// LocalAndExternalLinks are the Local Links and the External Links.
	LocalAndExternalLinks
)
