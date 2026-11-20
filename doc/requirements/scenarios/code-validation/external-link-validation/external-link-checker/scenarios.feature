Feature: External Link Checker
  Check the External Links in the documents of a Target Project on request.

  @id:scn~external-link-checker_off-by-default~1
  # Covers: req~external-link-validation~1
  # Needs: at
  Scenario Outline: Skip external Links by default
    Given a web server that answers `404 Not Found` for `/missing`
    And a target project path containing `<source>` with the Link `<web server>/missing`
    When the caller checks the links of the target project path
    Then the LinkValidation response has valid `true`
    And the web server receives 0 requests

    Examples:
      | source      |
      | README.md   |
      | doc/a.rst   |
      | notes.txt   |
      | sitemap.xml |

  @id:scn~external-link-checker_valid-links~1
  # Covers: req~external-link-validation~1
  # Needs: at
  Scenario Outline: Return valid for reachable External Links on request
    Given a web server that stands in for an external site, serves `/ok`, and answers `429 Too Many Requests` for `/rate-limited`
    And a target project path containing `<source>` with the Link `<web server><path>`
    When the caller checks the links of the target project path, including External Links
    Then the LinkValidation response has valid `true`

    Examples:
      | source      | path          |
      | README.md   | /ok           |
      | README.md   | /rate-limited |
      | doc/a.rst   | /ok           |
      | notes.txt   | /ok           |
      | sitemap.xml | /ok           |

  @id:scn~external-link-checker_broken-link~1
  # Covers: req~external-link-validation~1
  # Needs: at
  Scenario Outline: Report a broken External Link on request
    Given a web server that stands in for an external site and answers `404 Not Found` for `/missing`
    And a target project path containing `<source>` with the Link `<web server>/missing`
    When the caller checks the links of the target project path, including External Links
    Then the LinkValidation response has valid `false`
    And the LinkValidation response lists the Broken Links `[<web server>/missing]` in `<source>`

    Examples:
      | source      |
      | README.md   |
      | doc/a.rst   |
      | notes.txt   |
      | sitemap.xml |
