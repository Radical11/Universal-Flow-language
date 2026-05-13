package analysis

import (
	stdjson "encoding/json"
	"io"

	"github.com/Radical11/Universal-Flow-language/internal/ir"
)

type Document struct {
	Version      string               `json:"version"`
	Title        string               `json:"title,omitempty"`
	Profile      string               `json:"profile,omitempty"`
	Adjacency    map[string][]string  `json:"adjacency"`
	Dependencies map[string][]string  `json:"dependencies"`
	Reverse      map[string][]string  `json:"reverseDependencies"`
}

func Encode(w io.Writer, doc *ir.Document) error {
	out := Document{
		Version:      doc.Version,
		Title:        doc.Title,
		Adjacency:    map[string][]string{},
		Dependencies: map[string][]string{},
		Reverse:      map[string][]string{},
	}
	if doc.Metadata != nil {
		if profile, ok := doc.Metadata["profile"].(string); ok {
			out.Profile = profile
		}
	}
	for _, node := range doc.Nodes {
		out.Adjacency[node.ID] = []string{}
		out.Dependencies[node.ID] = []string{}
		out.Reverse[node.ID] = []string{}
	}
	for _, edge := range doc.Edges {
		out.Adjacency[edge.From] = append(out.Adjacency[edge.From], edge.To)
		out.Dependencies[edge.From] = append(out.Dependencies[edge.From], edge.To)
		out.Reverse[edge.To] = append(out.Reverse[edge.To], edge.From)
	}
	for _, flow := range doc.Flows {
		for _, step := range flow.Steps {
			if step.Target == "" {
				continue
			}
			key := flow.ID + ":" + step.ID
			out.Adjacency[key] = append(out.Adjacency[key], step.Target)
			out.Dependencies[key] = append(out.Dependencies[key], step.Target)
			out.Reverse[step.Target] = append(out.Reverse[step.Target], key)
		}
	}
	encoder := stdjson.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(out)
}
