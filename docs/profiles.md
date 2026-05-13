# Semantic Profiles

Semantic profiles are conventions layered on top of core UFL. They are not new syntax.

A profile can define recommended entity types, relationship types, metadata keys, validation rules, and visualization preferences for a domain.

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

## State Machine Profile

Suggested entity types are usually unnecessary because `state` blocks carry the model.

Suggested metadata keys:

- `initial`
- `terminal`
- `priority`
- `timeout`

## System Map Profile

Suggested entity types:

- `actor`
- `surface`
- `service`
- `store`
- `queue`
- `external`

Suggested relationship types:

- `calls`
- `reads_from`
- `writes_to`
- `publishes`
- `subscribes`
- `returns_to`

## Profile Compatibility

Profiles should prefer:

- Metadata over new syntax
- Validators over parser changes
- Backend conventions over IR changes

When a profile needs something the IR cannot represent, that is a signal to discuss the IR carefully rather than expanding syntax immediately.

