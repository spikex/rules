# `rules init`

Creates initial rule directory structure.

## Usage

```bash
rules init
rules init --format cursor
```

## Flags

- `--format string`: Set rule format

## Behavior

- Creates directory structure
- Initializes empty rules.json
- Sets up format-specific configuration
- Default format creates `.rules/` directory
- Custom formats create `.{format}/rules/` directories