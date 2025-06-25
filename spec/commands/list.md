# Rules List Command

## Overview

The `rules list` command displays all rules currently installed in the project, similar to `npm list`. It reads from the `rules.json` file and shows the installed rules with their versions in a tree-like format.

## Usage

```bash
rules list [flags]
```

## Behavior

1. **Check for rules.json**: Verifies that a `rules.json` file exists in the current directory
2. **Parse Dependencies**: Reads and parses the rules from the `rules.json` file
3. **Display Tree**: Shows the rules in a hierarchical tree format with versions

## Output Format

### Default Tree Format

```
new-rules@1.0.0
├── redis@0.0.1
├── workos/authkit-nextjs@0.0.1
└── gh:owner/repo@0.0.1
```

## Error Cases

### No rules.json Found

```
Error: No rules.json file found in current directory
Run 'rules init' to initialize a new project
```

**Exit Code**: 1

### Invalid rules.json Format

```
Error: Invalid rules.json format
Failed to parse JSON: <error details>
```

**Exit Code**: 1

### Empty Rules

When no rules are installed:

```
new-rules@1.0.0
(empty)
```

## Examples

### Basic Usage

```bash
$ rules list
new-rules@1.0.0
├── redis@0.0.1
├── workos/authkit-nextjs@0.0.1
└── gh:owner/repo@0.0.1
```

## Implementation Notes

- Should be fast and not require network calls (only reads local files)
- Should respect terminal width for formatting
