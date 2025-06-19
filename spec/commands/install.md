# `rules install`

Synchronizes the `.rules` directory with the contents of `rules.json`.

## Usage

```bash
rules install
```

## Behavior

- Performs a clean installation by:
  - Removing all existing rule files and directories first
  - Re-downloading and installing all rules specified in `rules.json`
- Ensures the `.rules` directory exactly matches what's defined in `rules.json`
- Reports on installation progress and any errors encountered