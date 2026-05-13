package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCompileWritesOutputFile(t *testing.T) {
	dir := t.TempDir()
	input := filepath.Join(dir, "input.ufl")
	output := filepath.Join(dir, "output.mmd")
	if err := os.WriteFile(input, []byte("entity a\nentity b\nrel a -> b\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := run([]string{"compile", input, "--target", "mermaid", "--out", output}); err != nil {
		t.Fatal(err)
	}
	content, err := os.ReadFile(output)
	if err != nil {
		t.Fatal(err)
	}
	if len(content) == 0 {
		t.Fatal("expected output content")
	}
}

func TestFmtWrite(t *testing.T) {
	dir := t.TempDir()
	input := filepath.Join(dir, "input.ufl")
	if err := os.WriteFile(input, []byte("entity a label \"A\" [z=last, a=first]\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := run([]string{"fmt", input, "--write"}); err != nil {
		t.Fatal(err)
	}
	content, err := os.ReadFile(input)
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != "# \n\nentity a label \"A\" [a=\"first\", z=\"last\"]\n" && string(content) != "entity a label \"A\" [a=\"first\", z=\"last\"]\n" {
		t.Fatalf("unexpected formatted content:\n%s", content)
	}
}

func TestValidateReportsDiagnosticLocation(t *testing.T) {
	dir := t.TempDir()
	input := filepath.Join(dir, "input.ufl")
	if err := os.WriteFile(input, []byte("entity a\nrel a -> missing\n"), 0644); err != nil {
		t.Fatal(err)
	}
	stderr := os.Stderr
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stderr = w
	defer func() { os.Stderr = stderr }()
	if err := run([]string{"validate", input}); err != nil {
		t.Fatal(err)
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r); err != nil {
		t.Fatal(err)
	}
	got := buf.String()
	if !strings.Contains(got, "warning:2:") {
		t.Fatalf("expected warning output with location, got %q", got)
	}
}
