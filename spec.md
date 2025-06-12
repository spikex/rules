# Go CLI Specification: "Rules" Tool

## Overview

A command-line tool to create, manage, and convert rule sets for code guidance across different AI assistant platforms (Continue, Cursor, Claude Code, Copilot, etc.). The tool allows for creating rule sets in different formats and locations, managing rules, and publishing them to a central registry.

## Technologies & Dependencies

1. **Language**: Go (1.20+)
2. **CLI Framework**: [Cobra](https://github.com/spf13/cobra) for command structure
3. **Configuration Management**: [Viper](https://github.com/spf13/viper) for configuration files
4. **File Operations**: Standard Go libraries (os, io/ioutil)
5. **JSON Handling**: encoding/json for parsing and writing rule configs
6. **Interactive Prompts**: [promptui](https://github.com/manifoldco/promptui) for interactive rule creation
7. **HTTP Client**: net/http for API calls to registry

## Project Structure

```
rules-cli/
├── cmd/
│ ├── root.go         # Main command definition
│   ├── init.go         # Init command
│   ├── create.go       # Create command
│   ├── add.go          # Add command
│   ├── remove.go       # Remove command
│   ├── render.go       # Render command
│   ├── install.go      # Install command
├── internal/
│   ├── config/         # Configuration management
│   ├── formats/        # Format handling (cursor, default, etc)
│   ├── registry/       # Registry client
│   ├── ruleset/        # Rule set management
│   └── generators/     # Rule generators
├── main.go             # Application entry point
├── go.mod              # Go module definition
└── go.sum              # Go module checksums
```

## Data Structures

### rules.json

The `rules.json` file goes in the root of the project, adjacent to the rules directory

```json
{
  "name": "ruleset-name",
  "description": "Description of the ruleset",
  "author": "Author Name",
  "license": "Apache-2.0",
  "version": "1.0.0",
  "rules": {
    "redis": "0.0.1",
    "workos/authkit-nextjs": "0.0.1",
    "gh:owner/repo": "0.0.1"
  }
}
```

### Rule file format (.md with front matter)

```md
---
# All of these fields are optional
description: Description of the rule
tags: [tag1, tag2]
globs: *.{jsx,tsx}
alwaysApply: false
---

This is the body of the rule. It supports Markdown syntax.
```

## Core Functionality

### 1. Rule Initialization (rules init)

- Creates initial rule directory structure
- Supports different formats via `--format` flag
- Default format creates `.rules/` directory
- Custom formats create `.{format}/rules/` directories
- Creates empty rules.json with basic structure

### 2. Rule Creation (rules create)

```bash
# Create new rules with interactive walkthrough that lets you choose triggers and write rules
rules create
rules create --tags frontend --globs *.{tsx,jsx} --description "Style guide for writing React components" "This is the body of the rule"
rules create --alwaysApply # Body not supplied, so will prompt for it interactively
```

- Interactive mode when parameters not supplied
- Supports flags for all rule properties (tags, globs, description, alwaysApply)
- Allows for stdin/editor input for rule body
- Creates a new rule (.md) file in the root of the rules directory
- Does not modify the rules.json file at all

### 3. Rule Importing (rules add)

```bash
rules add vercel/nextjs
rules add gh:owner/repo
```

- Adds the rule to rules.json "rules" object
- If rules.json doesn't exist, creates it with default structure and adds the rule
- Downloads rule files from the registry to appropriate folder (e.g. `.rules/vercel/nextjs/`)
- When using `gh:` prefix, downloads the rules from the GitHub repository:
  - By default, imports all files from the `src/` directory in the repository
  - Downloads from the main branch of the repository
- When rules.json doesn't exist:
  - Check for any top-level folder of the structure ".{folder-name}/rules"
  - If one exists, print a suggestion to the user to run `rules render {folder-name}` at the very end of the output

### 4. Rule Removal (rules remove)

```bash
rules remove vercel/nextjs
rules remove gh:owner/repo
```

- Removes the rule from rules.json "rules" object
- Optionally deletes rule files from the local directory (with --delete flag)
- Provides confirmation prompt before deletion (can be bypassed with --force flag)

### 5. Rule Rendering (rules render)

```bash
rules render foo
rules render cursor
rules render --all  # Renders to all formats specified in config
```

- Renders existing rules to a specified format
- Creates `.{format}/rules/` directory (e.g., `.foo/rules/`)
- Copies all rules from the default location (`.rules/`) to the target format location
- Can render to multiple formats simultaneously with `--all` flag
- Preserves directory structure of rule sets

### 6. Rule Installation (rules install)

```bash
rules install
rules install --force  # Skip confirmation prompts
```

- Synchronizes the `.rules` directory with the contents of `rules.json`
- Performs a clean installation by:
  - Removing all existing rule files and directories first
  - Re-downloading and installing all rules specified in `rules.json`
- Provides confirmation prompt before deleting existing rules (can be bypassed with --force flag)
- Useful for ensuring the rules directory matches exactly what's defined in `rules.json`

## Command Specifications

### `rules init`

- **Flags**:
  - `--format string`: Set rule format (default, cursor, etc.)
- **Behavior**:
  - Creates directory structure
  - Initializes empty rules.json
  - Sets up format-specific configuration

### `rules create`

- **Flags**:
  - `--tags`: Comma-separated list of tags
  - `--globs`: Glob patterns to match files
  - `--description`: Short description
  - `--alwaysApply`: Flag to always apply rule
- **Args**:
  - Optional rule body as last argument
- **Behavior**:
  - Prompts for missing fields if not provided
  - Creates rule file in root of the rules directory
  - Does not modify the rules.json file

### `rules add`

- **Args**:
  - Name of ruleset to add (with optional `gh:` prefix for GitHub repositories)
- **Behavior**:
  - Checks if rules.json exists:
    - If it exists, adds the rule to the existing file
    - If it doesn't exist, creates a new rules.json with default structure and adds the rule
  - Fetches ruleset from registry or GitHub based on prefix
  - For GitHub repos:
    - Downloads all files from the `src/` directory in the repository
    - Uses the main branch by default
  - Adds to rules.json "rules" object
  - Validates ruleset exists
  - When rules.json doesn't exist:
    - Check for any top-level folder of the structure ".{folder-name}/rules"
    - If one exists, print a suggestion to the user to run `rules render {folder-name}` at the very end of the output

### `rules remove`

- **Args**:
  - Name of ruleset to remove (including GitHub-sourced rules)
- **Flags**:
  - `--delete`: Also delete rule files from disk
  - `--force`: Skip confirmation prompts
- **Behavior**:
  - Removes rule reference from rules.json
  - Optionally deletes rule files from disk
  - Confirms before destructive operations

### `rules render`

- **Args**:
  - Name of format to render rules to (e.g., "foo", "continue")
- **Behavior**:
  - Creates `.{format}/rules/` directory structure
  - Copies all rules from the default location (`.rules/`) to the target format location
  - Preserves directory structure of rule sets
  - Can perform format-specific transformations if needed
  - Does not modify the original rule files

### `rules install`

- **Flags**:
  - `--force`: Skip confirmation prompts
- **Behavior**:
  - Removes all existing rule files and directories from `.rules/`
  - Re-downloads and installs all rules specified in `rules.json`
  - Confirms before deleting existing rules
  - Ensures the `.rules` directory exactly matches the specification in `rules.json`
  - Reports on installation progress and any errors encountered

## Error Handling

- Clear error messages with actionable advice
- Validation checks before operations
- Proper exit codes for different error scenarios

## Configuration

- Uses Viper for configuration management
- Supports environment variables
- Configuration file stored in user's config directory

# Task

Please implement the entire command line tool described above in the current directory.
