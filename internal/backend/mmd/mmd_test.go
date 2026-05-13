package mmd

import (
	"strings"
	"testing"

	"github.com/Radical11/Universal-Flow-language/internal/compiler"
)

func TestEncodeAddsProfileStyles(t *testing.T) {
	doc, err := compiler.CompileWithProfile(`# System
entity user as actor label "User"
entity api as service label "API"
rel user -> api as interacts_with
`, compiler.ProfileSystemMap)
	if err != nil {
		t.Fatal(err)
	}
	var out strings.Builder
	if err := Encode(&out, doc); err != nil {
		t.Fatal(err)
	}
	got := out.String()
	for _, want := range []string{"class user actor", "class api service", "classDef actor", "classDef service"} {
		if !strings.Contains(got, want) {
			t.Fatalf("missing %q in:\n%s", want, got)
		}
	}
}
