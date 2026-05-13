package runtime

import (
	"strings"
	"testing"

	"github.com/Radical11/Universal-Flow-language/internal/compiler"
)

func TestEncodeRuntime(t *testing.T) {
	doc, err := compiler.CompileWithProfile(`# Machine
state idle [initial=true]
{
  on start -> active
}
state active
{
  on finish -> done
}
state done [terminal=true]
{
}
`, compiler.ProfileStateMachine)
	if err != nil {
		t.Fatal(err)
	}
	var out strings.Builder
	if err := Encode(&out, doc); err != nil {
		t.Fatal(err)
	}
	got := out.String()
	for _, want := range []string{`"profile": "state-machine"`, `"initialStates"`, `"idle"`, `"terminalStates"`, `"done"`, `"event": "start"`} {
		if !strings.Contains(got, want) {
			t.Fatalf("missing %q in:\n%s", want, got)
		}
	}
}

