# `rules create`

Creates a new rule file in the rules directory.

## Usage

```bash
# Create new rules with interactive walkthrough that lets you choose triggers and write rules
rules create
rules create --tags frontend --globs *.{tsx,jsx} --description "Style guide for writing React components" "This is the body of the rule"
rules create --alwaysApply # Body not supplied, so will prompt for it interactively
```

## Flags

- `--tags`: Comma-separated list of tags
- `--globs`: Glob patterns to match files
- `--description`: Short description
- `--alwaysApply`: Flag to always apply rule

## Args

- Optional rule body as last argument

## Behavior

- Prompts for missing fields if not provided
- Allows for stdin/editor input for rule body
- Creates a new rule (.md) file in the root of the rules directory
- Does not modify the rules.json file