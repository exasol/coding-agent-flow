Feature: OFT Requirements Detector
  Identify OFT requirements files in a Target Project.

  @id:scn~oft-requirements-rst-in-doc-requirements~1
  # Covers: req~oft-requirements-detection~1
  # Needs: at
  Scenario: Return the req element ID when doc requirements contains an OFT Requirement in reStructuredText
    Given a target project path whose doc requirements directory contains an OFT Requirement in a reStructuredText file
    When the caller validates OFT tracing for the target project path
    Then the OFTValidation response has valid true and req element ID req~valid-doc-requirements-rst~1

  @id:scn~oft-requirements-rst-in-root~1
  # Covers: req~oft-requirements-detection~1
  # Needs: at
  Scenario: Return the req element ID when the target project root contains an OFT Requirement in reStructuredText
    Given a target project path whose root contains an OFT Requirement in a reStructuredText file
    When the caller validates OFT tracing for the target project path
    Then the OFTValidation response has valid true and req element ID req~valid-root-rst~1

  @id:scn~oft-requirements-markdown-in-root~1
  # Covers: req~oft-requirements-detection~1
  # Needs: at
  Scenario: Return the req element ID when the target project root contains a Markdown file
    Given a target project path whose root contains a Markdown file
    When the caller validates OFT tracing for the target project path
    Then the OFTValidation response has valid true and req element ID req~valid-root-markdown~1

  @id:scn~oft-requirements-text-in-doc-requirements~1
  # Covers: req~oft-requirements-detection~1
  # Needs: at
  Scenario: Return no req element IDs when doc requirements contains a text file
    Given a target project path whose doc requirements directory contains a text file with an OFT Requirement ID
    When the caller validates OFT tracing for the target project path
    Then the OFTValidation response has valid true and no req element IDs

  @id:scn~oft-requirements-text-in-root~1
  # Covers: req~oft-requirements-detection~1
  # Needs: at
  Scenario: Return no req element IDs when the target project root contains a text file
    Given a target project path whose root contains a text file with an OFT Requirement ID
    When the caller validates OFT tracing for the target project path
    Then the OFTValidation response has valid true and no req element IDs
