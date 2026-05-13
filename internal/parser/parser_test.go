package parser

import "testing"

const sample = `# Story
entity hero as person label "Hero" [role=lead]
entity gate as place label "Gate"
rel hero -> gate as approaches [weight=0.8]

flow journey label "Opening journey"
{
  step arrive uses hero label "Arrive"
  step cross uses gate label "Cross threshold"
}

state idle label "Idle"
{
  on begin -> moving if "ready"
}
state moving label "Moving"
{
  on stop -> idle
}
`

func TestParseDocument(t *testing.T) {
	doc, err := Parse(sample)
	if err != nil {
		t.Fatal(err)
	}
	if doc.Title != "Story" {
		t.Fatalf("title got %q", doc.Title)
	}
	if len(doc.Entities) != 2 || len(doc.Relations) != 1 || len(doc.Flows) != 1 || len(doc.States) != 2 {
		t.Fatalf("unexpected document shape: %#v", doc)
	}
	if doc.Flows[0].Steps[1].Target != "gate" {
		t.Fatalf("unexpected step target: %#v", doc.Flows[0].Steps[1])
	}
	if doc.States[0].Transitions[0].Condition != "ready" {
		t.Fatalf("unexpected condition: %#v", doc.States[0].Transitions[0])
	}
}

func TestParseRejectsMalformedRelation(t *testing.T) {
	if _, err := Parse("rel a b\n"); err == nil {
		t.Fatal("expected parse error")
	}
}
