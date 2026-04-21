package template

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	clientconfig "client/config"
)

/*
Category Partition Methodology (Steps 1–4)

Step 1 — Identify Parameters and Environment Conditions
Inputs:
- verifyAndHandlePath(path string)
- Create(name string)
- Remove(name string)

Outputs:
- (string, error) for Create/Remove
- error for verifyAndHandlePath

Relevant environment/state:
- Filesystem state (directory exists/missing, file exists/missing)
- client configuration path via clientconfig.DefaultPath
- Path flavor (absolute path vs name-like path)
- Name normalization performed by system.NormalizeNamePathExploit

Step 2 — Define Categories
A) Path existence category (verifyAndHandlePath):
- A1: path does not exist (valid create path)
- A2: path already exists (special/error path according to current implementation)
- A3: invalid/unwritable path (error path)

B) Create(name) category:
- B1: normal valid exploit name
- B2: name with extension/special but valid
- B3: empty name (invalid)
- B4: exploit directory missing before call (boundary/environment)

C) Remove(name) category:
- C1: absolute path to existing file
- C2: absolute path to missing file
- C3: name-based path (relative/name) — current implementation builds from empty base and should fail
- C4: base config path missing/already existing (via verifyAndHandlePath behavior)

Step 3 — Define Constraints
- Create/Remove read clientconfig.DefaultPath; tests must isolate by overriding it with a temporary directory.
- Remove first calls verifyAndHandlePath(clientconfig.DefaultPath):
  - If DefaultPath already exists, function returns an error by current code path.
  - Therefore Remove success path requires DefaultPath missing initially.
- For absolute-path remove tests, ensure Remove receives an absolute path (to trigger IsPath true branch).
- Invalid/unwritable path is platform-dependent; avoid brittle test relying on root permissions.

Step 4 — Generate Test Frames
1) verifyAndHandlePath with missing path -> expect nil and directory created (normal)
2) verifyAndHandlePath with existing path -> expect non-nil error (special/error)
3) Create with valid name and missing exploits dir -> expect success, file created, template content present (normal+boundary)
4) Create with empty name -> expect error (error)
5) Remove with existing absolute path and missing DefaultPath -> expect success and file removed (normal)
6) Remove with missing absolute path and missing DefaultPath -> expect error contains "does not exist" (error)
7) Remove with name-only input -> expect error due to current branch behavior (special/error)
*/

func withTempDefaultPath(t *testing.T) string {
	t.Helper()

	orig := clientconfig.DefaultPath
	tmp := t.TempDir()
	clientconfig.DefaultPath = tmp
	t.Cleanup(func() {
		clientconfig.DefaultPath = orig
	})

	return tmp
}

func TestVerifyAndHandlePath_should_create_directory_when_path_missing(t *testing.T) {
	t.Parallel()

	base := t.TempDir()
	target := filepath.Join(base, "missing", "nested")

	// Category A1: path does not exist
	if err := verifyAndHandlePath(target); err != nil {
		t.Fatalf("verifyAndHandlePath() error = %v; want nil", err)
	}

	info, err := os.Stat(target)
	if err != nil {
		t.Fatalf("expected directory to exist: %v", err)
	}
	if !info.IsDir() {
		t.Fatalf("expected %q to be a directory", target)
	}
}

func TestCreate_should_create_template_file_when_valid_name_and_dir_missing(t *testing.T) {
	t.Parallel()

	defaultPath := withTempDefaultPath(t)
	exploitName := "my_exploit.py"

	// Ensure exploits directory is missing before call (Category B4 + B1)
	exploitsDir := filepath.Join(defaultPath, "exploits")
	if _, err := os.Stat(exploitsDir); !os.IsNotExist(err) {
		t.Fatalf("expected exploits dir to be missing before test")
	}

	msg, err := Create(exploitName)
	if err != nil {
		t.Fatalf("Create() error = %v; want nil", err)
	}
	if !strings.Contains(msg, "Exploit file created successfully at") {
		t.Fatalf("unexpected success message: %q", msg)
	}

	createdPath := filepath.Join(exploitsDir, exploitName)
	data, readErr := os.ReadFile(createdPath)
	if readErr != nil {
		t.Fatalf("expected created exploit file at %q: %v", createdPath, readErr)
	}

	content := string(data)
	// Covers template integrity
	mustContain := []string{
		"#!/usr/bin/env python3",
		"from cookiefarm import exploit_manager",
		"@exploit_manager",
		"def exploit(ip, port, name_service, flag_ids: list):",
	}
	for _, token := range mustContain {
		if !strings.Contains(content, token) {
			t.Fatalf("created template missing token %q", token)
		}
	}
}

func TestCreate_should_return_error_when_name_empty(t *testing.T) {
	t.Parallel()

	withTempDefaultPath(t)

	// Category B3: empty/invalid name
	_, err := Create("")
	// Error text may vary based on NormalizeNamePathExploit implementation
	if !strings.Contains(err.Error(), "exploit name cannot be empty") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRemove_should_remove_file_when_absolute_path_exists_and_default_path_missing(t *testing.T) {
	t.Parallel()

	parent := t.TempDir()
	// Put DefaultPath under a missing child so verifyAndHandlePath can create it (required by current code)
	defaultPath := filepath.Join(parent, "cfg_missing")
	orig := clientconfig.DefaultPath
	clientconfig.DefaultPath = defaultPath
	t.Cleanup(func() {
		clientconfig.DefaultPath = orig
	})

	// Create target file outside DefaultPath and pass absolute path (Category C1)
	targetDir := t.TempDir()
	targetFile := filepath.Join(targetDir, "to_remove.py")
	if err := os.WriteFile(targetFile, []byte("print('x')"), 0o644); err != nil {
		t.Fatalf("failed to create target file: %v", err)
	}

	msg, err := Remove(targetFile)
	if err != nil {
		t.Fatalf("Remove(abs-existing) error = %v; want nil", err)
	}
	if !strings.Contains(msg, "Exploit file removed successfully") {
		t.Fatalf("unexpected success message: %q", msg)
	}

	if _, statErr := os.Stat(targetFile); !os.IsNotExist(statErr) {
		t.Fatalf("expected file to be removed, stat err = %v", statErr)
	}
}

func TestRemove_should_return_error_when_absolute_path_missing(t *testing.T) {
	t.Parallel()

	parent := t.TempDir()
	defaultPath := filepath.Join(parent, "cfg_missing")
	orig := clientconfig.DefaultPath
	clientconfig.DefaultPath = defaultPath
	t.Cleanup(func() {
		clientconfig.DefaultPath = orig
	})

	// Category C2: absolute path missing
	missingAbs := filepath.Join(t.TempDir(), "nope.py")
	_, err := Remove(missingAbs)
	if err == nil {
		t.Fatalf("Remove(abs-missing) error = nil; want non-nil")
	}
	if !strings.Contains(err.Error(), "exploit file does not exist") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRemove_should_return_error_when_name_only_path(t *testing.T) {
	t.Parallel()

	parent := t.TempDir()
	defaultPath := filepath.Join(parent, "cfg_missing")
	orig := clientconfig.DefaultPath
	clientconfig.DefaultPath = defaultPath
	t.Cleanup(func() {
		clientconfig.DefaultPath = orig
	})

	// Category C3: name-only input triggers non-absolute branch
	_, err := Remove("relative_name.py")
	if err == nil {
		t.Fatalf("Remove(name-only) error = nil; want non-nil")
	}
	if !strings.Contains(err.Error(), "exploit file does not exist") {
		t.Fatalf("unexpected error: %v", err)
	}
}
