# OFT Validation Use Cases

## OFT Requirements Detector

The CAF Caller supplies a Target Project Path. The OFT Requirements Detector reads `req` element IDs from an OFTValidation response.

Target Project contributors provide OFT Requirements in Markdown (`.md`) or reStructuredText (`.rst`) files. The OFT Requirements Detector invokes OFT's `trace` command to find them. OFT returns an ASpec document containing OFT objects and their OFT Element IDs.

Scenarios: [OFT Requirements Detector](../../scenarios/oft-validation/oft-requirements-detection/oft-requirements-detector/)

## OFT Tracing Validator

The CAF Caller supplies a Target Project Path. The OFT Tracing Validator validates OFT tracing in the Target Project and returns an OFTValidation response.

Target Project contributors define tracing validation rules through OFT traceability items and links. CAF invokes OFT's `trace` command to validate them.

Scenarios: [OFT Tracing Validator](../../scenarios/oft-validation/oft-tracing-validation/oft-tracing-validator/)
