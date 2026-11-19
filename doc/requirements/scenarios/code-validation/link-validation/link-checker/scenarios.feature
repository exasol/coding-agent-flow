Feature: Link Checker
  Check the Links in the documents of a Target Project.

  @id:scn~link-checker_local-links~1
  # Covers: req~link-validation~1
  # Needs: at
  Scenario Outline: Check local Links
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
