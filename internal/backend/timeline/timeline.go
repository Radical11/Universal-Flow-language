package timeline

import (
	stdjson "encoding/json"
	"io"

	"github.com/Radical11/Universal-Flow-language/internal/ir"
)

type Document struct {
	Version   string     `json:"version"`
	Title     string     `json:"title,omitempty"`
	Profile   string     `json:"profile,omitempty"`
	Timelines []Timeline `json:"timelines"`
}

type Timeline struct {
	ID    string      `json:"id"`
	Label string      `json:"label,omitempty"`
	Steps []TimelineStep `json:"steps"`
}

type TimelineStep struct {
	Index    int            `json:"index"`
	ID       string         `json:"id"`
	Label    string         `json:"label,omitempty"`
	Target   string         `json:"target,omitempty"`
	Metadata map[string]any `json:"metadata,omitempty"`
}

func Encode(w io.Writer, doc *ir.Document) error {
	out := Document{
		Version:   doc.Version,
		Title:     doc.Title,
		Timelines: make([]Timeline, 0, len(doc.Flows)),
	}
	if doc.Metadata != nil {
		if profile, ok := doc.Metadata["profile"].(string); ok {
			out.Profile = profile
		}
	}
	for _, flow := range doc.Flows {
		timeline := Timeline{
			ID:    flow.ID,
			Label: flow.Label,
			Steps: make([]TimelineStep, 0, len(flow.Steps)),
		}
		for i, step := range flow.Steps {
			timeline.Steps = append(timeline.Steps, TimelineStep{
				Index:    i,
				ID:       step.ID,
				Label:    step.Label,
				Target:   step.Target,
				Metadata: step.Metadata,
			})
		}
		out.Timelines = append(out.Timelines, timeline)
	}
	encoder := stdjson.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(out)
}
