# UFL IR Reference

The UFL compiler transforms source text into a universal semantic intermediate representation.

The first IR version is:

```text
ufl.ir.v0
```

## Document

```json
{
  "version": "ufl.ir.v0",
  "title": "Approval Workflow",
  "nodes": [],
  "edges": [],
  "flows": [],
  "states": [],
  "diagnostics": []
}
```

## Nodes

Nodes represent entities.

```json
{
  "id": "reviewer",
  "type": "actor",
  "label": "Reviewer",
  "metadata": {
    "team": "ops",
    "enabled": true,
    "retries": 3,
    "lanes": [
      "fast",
      "safe"
    ]
  }
}
```

## Edges

Edges represent relationships.

```json
{
  "id": "edge:1",
  "from": "request",
  "to": "reviewer",
  "type": "assigned_to"
}
```

## Flows

Flows represent ordered steps.

```json
{
  "id": "approval",
  "label": "Approval path",
  "steps": [
    {
      "id": "submit",
      "target": "request",
      "label": "Submit request"
    }
  ]
}
```

## States

States represent state-machine structure.

```json
{
  "id": "idle",
  "label": "Idle",
  "transitions": [
    {
      "on": "start",
      "to": "active"
    }
  ]
}
```

## Diagnostics

Diagnostics report semantic issues that do not necessarily prevent IR generation.

```json
{
  "severity": "warning",
  "message": "relation target \"reviewer\" is not declared",
  "line": 4,
  "column": 5
}
```

Examples:

- Relationship target not declared
- Flow step uses an undeclared target
- State transitions to an undeclared state

Parser errors still fail compilation.
