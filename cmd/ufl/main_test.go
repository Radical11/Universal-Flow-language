package main

import (
	"os"
	"path/filepath"
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
	if string(content) != "# \n\nentity a label \"A\" [a=first, z=last]\n" && string(content) != "entity a label \"A\" [a=first, z=last]\n" {
		t.Fatalf("unexpected formatted content:\n%s", content)
	}
}
