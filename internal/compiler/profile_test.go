package compiler

import "testing"

func TestWorkflowProfileWarnsWithoutFlow(t *testing.T) {
	doc, err := CompileWithProfile(`entity request as artifact
entity reviewer as actor
rel request -> reviewer as assigned_to
`, ProfileWorkflow)
	if err != nil {
		t.Fatal(err)
	}
	if len(doc.Diagnostics) == 0 {
		t.Fatal("expected workflow diagnostics")
	}
}

func TestStateMachineProfileRequiresInitialState(t *testing.T) {
	doc, err := CompileWithProfile(`state idle
{
  on start -> active
}
state active
{
  on stop -> idle
}
`, ProfileStateMachine)
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, diagnostic := range doc.Diagnostics {
		if diagnostic.Message == `state-machine profile expects one state with metadata [initial=true]` {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected missing initial state diagnostic, got %#v", doc.Diagnostics)
	}
}

func TestSystemMapProfileWarnsOnUnknownEntityType(t *testing.T) {
	doc, err := CompileWithProfile(`entity batch as artifact
entity api as service
rel api -> batch as calls
`, ProfileSystemMap)
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, diagnostic := range doc.Diagnostics {
		if diagnostic.Message == `system-map profile does not recognize entity type "artifact"` {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected system-map entity type diagnostic, got %#v", doc.Diagnostics)
	}
}

func TestNormalizeProfileRejectsUnknownValue(t *testing.T) {
	if _, err := NormalizeProfile("unknown"); err == nil {
		t.Fatal("expected error")
	}
}
