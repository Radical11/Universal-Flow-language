# Compiler Architecture

The UFL compiler is a semantic transformation engine. It does not compile to machine code. It turns readable source text into a structured IR that other systems can consume.

## Pipeline

```text
source .ufl
  -> lexer
  -> parser
  -> AST
  -> semantic analyzer
  -> IR builder
  -> backend
```

## Lexer

The lexer converts source text into tokens such as identifiers, strings, headings, arrows, braces, and metadata punctuation.

It is intentionally small and line-oriented so the language remains easy to reason about.

## Parser

The parser turns tokens into an AST. The current grammar supports:

- Headings
- Entities
- Relationships
- Flows
- Steps
- States
- Transitions
- Metadata

## AST

The AST preserves the source-level structure. It is useful for diagnostics, source-aware tooling, future formatters, and editor integrations.

## Semantic IR

The IR normalizes source concepts into a backend-friendly model:

- Entities become nodes.
- Relationships become edges.
- Flows become ordered step lists.
- States become transition systems.
- Metadata remains attached to each semantic object.

## Backends

Backends transform the IR into target formats. Current backends emit JSON, Mermaid, DOT, timeline JSON, runtime JSON, and analysis JSON.

Future backends can emit graph formats, diagrams, runtime data, engine integrations, or simulation inputs without changing the core parser.
