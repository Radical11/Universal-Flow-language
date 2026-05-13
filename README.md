# UFL: Universal Flow Language

UFL is a lightweight semantic scripting language for turning human-readable descriptions of stories, workflows, interactions, systems, journeys, and experiences into structured computational representations.

The language is intentionally small. It does not hardcode domains like games, cinema, UX, travel, education, or simulations. Instead, UFL provides a stable semantic core that compiles into an intermediate representation (IR) other tools can visualize, analyze, simulate, or transform.

## Quick Start

Create a file named `journey.ufl`:

```ufl
# Journey
entity hero as person label "Hero"
entity gate as place label "Gate"

rel hero -> gate as approaches

flow opening label "Opening movement"
{
  step arrive uses hero label "Arrive"
  step cross uses gate label "Cross the threshold"
}
```

Compile it to JSON IR:

```sh
go run ./cmd/ufl compile journey.ufl --target json
```

Validate an input:

```sh
go run ./cmd/ufl validate journey.ufl
```

Inspect the parsed AST:

```sh
go run ./cmd/ufl parse journey.ufl
```

## Core Ideas

- **Small syntax**: UFL uses readable declarations, Markdown-style headings, blocks, and metadata.
- **Semantic IR first**: the compiler output is a universal model of nodes, edges, flows, states, transitions, and metadata.
- **Domain agnostic**: the language core describes structure and meaning, not one industry.
- **Backend oriented**: the same IR can feed graph visualizers, simulation runtimes, diagrams, APIs, engines, or future plugins.
- **Extensible by ecosystem**: domains should live in libraries, validators, backends, and tooling rather than in the base language.

## Example Domains

UFL can describe:

- Story beats and character relationships
- Workflows and approvals
- UI interaction paths
- State machines
- Dependency graphs
- Journey maps
- System relationships
- Educational diagrams
- Simulation structures

These are all represented through the same low-level semantic core.

## Documentation

- [Language Manual](docs/language.md)
- [Grammar](docs/grammar.md)
- [IR Reference](docs/ir.md)
- [Compiler Architecture](docs/compiler.md)
- [Semantic Profiles](docs/profiles.md)
- [Plugin Model](docs/plugins.md)
- [CLI Reference](docs/cli.md)

## Current Status

This repository contains the first MVP:

- Go implementation
- Lexer and parser
- AST model
- Semantic IR builder
- JSON backend
- Mermaid and DOT graph backends
- CLI commands
- Examples
- Unit tests

Future milestones can add Tree-sitter editor support, graph visualizers, DOT/Mermaid backends, Unity and Unreal export targets, plugin loading, richer diagnostics, and domain-specific packages.
