Feature: Link Checker
  Check the Links in the documents of a Target Project.

  @id:scn~link-checker_local-links~1
  # Covers: req~link-validation~1
  # Needs: at
  Scenario Outline: Check local Links by default
    Given a target project path containing `README.md` with the heading `Setup`
    And `doc/guide.md` with the heading `Usage`
    And `<source>` with the Link `<link>`
    And the additional file `<extra file>`
    When the caller checks the links of the target project path
    Then the LinkValidation response has valid `<valid>`
    And the LinkValidation response lists the Broken Links `<broken links>` in `<source>`
    And the target project stays unchanged

    Examples:
      | source    | link               | extra file                                         | valid | broken links        |
      | README.md | doc/guide.md       | -                                                  | true  | []                  |
      | doc/a.md  | ../README.md       | -                                                  | true  | []                  |
      | README.md | #setup             | -                                                  | true  | []                  |
      | README.md | doc/guide.md#usage | -                                                  | true  | []                  |
      | README.md | doc/missing.md     | -                                                  | false | [doc/missing.md]    |
      | README.md | #nope              | -                                                  | false | [README.md#nope]    |
      | README.md | doc/guide.md#nope  | -                                                  | false | [doc/guide.md#nope] |
      | README.md | doc/missing.md     | `lychee.toml` that excludes `README.md` and caches | false | [doc/missing.md]    |

  @id:scn~link-checker_external-links-off-by-default~1
  # Covers: req~link-validation~1
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

  @id:scn~link-checker_external-links-opt-in~1
  # Covers: req~link-validation~1
  # Needs: at
  Scenario Outline: Check external Links on request
    Given a web server that stands in for an external site and serves `/ok`, answers `429 Too Many Requests` for `/rate-limited`, and `404 Not Found` for other paths
    And a target project path containing `<source>` with the Link `<web server><path>`
    When the caller checks the links of the target project path, including external Links
    Then the LinkValidation response has valid `<valid>`
    And the LinkValidation response lists the Broken Links `<broken links>` in `<source>`

    Examples:
      | source      | path          | valid | broken links           |
      | README.md   | /ok           | true  | []                     |
      | README.md   | /rate-limited | true  | []                     |
      | README.md   | /missing      | false | [<web server>/missing] |
      | doc/a.rst   | /ok           | true  | []                     |
      | doc/a.rst   | /missing      | false | [<web server>/missing] |
      | notes.txt   | /ok           | true  | []                     |
      | notes.txt   | /missing      | false | [<web server>/missing] |
      | sitemap.xml | /ok           | true  | []                     |
      | sitemap.xml | /missing      | false | [<web server>/missing] |
