package lychee

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"

	exaerror "github.com/exasol/error-reporting-go"
)

const lycheeVersion = "0.24.2"

const binaryEnvironmentVariable = "LYCHEE_BINARY"

const releaseURL = "https://github.com/lycheeverse/lychee/releases/download/lychee-v" + lycheeVersion + "/"

// releaseArchive is a lychee release archive for one platform with its SHA-256 checksum from the release page.
type releaseArchive struct {
	target   string
	checksum string
}

var releaseArchives = map[string]releaseArchive{
	"linux/amd64":  {"x86_64-unknown-linux-musl", "73657a111819a30c47c08352896796f23d64e4eb2b3ed39b6d32149241566fc5"},
	"linux/arm64":  {"aarch64-unknown-linux-musl", "5d0b0e3aeab240f41920c633a6eaf97599be6eedda034b36e858ede7dba5e535"},
	"darwin/amd64": {"x86_64-apple-darwin", "887503a9cff667d322b8d0892b40bf49976eb9507af8483220a3706cdad55978"},
	"darwin/arm64": {"aarch64-apple-darwin", "c9d3740ea2d891854d37116c9fba840f37b6e7c89d330e7db84ac333631c4977"},
}

// versionPattern matches the output of lychee --version, such as "lychee 0.24.2".
var versionPattern = regexp.MustCompile(`^lychee (\d+\.\d+\.\d+)$`)

// versionTimeout limits the time lychee --version may take.
const versionTimeout = 10 * time.Second

// maximumDiagnosticLength limits the program output that an error message quotes.
const maximumDiagnosticLength = 200

// locateBinary returns the lychee binary and verifies that it reports the supported version.
func locateBinary(ctx context.Context) (string, error) {
	binary, err := findBinary(ctx)
	if err != nil {
		return "", err
	}
	if err := verifyVersion(ctx, binary); err != nil {
		return "", binaryError(err)
	}
	return binary, nil
}

// findBinary returns the lychee binary from LYCHEE_BINARY if set, else from the user cache directory, where it
// downloads the release archive on first use and verifies its checksum.
func findBinary(ctx context.Context) (string, error) {
	if configured := os.Getenv(binaryEnvironmentVariable); configured != "" {
		if _, err := os.Stat(configured); err != nil {
			return "", binaryError(fmt.Errorf("inspect lychee binary from %s: %w", binaryEnvironmentVariable, err))
		}
		return configured, nil
	}
	cacheDirectory, err := os.UserCacheDir()
	if err != nil {
		return "", binaryError(fmt.Errorf("determine the cache directory: %w", err))
	}
	binary := filepath.Join(cacheDirectory, "coding-agent-flow", "lychee", lycheeVersion, "lychee")
	if _, err := os.Stat(binary); err == nil {
		return binary, nil
	} else if !errors.Is(err, os.ErrNotExist) {
		return "", binaryError(fmt.Errorf("inspect lychee binary: %w", err))
	}
	archive, supported := releaseArchives[runtime.GOOS+"/"+runtime.GOARCH]
	if !supported {
		return "", binaryError(fmt.Errorf("no lychee release archive for %s/%s", runtime.GOOS, runtime.GOARCH))
	}
	if err := download(ctx, releaseURL+"lychee-"+archive.target+".tar.gz", archive.checksum, binary); err != nil {
		return "", binaryError(err)
	}
	return binary, nil
}

// verifyVersion requires lychee --version to report lycheeVersion, the version whose JSON report the adapter reads.
func verifyVersion(ctx context.Context, binary string) error {
	ctx, cancel := context.WithTimeout(ctx, versionTimeout)
	defer cancel()
	output, err := exec.CommandContext(ctx, binary, "--version").Output()
	if err != nil {
		return fmt.Errorf("run %s --version: %w", binary, err)
	}
	version, found := parseVersion(output)
	if !found {
		return fmt.Errorf("%s --version reports no lychee version: %q", binary, truncate(output))
	}
	if version != lycheeVersion {
		return fmt.Errorf("%s reports version %s", binary, version)
	}
	return nil
}

// parseVersion reads the version from the output of lychee --version.
func parseVersion(output []byte) (string, bool) {
	match := versionPattern.FindStringSubmatch(strings.TrimSpace(string(output)))
	if match == nil {
		return "", false
	}
	return match[1], true
}

func truncate(output []byte) string {
	if len(output) > maximumDiagnosticLength {
		return string(output[:maximumDiagnosticLength]) + "..."
	}
	return string(output)
}

// download fetches the release archive, verifies its checksum, and writes the lychee binary from it to binary.
func download(ctx context.Context, archiveURL, checksum, binary string) error {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, archiveURL, nil)
	if err != nil {
		return err
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return fmt.Errorf("download lychee %s: %w", lycheeVersion, err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("download lychee %s: %s", lycheeVersion, response.Status)
	}
	archive, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("download lychee %s: %w", lycheeVersion, err)
	}
	return installBinary(archive, checksum, binary)
}

func installBinary(archive []byte, checksum, binary string) error {
	if sum := sha256.Sum256(archive); hex.EncodeToString(sum[:]) != checksum {
		return fmt.Errorf("lychee archive has checksum %x, want %s", sum, checksum)
	}
	content, err := extractBinary(archive)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(binary), 0o755); err != nil {
		return err
	}
	temporary, err := os.CreateTemp(filepath.Dir(binary), "lychee-*")
	if err != nil {
		return err
	}
	defer os.Remove(temporary.Name())
	if _, err := temporary.Write(content); err != nil {
		temporary.Close()
		return err
	}
	if err := temporary.Close(); err != nil {
		return err
	}
	if err := os.Chmod(temporary.Name(), 0o755); err != nil {
		return err
	}
	return os.Rename(temporary.Name(), binary)
}

// extractBinary returns the file named lychee from a .tar.gz release archive.
func extractBinary(archive []byte) ([]byte, error) {
	compressed, err := gzip.NewReader(bytes.NewReader(archive))
	if err != nil {
		return nil, fmt.Errorf("read lychee archive: %w", err)
	}
	files := tar.NewReader(compressed)
	for {
		header, err := files.Next()
		if errors.Is(err, io.EOF) {
			return nil, errors.New("lychee archive contains no lychee binary")
		}
		if err != nil {
			return nil, fmt.Errorf("read lychee archive: %w", err)
		}
		if header.Typeflag == tar.TypeReg && filepath.Base(header.Name) == "lychee" {
			return io.ReadAll(files)
		}
	}
}

func binaryError(err error) error {
	return exaerror.New("E-CAF-4").
		Message("Failed to provide lychee {{version}}: {{error|uq}}.").
		Parameter("version", lycheeVersion).
		Parameter("error", err).
		Mitigation("Set " + binaryEnvironmentVariable + " to the path of a lychee " + lycheeVersion + " binary, or unset it so that CAF downloads lychee " + lycheeVersion + ".")
}
