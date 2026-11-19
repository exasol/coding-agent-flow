package lychee

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func fakeBinary(t *testing.T, output string, exitCode int) string {
	t.Helper()
	binary := filepath.Join(t.TempDir(), "lychee")
	script := fmt.Sprintf("#!/bin/sh\nprintf '%%s\\n' '%s'\nexit %d\n", output, exitCode)
	if err := os.WriteFile(binary, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv(binaryEnvironmentVariable, binary)
	return binary
}

func TestLocateBinaryAcceptsTheSupportedVersion(t *testing.T) {
	expected := fakeBinary(t, "lychee "+lycheeVersion, 0)

	binary, err := locateBinary(context.Background())

	if err != nil || binary != expected {
		t.Fatalf("got %q, %v; want %q", binary, err, expected)
	}
}

func TestLocateBinaryRefusesOtherVersionsAndUnexpectedOutput(t *testing.T) {
	for name, test := range map[string]struct {
		output   string
		exitCode int
		want     string
	}{
		"other version":     {"lychee 0.18.1", 0, "reports version 0.18.1"},
		"no lychee version": {"garbage", 0, "reports no lychee version"},
		"failing binary":    {"lychee " + lycheeVersion, 1, "--version"},
	} {
		fakeBinary(t, test.output, test.exitCode)

		_, err := locateBinary(context.Background())

		if err == nil || !strings.Contains(err.Error(), "E-CAF-4") || !strings.Contains(err.Error(), test.want) {
			t.Errorf("%s: expected E-CAF-4 with %q, got %v", name, test.want, err)
		}
	}
}

func TestParseVersion(t *testing.T) {
	for output, want := range map[string]string{
		"lychee 0.24.2\n":      "0.24.2",
		"  lychee 1.0.0  ":     "1.0.0",
		"lychee 0.24.2-beta":   "",
		"lychee-lib 0.24.2":    "",
		"lychee 0.24.2\nextra": "",
		"":                     "",
	} {
		version, found := parseVersion([]byte(output))
		if version != want || found != (want != "") {
			t.Errorf("%q: got %q, %t; want %q", output, version, found, want)
		}
	}
}

func TestLocateBinaryReportsAMissingConfiguredBinary(t *testing.T) {
	t.Setenv(binaryEnvironmentVariable, filepath.Join(t.TempDir(), "lychee"))

	_, err := locateBinary(context.Background())

	if err == nil || !strings.Contains(err.Error(), binaryEnvironmentVariable) || !strings.Contains(err.Error(), "E-CAF-4") {
		t.Fatalf("expected E-CAF-4 naming %s, got %v", binaryEnvironmentVariable, err)
	}
}

func releaseArchiveWith(t *testing.T, name, content string) []byte {
	t.Helper()
	var archive bytes.Buffer
	compressed := gzip.NewWriter(&archive)
	files := tar.NewWriter(compressed)
	if err := files.WriteHeader(&tar.Header{Name: name, Mode: 0o755, Size: int64(len(content)), Typeflag: tar.TypeReg}); err != nil {
		t.Fatal(err)
	}
	if _, err := files.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	if err := files.Close(); err != nil {
		t.Fatal(err)
	}
	if err := compressed.Close(); err != nil {
		t.Fatal(err)
	}
	return archive.Bytes()
}

func checksumOf(content []byte) string {
	sum := sha256.Sum256(content)
	return hex.EncodeToString(sum[:])
}

func TestInstallBinaryWritesTheExecutableFromTheArchive(t *testing.T) {
	archive := releaseArchiveWith(t, "lychee-x86_64-unknown-linux-musl/lychee", "binary")
	binary := filepath.Join(t.TempDir(), "cache", "lychee")

	if err := installBinary(archive, checksumOf(archive), binary); err != nil {
		t.Fatal(err)
	}
	content, err := os.ReadFile(binary)
	if err != nil || string(content) != "binary" {
		t.Fatalf("got %q, %v", content, err)
	}
	if info, err := os.Stat(binary); err != nil || info.Mode().Perm()&0o100 == 0 {
		t.Fatalf("expected an executable, got %v, %v", info, err)
	}
}

func TestInstallBinaryRefusesAnArchiveWithAnotherChecksum(t *testing.T) {
	archive := releaseArchiveWith(t, "lychee", "binary")
	binary := filepath.Join(t.TempDir(), "lychee")

	if err := installBinary(archive, checksumOf([]byte("other")), binary); err == nil || !strings.Contains(err.Error(), "checksum") {
		t.Fatalf("expected a checksum error, got %v", err)
	}
	if _, err := os.Stat(binary); !os.IsNotExist(err) {
		t.Fatal("expected no binary after a checksum mismatch")
	}
}

func TestInstallBinaryRefusesAnArchiveWithoutLychee(t *testing.T) {
	archive := releaseArchiveWith(t, "README.md", "text")

	if err := installBinary(archive, checksumOf(archive), filepath.Join(t.TempDir(), "lychee")); err == nil {
		t.Fatal("expected an error for an archive without lychee")
	}
}

// fakeLychee configures as LYCHEE_BINARY a shell script that reports the supported version for --version. For a
// check, it records its arguments and working directory, prints stdout and stderr, and exits with exitCode.
