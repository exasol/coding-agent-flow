# Code Validation Use Cases

## Link Checker

The CAF Caller supplies a Target Project Path. The Link Checker checks the Links in the Target Project documents and returns a LinkValidation response.

| Extension | Read as | Local Links (default) | External Links (on request) |
| --- | --- | --- | --- |
| `.md` | Markdown | Files and Anchors | URLs |
| `.rst` | Plain text | – | URLs |
| `.txt` | Plain text | – | URLs |
| `.xml` | Sitemap | – | URLs in `<loc>` elements |

The Link Checker skips:

* hidden files and directories, such as `.github/`;
* files listed in `.ignore`, and in git repositories, files listed in `.gitignore`;
* Links listed in `.lycheeignore`.

The Link Checker leaves the Target Project unchanged and ignores link-checker configuration files in it.

Scenarios: [Link Checker](../../scenarios/code-validation/link-validation/link-checker/)

## External Link Checker

The CAF Caller supplies a Target Project Path and requests External Links. The External Link Checker checks the Local Links and the External Links and returns a LinkValidation response. It counts a URL that answers `429 Too Many Requests` as valid. The table above shows the External Links per extension.

Scenarios: [External Link Checker](../../scenarios/code-validation/external-link-validation/external-link-checker/)
