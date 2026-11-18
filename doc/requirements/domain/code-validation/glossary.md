# Code Validation Glossary

Shared CAF terms are defined in the [Coding Agent Flow Domain Glossary](../glossary.md).

| Term | Meaning |
| --- | --- |
| Code Validation | The checks CAF runs on the content of a Target Project. |
| lychee | The link checker that CAF runs. |
| Link | A reference from a Target Project document to a file, an Anchor, or a URL. |
| Link Target | A project-relative path with an optional Anchor, such as `doc/guide.md#usage`, or a URL. |
| Local Link | A Link to a file or Anchor in the Target Project, such as `doc/guide.md#usage`. lychee checks it in the file system. |
| External Link | A Link to an `http` or `https` URL on any host, including `localhost`, such as `https://exasol.com`. lychee checks it with an HTTP request. |
| Link Scope | The Links the Link Checker checks: Local Links (default), or Local Links and External Links. |
| Anchor | A heading ID in a Markdown document, referenced with `#`. |
| Broken Link | A Link whose target is missing or answers with an error status. A URL that answers `429 Too Many Requests` counts as reachable. |
| Broken Link Status | The lychee status of a Broken Link, such as `File not found`, `Cannot find fragment`, or `Rejected status code: 404 Not Found`. |
| Link Checker | The use case that checks the Links in a Target Project. |
| LinkValidation response | The immutable result of the Link Checker. `valid` is true when all Links resolve. `brokenLinks` lists each Broken Link with source, line, target, and status. |
