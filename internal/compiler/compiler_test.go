package compiler

import "testing"

func TestCompileBuildsSemanticIR(t *testing.T) {
	doc, err := Compile(`# Workflow
entity request as input label "Request"
entity response as output label "Response"
rel request -> response as produces
`)
	if err != nil {
		t.Fatal(err)
	}
	if doc.Version != IRVersion {
		t.Fatalf("version got %q", doc.Version)
	}
	if len(doc.Nodes) != 2 || len(doc.Edges) != 1 {
		t.Fatalf("unexpected IR shape: %#v", doc)
	}
	if doc.Edges[0].Type != "produces" {
		t.Fatalf("edge type got %q", doc.Edges[0].Type)
	}
}

func TestCompileWarnsForUndeclaredRelationTarget(t *testing.T) {
	doc, err := Compile("entity a\nrel a -> missing\n")
	if err != nil {
		t.Fatal(err)
	}
	if len(doc.Diagnostics) != 1 {
		t.Fatalf("expected one diagnostic, got %#v", doc.Diagnostics)
	}
}
