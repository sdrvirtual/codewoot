package services

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRemoveInstanceFilesRemovesOnlyUUIDDirectory(t *testing.T) {
	base := t.TempDir()
	instance := "d4b9e6aa-f0b7-4fdf-9b91-24d080664bae"
	target := filepath.Join(base, instance)
	if err := os.MkdirAll(target, 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	if err := os.WriteFile(filepath.Join(target, "creds.json"), []byte("{}"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	if err := removeInstanceFiles(base, instance); err != nil {
		t.Fatalf("removeInstanceFiles() error = %v", err)
	}
	if _, err := os.Stat(target); !os.IsNotExist(err) {
		t.Fatalf("expected target directory to be removed, stat error = %v", err)
	}
	if _, err := os.Stat(base); err != nil {
		t.Fatalf("expected base directory to remain, stat error = %v", err)
	}
}

func TestRemoveInstanceFilesRejectsInvalidInstance(t *testing.T) {
	base := t.TempDir()
	keep := filepath.Join(base, "keep")
	if err := os.MkdirAll(keep, 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}

	if err := removeInstanceFiles(base, "../keep"); err == nil {
		t.Fatalf("expected invalid uuid error")
	}
	if _, err := os.Stat(keep); err != nil {
		t.Fatalf("expected unrelated directory to remain, stat error = %v", err)
	}
}

func TestRemoveInstanceFilesRequiresConfiguredDirectoryToExist(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "missing")
	err := removeInstanceFiles(missing, "d4b9e6aa-f0b7-4fdf-9b91-24d080664bae")
	if err == nil {
		t.Fatalf("expected missing directory error")
	}
}

func TestRemoveInstanceFilesNoopWhenDirectoryNotConfigured(t *testing.T) {
	err := removeInstanceFiles("", "not-a-uuid")
	if err != nil {
		t.Fatalf("expected no-op without configured directory, got %v", err)
	}
}
