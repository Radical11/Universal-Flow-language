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
	profile := profileName(doc)
	for _, node := range doc.Nodes {
		nodeID := mermaidID(node.ID)
		if _, err := fmt.Fprintf(w, "  %s[%q]\n", nodeID, nodeLabel(node)); err != nil {
			return err
		}
		if className := nodeClass(profile, node.Type); className != "" {
			if _, err := fmt.Fprintf(w, "  class %s %s\n", nodeID, className); err != nil {
				return err
			}
		}
		if _, err := writeNodeAnnotations(w, nodeID, node.Metadata); err != nil {
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
	if _, err := writeProfileStyleBlock(w, profile); err != nil {
		return err
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

func profileName(doc *ir.Document) string {
	if doc.Metadata == nil {
		return ""
	}
	if profile, ok := doc.Metadata["profile"].(string); ok {
		return profile
	}
	return ""
}

func nodeClass(profile, typ string) string {
	switch profile {
	case "workflow":
		switch typ {
		case "actor":
			return "actor"
		case "artifact":
			return "artifact"
		case "decision":
			return "decision"
		case "system":
			return "system"
		}
	case "story":
		switch typ {
		case "person":
			return "person"
		case "place":
			return "place"
		case "event":
			return "event"
		}
	case "system-map":
		switch typ {
		case "actor":
			return "actor"
		case "surface":
			return "surface"
		case "service":
			return "service"
		case "store":
			return "store"
		}
	}
	return ""
}

func writeProfileStyleBlock(w io.Writer, profile string) (int, error) {
	switch profile {
	case "workflow":
		if _, err := fmt.Fprintln(w, "  classDef actor fill:#d6f5e8,stroke:#1f7a4c"); err != nil {
			return 0, err
		}
		if _, err := fmt.Fprintln(w, "  classDef artifact fill:#fff4d6,stroke:#9a6b00"); err != nil {
			return 0, err
		}
		if _, err := fmt.Fprintln(w, "  classDef decision fill:#ffdede,stroke:#8c2f39"); err != nil {
			return 0, err
		}
		if _, err := fmt.Fprintln(w, "  classDef system fill:#dce9ff,stroke:#2d5ea8"); err != nil {
			return 0, err
		}
	case "story":
		if _, err := fmt.Fprintln(w, "  classDef person fill:#ffe3d1,stroke:#a85724"); err != nil {
			return 0, err
		}
		if _, err := fmt.Fprintln(w, "  classDef place fill:#dff3ff,stroke:#2f6f8f"); err != nil {
			return 0, err
		}
		if _, err := fmt.Fprintln(w, "  classDef event fill:#f4e0ff,stroke:#7d3ea1"); err != nil {
			return 0, err
		}
	case "system-map":
		if _, err := fmt.Fprintln(w, "  classDef actor fill:#e9f7da,stroke:#567d1d"); err != nil {
			return 0, err
		}
		if _, err := fmt.Fprintln(w, "  classDef surface fill:#dfefff,stroke:#35639f"); err != nil {
			return 0, err
		}
		if _, err := fmt.Fprintln(w, "  classDef service fill:#ffe7c8,stroke:#9b5e14"); err != nil {
			return 0, err
		}
		if _, err := fmt.Fprintln(w, "  classDef store fill:#f4def7,stroke:#804393"); err != nil {
			return 0, err
		}
	}
	return 0, nil
}

func writeNodeAnnotations(w io.Writer, nodeID string, meta map[string]any) (int, error) {
	if value, ok := meta["terminal"].(bool); ok && value {
		return fmt.Fprintf(w, "  style %s stroke-width:3px\n", nodeID)
	}
	return 0, nil
}
