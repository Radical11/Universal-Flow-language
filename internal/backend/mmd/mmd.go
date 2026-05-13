package mmd

import (
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/Radical11/Universal-Flow-language/internal/ir"
)

var unsafeID = regexp.MustCompile(`[^A-Za-z0-9_]`)

func Encode(w io.Writer, doc *ir.Document) error {
	if len(doc.Nodes) > 0 || len(doc.Edges) > 0 || len(doc.Flows) > 0 {
		if err := encodeGraph(w, doc); err != nil {
			return err
		}
	}
	if len(doc.States) > 0 {
		if len(doc.Nodes) > 0 || len(doc.Edges) > 0 || len(doc.Flows) > 0 {
			if _, err := fmt.Fprintln(w); err != nil {
				return err
			}
		}
		return encodeStateDiagram(w, doc)
	}
	return nil
}

func encodeGraph(w io.Writer, doc *ir.Document) error {
	if _, err := fmt.Fprintln(w, "flowchart TD"); err != nil {
		return err
	}
	for _, node := range doc.Nodes {
		if _, err := fmt.Fprintf(w, "  %s[%q]\n", mermaidID(node.ID), nodeLabel(node)); err != nil {
			return err
		}
	}
	for _, edge := range doc.Edges {
		if _, err := fmt.Fprintf(w, "  %s -->|%s| %s\n", mermaidID(edge.From), escapeLabel(edge.Type), mermaidID(edge.To)); err != nil {
			return err
		}
	}
	for _, flow := range doc.Flows {
		previous := ""
		for _, step := range flow.Steps {
			stepID := mermaidID(flow.ID + "_" + step.ID)
			label := step.Label
			if label == "" {
				label = step.ID
			}
			if _, err := fmt.Fprintf(w, "  %s[%q]\n", stepID, label); err != nil {
				return err
			}
			if step.Target != "" {
				if _, err := fmt.Fprintf(w, "  %s -.uses.-> %s\n", stepID, mermaidID(step.Target)); err != nil {
					return err
				}
			}
			if previous != "" {
				if _, err := fmt.Fprintf(w, "  %s --> %s\n", previous, stepID); err != nil {
					return err
				}
			}
			previous = stepID
		}
	}
	return nil
}

func encodeStateDiagram(w io.Writer, doc *ir.Document) error {
	if _, err := fmt.Fprintln(w, "stateDiagram-v2"); err != nil {
		return err
	}
	for _, state := range doc.States {
		if state.Label != "" && state.Label != state.ID {
			if _, err := fmt.Fprintf(w, "  %s: %s\n", mermaidID(state.ID), escapeLabel(state.Label)); err != nil {
				return err
			}
		}
		for _, transition := range state.Transitions {
			label := transition.On
			if transition.Condition != "" {
				label += " if " + transition.Condition
			}
			if _, err := fmt.Fprintf(w, "  %s --> %s: %s\n", mermaidID(state.ID), mermaidID(transition.To), escapeLabel(label)); err != nil {
				return err
			}
		}
	}
	return nil
}

func nodeLabel(node ir.Node) string {
	if node.Label != "" {
		return node.Label
	}
	return node.ID
}

func mermaidID(id string) string {
	id = unsafeID.ReplaceAllString(id, "_")
	if id == "" {
		return "_"
	}
	return id
}

func escapeLabel(value string) string {
	return strings.ReplaceAll(value, `"`, `'`)
}
