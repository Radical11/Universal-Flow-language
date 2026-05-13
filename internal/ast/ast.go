package ast

type Document struct {
	Title     string
	Sections  []Section
	Entities  []Entity
	Relations []Relation
	Flows     []Flow
	States    []State
}

type Section struct {
	Level int
	Title string
}

type Metadata map[string]string

type Entity struct {
	ID       string
	Type     string
	Label    string
	Metadata Metadata
}

type Relation struct {
	From     string
	To       string
	Type     string
	Metadata Metadata
}

type Flow struct {
	ID       string
	Label    string
	Steps    []Step
	Metadata Metadata
}

type Step struct {
	ID       string
	Target   string
	Label    string
	Metadata Metadata
}

type State struct {
	ID          string
	Label       string
	Transitions []Transition
	Metadata    Metadata
}

type Transition struct {
	To        string
	On        string
	Condition string
	Metadata  Metadata
}
