# `rules create`

Creates a new rule file in the rules directory using the Continue format specification.

## Usage

```bash
# Create new rules with interactive walkthrough that lets you choose triggers and write rules
rules create
rules create --tags frontend --globs "**/*.{tsx,jsx}" --description "Style guide for writing React components" "This is the body of the rule"
rules create --alwaysApply "Always apply this rule" # Body not supplied, so will prompt for it interactively
```

## Flags

- `--tags`: Comma-separated list of tags
- `--globs`: Glob patterns to match files (defaults to "\*\*/\*" if not specified)
- `--description`: Short description
- `--alwaysApply`: Flag to always apply rule (creates `alwaysApply: true` in frontmatter)

## Args

- Optional rule name as first argument (if not provided, will be prompted)
- Optional rule body as second argument (if not provided, will be prompted interactively)

## Behavior

- Prompts for missing fields if not provided
- Allows for stdin/editor input for rule body
- Creates a new rule (.md) file in the current directory following Continue format
- Uses Continue frontmatter format with `alwaysApply`, `description`, and `globs` fields
- Does not modify the rules.json file

## Continue Format Output

The generated rule files follow the Continue format specification (see [render-formats.md](../render-formats.md) for full format details):

```md
---
alwaysApply: true
description: This is a description
globs: "**/*.tsx"
---

# Rule content here
```

### Frontmatter Fields:

- `alwaysApply`: boolean - Whether to always apply the rule (set via `--alwaysApply` flag)
- `description`: string - Short description of the rule (set via `--description` flag)
- `globs`: string - Glob patterns to match files (set via `--globs` flag, defaults to "\*\*/\*")

Note: Tags are not part of the Continue format specification and will be omitted from the frontmatter.
