package runtime

import (
	stdjson "encoding/json"
	"io"

	"github.com/Radical11/Universal-Flow-language/internal/ir"
)

type Document struct {
	Version       string            `json:"version"`
	Title         string            `json:"title,omitempty"`
	Profile       string            `json:"profile,omitempty"`
	InitialStates []string          `json:"initialStates,omitempty"`
	TerminalStates []string         `json:"terminalStates,omitempty"`
	States        []StateRuntime    `json:"states,omitempty"`
	Flows         []FlowRuntime     `json:"flows,omitempty"`
}

type StateRuntime struct {
	ID          string              `json:"id"`
	Label       string              `json:"label,omitempty"`
	Terminal    bool                `json:"terminal,omitempty"`
	Transitions []TransitionRuntime `json:"transitions,omitempty"`
}

type TransitionRuntime struct {
	Event     string         `json:"event"`
	Target    string         `json:"target"`
	Condition string         `json:"condition,omitempty"`
	Metadata  map[string]any `json:"metadata,omitempty"`
}

type FlowRuntime struct {
	ID       string            `json:"id"`
	Label    string            `json:"label,omitempty"`
	Sequence []FlowStepRuntime `json:"sequence"`
}

type FlowStepRuntime struct {
	Index    int            `json:"index"`
	ID       string         `json:"id"`
	Target   string         `json:"target,omitempty"`
	Label    string         `json:"label,omitempty"`
	Metadata map[string]any `json:"metadata,omitempty"`
}

func Encode(w io.Writer, doc *ir.Document) error {
	out := Document{
		Version: doc.Version,
		Title:   doc.Title,
		States:  make([]StateRuntime, 0, len(doc.States)),
		Flows:   make([]FlowRuntime, 0, len(doc.Flows)),
	}
	if doc.Metadata != nil {
		if profile, ok := doc.Metadata["profile"].(string); ok {
			out.Profile = profile
		}
	}
	for _, state := range doc.States {
		initial, _ := state.Metadata["initial"].(bool)
		terminal, _ := state.Metadata["terminal"].(bool)
		if initial {
			out.InitialStates = append(out.InitialStates, state.ID)
		}
		if terminal {
			out.TerminalStates = append(out.TerminalStates, state.ID)
		}
		runtimeState := StateRuntime{
			ID:       state.ID,
			Label:    state.Label,
			Terminal: terminal,
		}
		for _, transition := range state.Transitions {
			runtimeState.Transitions = append(runtimeState.Transitions, TransitionRuntime{
				Event:     transition.On,
				Target:    transition.To,
				Condition: transition.Condition,
				Metadata:  transition.Metadata,
			})
		}
		out.States = append(out.States, runtimeState)
	}
	for _, flow := range doc.Flows {
		runtimeFlow := FlowRuntime{
			ID:       flow.ID,
			Label:    flow.Label,
			Sequence: make([]FlowStepRuntime, 0, len(flow.Steps)),
		}
		for i, step := range flow.Steps {
			runtimeFlow.Sequence = append(runtimeFlow.Sequence, FlowStepRuntime{
				Index:    i,
				ID:       step.ID,
				Target:   step.Target,
				Label:    step.Label,
				Metadata: step.Metadata,
			})
		}
		out.Flows = append(out.Flows, runtimeFlow)
	}
	encoder := stdjson.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(out)
}

