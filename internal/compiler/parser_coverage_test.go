package compiler

import (
	"testing"

	"github.com/Radical11/Universal-Flow-language/internal/ast"
	"github.com/Radical11/Universal-Flow-language/internal/parser"
)

const parserSample = `# Story
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

func TestParserBuildsDocumentShape(t *testing.T) {
	doc, err := parser.Parse(parserSample)
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
	if doc.Relations[0].Metadata["weight"].Kind != ast.ValueNumber {
		t.Fatalf("expected numeric metadata, got %#v", doc.Relations[0].Metadata["weight"])
	}
}

func TestParserRejectsMalformedRelation(t *testing.T) {
	if _, err := parser.Parse("rel a b\n"); err == nil {
		t.Fatal("expected parse error")
	}
}

func TestParserSupportsTypedMetadataValues(t *testing.T) {
	doc, err := parser.Parse(`entity user [enabled=true, retries=3, tags=["alpha", beta, 2]]`)
	if err != nil {
		t.Fatal(err)
	}
	meta := doc.Entities[0].Metadata
	if !meta["enabled"].Bool {
		t.Fatalf("expected true bool, got %#v", meta["enabled"])
	}
	if meta["retries"].Number != 3 {
		t.Fatalf("expected number 3, got %#v", meta["retries"])
	}
	if meta["tags"].Kind != ast.ValueArray || len(meta["tags"].Array) != 3 {
		t.Fatalf("expected array metadata, got %#v", meta["tags"])
	}
}
