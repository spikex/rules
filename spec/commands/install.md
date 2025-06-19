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

## Error Handling

### Invalid Usage with Arguments

If the user provides arguments to the `install` command (e.g., `rules install starter/nextjs-rules`), the command should:

1. Display an error message indicating that `rules install` does not accept arguments
2. Suggest that they may have meant to use `rules add <rule>` instead
3. Exit with a non-zero status code
