package timeline

import (
	"strings"
	"testing"

	"github.com/Radical11/Universal-Flow-language/internal/compiler"
)

func TestEncodeTimeline(t *testing.T) {
	doc, err := compiler.CompileWithProfile(`# Workflow
entity request as artifact
flow approval label "Approval path"
{
  step submit uses request label "Submit request"
}
`, compiler.ProfileWorkflow)
	if err != nil {
		t.Fatal(err)
	}
	var out strings.Builder
	if err := Encode(&out, doc); err != nil {
		t.Fatal(err)
	}
	got := out.String()
	for _, want := range []string{`"profile": "workflow"`, `"timelines"`, `"id": "approval"`, `"index": 0`, `"target": "request"`} {
		if !strings.Contains(got, want) {
			t.Fatalf("missing %q in:\n%s", want, got)
		}
	}
}
