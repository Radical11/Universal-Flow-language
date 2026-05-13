package dot

import (
	"strings"
	"testing"

	"github.com/Radical11/Universal-Flow-language/internal/compiler"
)

func TestEncodeDOT(t *testing.T) {
	doc, err := compiler.Compile(`# Network
entity a as actor label "A"
entity b as system label "B"
rel a -> b as calls
`)
	if err != nil {
		t.Fatal(err)
	}
	var out strings.Builder
	if err := Encode(&out, doc); err != nil {
		t.Fatal(err)
	}
	got := out.String()
	for _, want := range []string{"digraph UFL", "rankdir=LR", "a [label=\"A\"", "a -> b [label=\"calls\"]"} {
		if !strings.Contains(got, want) {
			t.Fatalf("missing %q in:\n%s", want, got)
		}
	}
}
