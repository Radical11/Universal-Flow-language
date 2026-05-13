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

func TestCompileTimelineWithProfile(t *testing.T) {
	dir := t.TempDir()
	input := filepath.Join(dir, "input.ufl")
	output := filepath.Join(dir, "output.json")
	source := "# Workflow\nentity request as artifact\nflow approval\n{\n  step submit uses request label \"Submit request\"\n}\n"
	if err := os.WriteFile(input, []byte(source), 0644); err != nil {
		t.Fatal(err)
	}
	if err := run([]string{"compile", input, "--target", "timeline", "--profile", "workflow", "--out", output}); err != nil {
		t.Fatal(err)
	}
	content, err := os.ReadFile(output)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(content), `"profile": "workflow"`) {
		t.Fatalf("expected timeline profile in output:\n%s", content)
	}
}

func TestCompileRuntimeTarget(t *testing.T) {
	dir := t.TempDir()
	input := filepath.Join(dir, "machine.ufl")
	output := filepath.Join(dir, "runtime.json")
	source := "# Machine\nstate idle [initial=true]\n{\n  on start -> done\n}\nstate done [terminal=true]\n{\n}\n"
	if err := os.WriteFile(input, []byte(source), 0644); err != nil {
		t.Fatal(err)
	}
	if err := run([]string{"compile", input, "--target", "runtime", "--profile", "state-machine", "--out", output}); err != nil {
		t.Fatal(err)
	}
	content, err := os.ReadFile(output)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(content), `"initialStates"`) || !strings.Contains(string(content), `"terminalStates"`) {
		t.Fatalf("expected runtime output:\n%s", content)
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

func TestValidateWithProfileReportsProfileWarning(t *testing.T) {
	dir := t.TempDir()
	input := filepath.Join(dir, "input.ufl")
	if err := os.WriteFile(input, []byte("entity request as artifact\nentity reviewer as actor\nrel request -> reviewer as assigned_to\n"), 0644); err != nil {
		t.Fatal(err)
	}
	stderr := os.Stderr
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stderr = w
	defer func() { os.Stderr = stderr }()
	if err := run([]string{"validate", input, "--profile", "workflow"}); err != nil {
		t.Fatal(err)
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), `workflow profile expects at least one "flow" block`) {
		t.Fatalf("expected profile warning, got %q", buf.String())
	}
}

func TestInspectShowsProfile(t *testing.T) {
	dir := t.TempDir()
	input := filepath.Join(dir, "input.ufl")
	if err := os.WriteFile(input, []byte("state idle [initial=true]\n{\n  on start -> active\n}\nstate active\n{\n  on stop -> idle\n}\n"), 0644); err != nil {
		t.Fatal(err)
	}
	stdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout = w
	defer func() { os.Stdout = stdout }()
	if err := run([]string{"inspect", input, "--profile", "state-machine"}); err != nil {
		t.Fatal(err)
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "profile: state-machine") {
		t.Fatalf("expected profile in inspect output, got %q", buf.String())
	}
}
