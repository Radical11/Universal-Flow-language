package compiler

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	jsonbackend "github.com/Radical11/Universal-Flow-language/internal/backend/json"
)

func TestWorkflowJSONSnapshot(t *testing.T) {
	source, err := os.ReadFile(filepath.Join("..", "..", "examples", "workflow.ufl"))
	if err != nil {
		t.Fatal(err)
	}
	doc, err := Compile(string(source))
	if err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	if err := jsonbackend.Encode(&out, doc); err != nil {
		t.Fatal(err)
	}
	want, err := os.ReadFile(filepath.Join("testdata", "workflow.golden.json"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.TrimSpace(out.String()) != strings.TrimSpace(string(want)) {
		t.Fatalf("snapshot mismatch\nwant:\n%s\n\ngot:\n%s", want, out.String())
	}
}
