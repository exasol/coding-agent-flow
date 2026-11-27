# Validation Entry Point Design

## Configuration File

`dsn~configuration-file~1`

The adapter `internal/validationentrypoint/adapter/yamlconfig` reads `.caf.yaml` as a regular file of at most 64 KiB in the Target Project root.

Covers:

- `req~validation-selection~1`

Needs: impl, utest

Rationale: `.caf.yaml` belongs to the Target Project, so CAF treats it as untrusted input.

- Regular file: CAF checks the file type before it opens the file and checks the opened file again. A symbolic link cannot point CAF at a file outside the Target Project, and a FIFO or device cannot block the read.
- 64 KiB: the file holds a version and a list of Validation Names, which stays well below 1 KiB even with many validators. The limit leaves more than 50 times that headroom and bounds the memory and YAML parsing a broken or hostile file can cause. A larger file fails with `E-CAF-11`, which names the limit.
