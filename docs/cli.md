# UFL CLI Reference

The CLI lives at:

```sh
go run ./cmd/ufl
```

After installing or building a binary, the same commands can be run as `ufl`.

## Parse

Parse source into the AST.

```sh
go run ./cmd/ufl parse examples/story.ufl
```

Use this when debugging syntax-level structure.

## Validate

Validate source and report diagnostics.

```sh
go run ./cmd/ufl validate examples/workflow.ufl
```

Successful output:

```text
valid
```

Warnings are printed to stderr.

## Compile

Compile source to a backend target.

```sh
go run ./cmd/ufl compile examples/network.ufl --target json
```

Current targets:

- `json`
- `mermaid`
- `dot`

Write output to a file:

```sh
go run ./cmd/ufl compile examples/workflow.ufl --target mermaid --out workflow.mmd
```

## Inspect

Print a quick semantic summary.

```sh
go run ./cmd/ufl inspect examples/story.ufl
```

Example output:

```text
title: Story Journey
nodes: 3
edges: 2
flows: 1
states: 0
diagnostics: 0
```

## Format

Print canonical UFL formatting:

```sh
go run ./cmd/ufl fmt examples/workflow.ufl
```

Rewrite a file in place:

```sh
go run ./cmd/ufl fmt examples/workflow.ufl --write
```

The formatter emits a canonical representation. Metadata strings are always written with quotes, even if the source used a bare identifier.

## Common Development Commands

Run tests:

```sh
go test ./...
```

Format code:

```sh
gofmt -w cmd internal
```
