package analysis

import (
	"strings"
	"testing"

	"github.com/Radical11/Universal-Flow-language/internal/compiler"
)

func TestEncodeAnalysis(t *testing.T) {
	doc, err := compiler.CompileWithProfile(`# Network
entity user as actor
entity api as service
rel user -> api as calls
flow access
{
  step request uses api
}
`, compiler.ProfileSystemMap)
	if err != nil {
		t.Fatal(err)
	}
	var out strings.Builder
	if err := Encode(&out, doc); err != nil {
		t.Fatal(err)
	}
	got := out.String()
	for _, want := range []string{`"adjacency"`, `"dependencies"`, `"reverseDependencies"`, `"user"`, `"api"`, `"access:request"`} {
		if !strings.Contains(got, want) {
			t.Fatalf("missing %q in:\n%s", want, got)
		}
	}
}

