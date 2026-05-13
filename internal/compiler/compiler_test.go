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
	if doc.Diagnostics[0].Line == 0 || doc.Diagnostics[0].Column == 0 {
		t.Fatalf("expected source location on diagnostic, got %#v", doc.Diagnostics[0])
	}
}

func TestCompilePreservesTypedMetadataInIR(t *testing.T) {
	doc, err := Compile(`entity task [enabled=true, retries=3, tags=["ops", core]]`)
	if err != nil {
		t.Fatal(err)
	}
	meta := doc.Nodes[0].Metadata
	if meta["enabled"] != true {
		t.Fatalf("expected bool metadata, got %#v", meta["enabled"])
	}
	if meta["retries"] != float64(3) {
		t.Fatalf("expected numeric metadata, got %#v", meta["retries"])
	}
	tags, ok := meta["tags"].([]any)
	if !ok || len(tags) != 2 {
		t.Fatalf("expected array metadata, got %#v", meta["tags"])
	}
}
