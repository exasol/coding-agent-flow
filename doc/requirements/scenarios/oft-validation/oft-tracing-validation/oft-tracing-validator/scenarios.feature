Feature: OFT Tracing Validator
  Validate OFT tracing in a Target Project.

  @id:scn~oft-tracing-validator_valid~1
  # Covers: req~oft-tracing-validation~1
  # Needs: at
  Scenario: Return valid when an OFT Requirement covers the OFT Feature that needs it
    Given a target project path containing `features.md` with the OFT Feature `feat~valid-fixture-feature~1` that needs an OFT Requirement
    And `requirements.md` with the OFT Requirement `req~valid-fixture-requirement~1` that covers `feat~valid-fixture-feature~1`
    When the caller validates OFT tracing for the target project path
    Then the OFT Validation response has valid `true`

  @id:scn~oft-tracing-validator_invalid~1
  # Covers: req~oft-tracing-validation~1
  # Needs: at
  Scenario: Return invalid when an OFT Feature that needs an OFT Requirement stays uncovered
    Given a target project path containing `features.md` with the OFT Feature `feat~invalid-fixture-feature~1` that needs an OFT Requirement
    And `feat~invalid-fixture-feature~1` stays uncovered
    When the caller validates OFT tracing for the target project path
    Then the OFT Validation response has valid `false`

  @id:scn~oft-tracing-validator_requirement-covers-other-feature~1
  # Covers: req~oft-tracing-validation~1
  # Needs: at
  Scenario: Return invalid when the OFT Requirement covers another OFT Feature than the one that needs it
    Given a target project path containing `features.md` with the OFT Feature `feat~needed-fixture-feature~1` that needs an OFT Requirement
    And `requirements.md` with the OFT Requirement `req~other-feature-fixture-requirement~1` that covers `feat~other-fixture-feature~1`
    When the caller validates OFT tracing for the target project path
    Then the OFT Validation response has valid `false`

  @id:scn~oft-tracing-validator_requirement-covers-other-revision~1
  # Covers: req~oft-tracing-validation~1
  # Needs: at
  Scenario Outline: Return invalid when the OFT Requirement covers another revision of the OFT Feature
    Given a target project path containing `features.md` with the OFT Feature `<feature>` that needs an OFT Requirement
    And `requirements.md` with the OFT Requirement `req~revision-fixture-requirement~1` that covers `<covered>`
    When the caller validates OFT tracing for the target project path
    Then the OFT Validation response has valid `false`

    Examples:
      | feature                        | covered                        |
      | feat~revised-fixture-feature~2 | feat~revised-fixture-feature~1 |
      | feat~revised-fixture-feature~1 | feat~revised-fixture-feature~2 |

  @id:scn~oft-tracing-validator_unwanted-coverage~1
  # Covers: req~oft-tracing-validation~1
  # Needs: at
  Scenario: Return invalid when an OFT Specification Item covers an OFT Feature that needs another OFT Artifact Type
    Given a target project path containing `features.md` with the OFT Feature `feat~unwanted-fixture-feature~1` that needs an OFT Requirement
    And `requirements.md` with the OFT Requirement `req~unwanted-fixture-requirement~1` that covers `feat~unwanted-fixture-feature~1`
    And `impl.go` with the tag `[impl->feat~unwanted-fixture-feature~1]`
    When the caller validates OFT tracing for the target project path
    Then the OFT Validation response has valid `false`

  @id:scn~oft-tracing-validator_duplicate-item~1
  # Covers: req~oft-tracing-validation~1
  # Needs: at
  Scenario: Return invalid when an OFT Specification Item ID appears twice
    Given a target project path containing `features.md` and `more_features.md`, each with the OFT Feature `feat~duplicate-fixture-feature~1` that needs an OFT Requirement
    And `requirements.md` with the OFT Requirement `req~duplicate-fixture-requirement~1` that covers `feat~duplicate-fixture-feature~1`
    When the caller validates OFT tracing for the target project path
    Then the OFT Validation response has valid `false`

  @id:scn~oft-tracing-validator_implementation-coverage~1
  # Covers: req~oft-tracing-validation~1
  # Needs: at
  Scenario Outline: Validate coverage of an OFT Requirement by an implementation tag
    Given a target project path containing `features.md` with the OFT Feature `feat~chain-fixture-feature~1` that needs an OFT Requirement
    And `requirements.md` with the OFT Requirement `req~chain-fixture-requirement~1` that covers `feat~chain-fixture-feature~1` and needs `impl`
    And <source>
    When the caller validates OFT tracing for the target project path
    Then the OFT Validation response has valid `<valid>`

    Examples:
      | source                                                                 | valid |
      | `src/main.go` with the tag `[impl->req~chain-fixture-requirement~1]`   | true  |
      | no implementation tag                                                  | false |
      | `src/main.go` with the tag `[impl->req~missing-fixture-requirement~1]` | false |

  @id:scn~oft-tracing-validator_no-specification-items~1
  # Covers: req~oft-tracing-validation~1
  # Needs: at
  Scenario Outline: Return valid without OFT Specification Item IDs when the target project has no OFT Specification Items
    Given a target project path containing <content>
    When the caller validates OFT tracing for the target project path
    Then the OFT Validation response has valid `true`
    And the OFT Validation response has no OFT Specification Item IDs

    Examples:
      | content                                          |
      | no files                                         |
      | only `README.md` without OFT Specification Items |

  @id:scn~oft-tracing-validator_nested-files~1
  # Covers: req~oft-tracing-validation~1
  # Needs: at
  Scenario: Return valid when the OFT Specification Items are in nested directories
    Given a target project path containing `doc/deep/er/features.md` with the OFT Feature `feat~nested-fixture-feature~1` that needs an OFT Requirement
    And `requirements.md` with the OFT Requirement `req~nested-fixture-requirement~1` that covers `feat~nested-fixture-feature~1`
    When the caller validates OFT tracing for the target project path
    Then the OFT Validation response has valid `true`

  @id:scn~oft-tracing-validator_target-project-path-no-directory~1
  # Covers: req~oft-tracing-validation~1
  # Needs: at
  Scenario Outline: Return an error when the target project path is no directory
    Given a target project path that is <path>
    When the caller validates OFT tracing for the target project path
    Then the caller receives an error containing `Target Project`

    Examples:
      | path           |
      | missing        |
      | a regular file |

  @id:scn~oft-tracing-validator_symlinked-target-project~1
  # Covers: req~oft-tracing-validation~1
  # Needs: at
  Scenario: Validate the directory that a symbolic link as target project path points to
    Given a target project path that is a symbolic link to a directory containing `features.md` with the OFT Feature `feat~linked-fixture-feature~1` that needs an OFT Requirement
    And `requirements.md` with the OFT Requirement `req~linked-fixture-requirement~1` that covers `feat~linked-fixture-feature~1`
    When the caller validates OFT tracing for the target project path
    Then the OFT Validation response has valid `true`
    And the OFT Validation response reports `true` for `feat` specification item ID `feat~linked-fixture-feature~1`

  @id:scn~oft-tracing-validator_unreadable-file~1
  # Covers: req~oft-tracing-validation~1
  # Needs: at
  Scenario: Return an error when OFT can't read a file of the target project
    Given a target project path containing `features.md` with the OFT Feature `feat~unreadable-fixture-feature~1` that needs an OFT Requirement
    And an unreadable `requirements.md`
    When the caller validates OFT tracing for the target project path
    Then the caller receives an error containing `Error reading`

  @id:scn~oft-tracing-validator_unusable-jar~1
  # Covers: req~oft-tracing-validation~1
  # Needs: at
  Scenario Outline: Return an error when OPENFASTTRACE_JAR is unusable
    Given `OPENFASTTRACE_JAR` points to <jar>
    And a target project path containing `features.md` with the OFT Feature `feat~jar-fixture-feature~1` that needs an OFT Requirement
    When the caller validates OFT tracing for the target project path
    Then the caller receives an error containing `OPENFASTTRACE_JAR`

    Examples:
      | jar                   |
      | a missing file        |
      | a file that is no JAR |

  @id:scn~oft-tracing-validator_read-only~1
  # Covers: req~oft-tracing-validation~1
  # Needs: at
  Scenario: Leave the target project unchanged
    Given a target project path containing `features.md` with the OFT Feature `feat~read-only-fixture-feature~1` that needs an OFT Requirement
    And `requirements.md` with the OFT Requirement `req~read-only-fixture-requirement~1` that covers `feat~read-only-fixture-feature~1`
    When the caller validates OFT tracing for the target project path
    Then the target project is unchanged
