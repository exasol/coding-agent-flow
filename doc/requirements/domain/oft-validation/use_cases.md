# OFT Validation Use Cases

## OFT Requirements Detector

The CAF Caller supplies a Target Project Path. The OFT Requirements Detector reads `req` specification item IDs from an OFT Validation response.

Target Project contributors provide OFT Requirements in Markdown (`.md`) or reStructuredText (`.rst`) files. The OFT Requirements Detector invokes OFT's `trace` command to find them. OFT returns an ASpec document containing OFT objects and their OFT Specification Item IDs.

Scenarios: [OFT Requirements Detector](../../scenarios/oft-validation/oft-requirements-detection/oft-requirements-detector/)

## OFT Tracing Validator

The CAF Caller supplies a Target Project Path. The OFT Tracing Validator validates OFT tracing in the Target Project and returns an OFT Validation response.

Target Project contributors define tracing validation rules through OFT traceability items and links. CAF invokes OFT's `trace` command and decides from the OFT Trace whether tracing is valid: each OFT Specification Item ID is unique, each coverage link reaches exactly one OFT Specification Item in the same revision, each covered item needs the OFT Artifact Type of the covering one, and each needed OFT Artifact Type has a covering item. CAF reports an error when its decision differs from the OFT Result.

Scenarios: [OFT Tracing Validator](../../scenarios/oft-validation/oft-tracing-validation/oft-tracing-validator/)
