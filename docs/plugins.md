# UFL Plugin Model

UFL is designed so the core language stays small while the ecosystem grows around it.

## Core Rule

The base language should define structure, not domains.

Core syntax should not know about films, games, UI flows, analytics, simulations, education, APIs, or travel. These belong in extensions.

## Extension Types

Future UFL extensions can provide:

- Domain validators
- Metadata conventions
- IR transforms
- Custom backends
- Runtime adapters
- Visualization systems
- Editor tooling
- Simulation engines

## Examples

A cinematic package could define metadata conventions like:

```ufl
entity opening_shot as beat [camera=wide, mood=quiet]
```

A workflow package could validate:

```ufl
entity approval as decision [sla=48h]
```

The language parser does not need to understand `camera`, `mood`, or `sla`. Those meanings are interpreted by external tooling.

## Backend Interface Direction

The current JSON backend is built directly against the IR. Future backends should follow the same idea:

```text
IR document -> backend -> target output
```

Potential targets:

- Graph JSON
- Mermaid
- DOT/Graphviz
- Timeline JSON
- Unity C# data
- Unreal C++ data
- Web visualization data
- Simulation runtime data

## Stability Goal

The core IR should evolve carefully. Extensions should prefer metadata and transforms before requesting new syntax.
