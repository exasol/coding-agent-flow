# Code Validation Use Cases

## Link Checker

The CAF Caller supplies a Target Project Path. The Link Checker checks the Links in the Target Project documents with [lychee](https://github.com/lycheeverse/lychee) and returns a LinkValidation response.

| Extension | lychee reads it as | Local Links (default) | External Links (on request) |
| --- | --- | --- | --- |
| `.md` | Markdown | Files and Anchors | URLs |
| `.rst` | Plain text | – | URLs |
| `.txt` | Plain text | – | URLs |
| `.xml` | Sitemap | – | URLs in `<loc>` elements |

lychee skips:

* hidden files and directories, such as `.github/`;
* files listed in `.ignore`, and in git repositories, files listed in `.gitignore`;
* Links listed in `.lycheeignore`.

CAF runs lychee with its own configuration and leaves the Target Project unchanged. lychee ignores a `lychee.toml` in the Target Project.

Scenarios: [Link Checker](../../scenarios/code-validation/link-validation/link-checker/)

## External Link Checker

The CAF Caller supplies a Target Project Path and requests External Links. The External Link Checker checks the Local Links and the External Links and returns a LinkValidation response. It counts a URL that answers `429 Too Many Requests` as valid. Without the request, the Link Checker checks Local Links offline. The table above shows the External Links per extension.

Scenarios: [External Link Checker](../../scenarios/code-validation/external-link-validation/external-link-checker/)
