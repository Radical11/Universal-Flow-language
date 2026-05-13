package dot

import (
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/Radical11/Universal-Flow-language/internal/ir"
)

var unsafeID = regexp.MustCompile(`[^A-Za-z0-9_]`)

func Encode(w io.Writer, doc *ir.Document) error {
	if _, err := fmt.Fprintln(w, "digraph UFL {"); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w, "  rankdir=LR;"); err != nil {
		return err
	}
	for _, node := range doc.Nodes {
		if _, err := fmt.Fprintf(w, "  %s [label=%q, shape=box];\n", dotID(node.ID), nodeLabel(node)); err != nil {
			return err
		}
	}
	for _, edge := range doc.Edges {
		if _, err := fmt.Fprintf(w, "  %s -> %s [label=%q];\n", dotID(edge.From), dotID(edge.To), edge.Type); err != nil {
			return err
		}
	}
	for _, flow := range doc.Flows {
		previous := ""
		for _, step := range flow.Steps {
			stepID := dotID(flow.ID + "_" + step.ID)
			label := step.Label
			if label == "" {
				label = step.ID
			}
			if _, err := fmt.Fprintf(w, "  %s [label=%q, shape=ellipse];\n", stepID, label); err != nil {
				return err
			}
			if step.Target != "" {
				if _, err := fmt.Fprintf(w, "  %s -> %s [label=%q, style=dashed];\n", stepID, dotID(step.Target), "uses"); err != nil {
					return err
				}
			}
			if previous != "" {
				if _, err := fmt.Fprintf(w, "  %s -> %s [label=%q];\n", previous, stepID, "next"); err != nil {
					return err
				}
			}
			previous = stepID
		}
	}
	for _, state := range doc.States {
		if _, err := fmt.Fprintf(w, "  %s [label=%q, shape=oval];\n", dotID(state.ID), stateLabel(state)); err != nil {
			return err
		}
		for _, transition := range state.Transitions {
			label := transition.On
			if transition.Condition != "" {
				label += " if " + transition.Condition
			}
			if _, err := fmt.Fprintf(w, "  %s -> %s [label=%q];\n", dotID(state.ID), dotID(transition.To), label); err != nil {
				return err
			}
		}
	}
	_, err := fmt.Fprintln(w, "}")
	return err
}

func nodeLabel(node ir.Node) string {
	if node.Label != "" {
		return node.Label
	}
	return node.ID
}

func stateLabel(state ir.State) string {
	if state.Label != "" {
		return state.Label
	}
	return state.ID
}

func dotID(id string) string {
	id = unsafeID.ReplaceAllString(id, "_")
	if id == "" {
		return "_"
	}
	if id[0] >= '0' && id[0] <= '9' {
		return "_" + id
	}
	return strings.TrimSpace(id)
}
