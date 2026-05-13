package compiler

import (
	"strings"
	"testing"
)

func TestFormatCanonicalizesDocument(t *testing.T) {
	doc, err := Parse(`# Example
entity b label "Bee" [z=last, a=first]
entity a as actor
rel a -> b
`)
	if err != nil {
		t.Fatal(err)
	}
	var out strings.Builder
	if err := Format(&out, doc); err != nil {
		t.Fatal(err)
	}
	got := out.String()
	if !strings.Contains(got, `entity b label "Bee" [a=first, z=last]`) {
		t.Fatalf("metadata was not sorted or formatted:\n%s", got)
	}
	if !strings.Contains(got, "rel a -> b") {
		t.Fatalf("relation missing:\n%s", got)
	}
}
