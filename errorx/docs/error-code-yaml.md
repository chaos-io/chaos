# Error Code YAML

This document describes the YAML input accepted by `gen_error_code`.

## Fields

```yaml
appCode: 6
bizCode: 12
errorCode:
  - name: TaskNotFound
    code: 1001
    message: task {task_id} not found
    description: task does not exist
    affectsStability: false
```

- `appCode`: required, `1..9`.
- `bizCode`: required, `1..9999`.
- `errorCode`: required, non-empty list.
- `name`: required, exported Go identifier. It becomes the generated variable name.
- `code`: required, `1..9999`.
- `message`: optional business message. Empty values use `errorx.DefaultMessage`.
- `description`: optional generated Go comment.
- `affectsStability`: optional, defaults to `true`.

## Code Layout

The generated full code is a nine-digit integer:

```text
appCode * 100000000 + bizCode * 10000 + code
```

Example:

```text
6 * 100000000 + 12 * 10000 + 1001 = 600121001
```

## Naming Rules

Use stable business names:

```yaml
- name: TaskNotFound
- name: UserAlreadyExists
- name: PaymentProviderUnavailable
```

Avoid transport-specific names:

```yaml
- name: HTTP404
- name: GRPCInvalidArgument
```

Transport mapping should live in each service's adapter layer, not in `errorx`.

## Conflict Rules

Generation fails when one input batch contains:

- Duplicate generated file names, such as `foo-bar.yaml` and `foo_bar.yaml`.
- Duplicate error names across files.
- Duplicate full error codes across files.
- Invalid code ranges.
- Non-exported or invalid Go identifiers.

Runtime registration also rejects conflicting definitions for the same code. Registering the exact same definition more than once is allowed.
