# `rules render`

Renders existing rules to a specified format.

## Usage

```bash
rules render cursor
rules render continue
```

## Args

- Name of format to render rules to (e.g. "continue", "cursor")

## Behavior

- Copies all rules from the default location (`.rules/`) to the target format as described in [render-formats.md](../render-formats.md)
- Does not modify the original rule files