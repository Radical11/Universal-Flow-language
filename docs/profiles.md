# Semantic Profiles

Semantic profiles are conventions layered on top of core UFL. They are not new syntax.

A profile can define recommended entity types, relationship types, metadata keys, validation rules, and visualization preferences for a domain.

Use profiles from the CLI with:

```sh
go run ./cmd/ufl validate input.ufl --profile workflow
go run ./cmd/ufl inspect input.ufl --profile story
go run ./cmd/ufl compile input.ufl --target mermaid --profile system-map
go run ./cmd/ufl compile input.ufl --target runtime --profile state-machine
```

## Why Profiles Exist

The language core should stay small. Profiles let communities build richer meanings without asking the parser to understand every domain.

## Workflow Profile

Suggested entity types:

- `actor`
- `artifact`
- `system`
- `decision`

Suggested relationship types:

- `assigned_to`
- `produces`
- `depends_on`
- `approves`
- `rejects`

Suggested metadata keys:

- `owner`
- `sla`
- `priority`
- `status`

Current validator behavior:

- Warn when no `flow` blocks exist
- Warn on unrecognized entity types
- Warn on unrecognized relationship types

## Story Profile

Suggested entity types:

- `person`
- `place`
- `event`
- `object`
- `theme`

Suggested relationship types:

- `begins_at`
- `travels_to`
- `conflicts_with`
- `reveals`
- `changes`

Suggested metadata keys:

- `role`
- `mood`
- `stakes`
- `time`

Current validator behavior:

- Warn when no relations or flows exist
- Warn on unrecognized entity types
- Warn on unrecognized relationship types

## State Machine Profile

Suggested entity types are usually unnecessary because `state` blocks carry the model.

Suggested metadata keys:

- `initial`
- `terminal`
- `priority`
- `timeout`

Current validator behavior:

- Warn when no `state` blocks exist
- Warn when no state has `[initial=true]`
- Warn when more than one state has `[initial=true]`
- Warn when a state has no outgoing transitions

## System Map Profile

Suggested entity types:

- `actor`
- `surface`
- `service`
- `store`
- `queue`
- `external`

Suggested relationship types:

- `interacts_with`
- `calls`
- `reads_from`
- `writes_to`
- `publishes`
- `subscribes`
- `returns_to`

Current validator behavior:

- Warn when entities or relations are missing
- Warn on unrecognized entity types
- Warn on unrecognized relationship types

## Profile Compatibility

Profiles should prefer:

- Metadata over new syntax
- Validators over parser changes
- Backend conventions over IR changes

When a profile needs something the IR cannot represent, that is a signal to discuss the IR carefully rather than expanding syntax immediately.
