package domain

// BrokenLink is a Link that failed the check, with the status that tells why.
type BrokenLink struct {
	link   Link
	status string
}

func NewBrokenLink(link Link, status string) BrokenLink {
	return BrokenLink{link: link, status: status}
}

func (brokenLink BrokenLink) Link() Link {
	return brokenLink.link
}

// Status tells why the Link is broken, for example "File not found" or "404 Not Found".
func (brokenLink BrokenLink) Status() string {
	return brokenLink.status
}
