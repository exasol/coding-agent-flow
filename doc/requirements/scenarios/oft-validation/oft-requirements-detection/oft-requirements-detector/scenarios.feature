Feature: OFT Requirements Detector
  Identify OFT requirements files in a Target Project.

  @id:scn~oft-requirements-detector~1
  # Covers: req~oft-requirements-detection~1
  # Needs: at
  Scenario Outline: Detect OFT Requirements by file location and extension
    Given a target project path containing the OFT Requirement `<req_id>` in file `<location>`
    When the caller validates OFT tracing for the target project path
    Then the OFT Validation response reports `<found>` for `req` specification item ID `<req_id>`

    Examples:
      | location                             | req_id                              | found |
      | doc/requirements/requirements.rst    | req~valid-doc-requirements-rst~1    | true  |
      | requirements.rst                     | req~valid-root-rst~1                | true  |
      | requirements.md                      | req~valid-root-markdown~1           | true  |
      | doc/requirements/requirements.txt    | req~invalid-doc-requirements-text~1 | false |
      | requirements.txt                     | req~invalid-root-text~1             | false |
