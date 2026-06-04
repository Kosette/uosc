package lib

import (
	"archive/zip"
	"os"
	"path/filepath"
	"testing"
)

func TestZipFilesWithHeadersNilHeaderMod(t *testing.T) {
	root := t.TempDir()
	srcDir := filepath.Join(root, "src")
	outPath := filepath.Join(root, "release", "uosc.zip")

	if err := os.MkdirAll(filepath.Join(srcDir, "nested"), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	fileA := filepath.Join(srcDir, "root.txt")
	if err := os.WriteFile(fileA, []byte("root"), 0o644); err != nil {
		t.Fatalf("write root file: %v", err)
	}

	fileB := filepath.Join(srcDir, "nested", "child.txt")
	if err := os.WriteFile(fileB, []byte("child"), 0o644); err != nil {
		t.Fatalf("write nested file: %v", err)
	}

	stats, err := ZipFilesWithHeaders(map[string]string{
		srcDir: "scripts/",
	}, outPath, nil)
	if err != nil {
		t.Fatalf("ZipFilesWithHeaders returned error: %v", err)
	}

	if stats.FilesNum != 2 {
		t.Fatalf("expected 2 archived files, got %d", stats.FilesNum)
	}

	reader, err := zip.OpenReader(outPath)
	if err != nil {
		t.Fatalf("open zip: %v", err)
	}
	defer reader.Close()

	got := make(map[string]bool, len(reader.File))
	for _, file := range reader.File {
		got[file.Name] = true
	}

	for _, want := range []string{
		"scripts/root.txt",
		"scripts/nested/child.txt",
	} {
		if !got[want] {
			t.Fatalf("expected archive entry %q, got %#v", want, got)
		}
	}
}
