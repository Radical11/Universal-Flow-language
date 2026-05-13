package ir

type Document struct {
	Version     string            `json:"version"`
	Title       string            `json:"title,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
	Nodes       []Node            `json:"nodes"`
	Edges       []Edge            `json:"edges"`
	Flows       []Flow            `json:"flows,omitempty"`
	States      []State           `json:"states,omitempty"`
	Diagnostics []Diagnostic      `json:"diagnostics,omitempty"`
}

type Node struct {
	ID       string            `json:"id"`
	Type     string            `json:"type"`
	Label    string            `json:"label,omitempty"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

type Edge struct {
	ID       string            `json:"id"`
	From     string            `json:"from"`
	To       string            `json:"to"`
	Type     string            `json:"type"`
	Weight   float64           `json:"weight,omitempty"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

type Flow struct {
	ID       string            `json:"id"`
	Label    string            `json:"label,omitempty"`
	Steps    []Step            `json:"steps"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

type Step struct {
	ID       string            `json:"id"`
	Target   string            `json:"target,omitempty"`
	Label    string            `json:"label,omitempty"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

type State struct {
	ID          string            `json:"id"`
	Label       string            `json:"label,omitempty"`
	Transitions []Transition      `json:"transitions,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

type Transition struct {
	To        string            `json:"to"`
	On        string            `json:"on"`
	Condition string            `json:"condition,omitempty"`
	Metadata  map[string]string `json:"metadata,omitempty"`
}

type Diagnostic struct {
	Severity string `json:"severity"`
	Message  string `json:"message"`
}
