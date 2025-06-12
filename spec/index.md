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
│   ├── publish.go      # Publish command
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
  "name": "new-rules",
  "description": "",
  "author": "",
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

## Commands

### `rules init`

Creates initial rule directory structure.

- **Flags**:
  - `--format string`: Set rule format (default, cursor, etc.)
- **Behavior**:
  - Creates directory structure
  - Initializes empty rules.json
  - Sets up format-specific configuration
  - Default format creates `.rules/` directory
  - Custom formats create `.{format}/rules/` directories

### `rules create`

Creates a new rule file in the rules directory.

```bash
# Create new rules with interactive walkthrough that lets you choose triggers and write rules
rules create
rules create --tags frontend --globs *.{tsx,jsx} --description "Style guide for writing React components" "This is the body of the rule"
rules create --alwaysApply # Body not supplied, so will prompt for it interactively
```

- **Flags**:
  - `--tags`: Comma-separated list of tags
  - `--globs`: Glob patterns to match files
  - `--description`: Short description
  - `--alwaysApply`: Flag to always apply rule
- **Args**:
  - Optional rule body as last argument
- **Behavior**:
  - Prompts for missing fields if not provided
  - Allows for stdin/editor input for rule body
  - Creates a new rule (.md) file in the root of the rules directory
  - Does not modify the rules.json file

### `rules add`

Adds a rule to the project.

```bash
rules add vercel/nextjs
rules add gh:owner/repo
```

- **Args**:
  - Name of ruleset to add (with optional `gh:` prefix for GitHub repositories)
- **Behavior**:
  - Adds the rule to rules.json "rules" object
  - If rules.json doesn't exist, creates it with default structure and adds the rule
  - Downloads rule files from the registry to appropriate folder (e.g. `.rules/vercel/nextjs/`) using the [registry API GET endpoint](registry-api.md#get)
  - For GitHub repos (`gh:` prefix):
    - Downloads all files from the `src/` directory in the repository
    - Uses the main branch by default
  - When rules.json doesn't exist:
    - Checks for any top-level folder of the structure ".{folder-name}/rules"
    - If one exists, suggests to the user to run `rules render {folder-name}`

### `rules remove`

Removes a rule from the project.

```bash
rules remove vercel/nextjs
rules remove gh:owner/repo
```

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

Renders existing rules to a specified format.

```bash
rules render foo
rules render cursor
rules render --all  # Renders to all formats specified in config
```

- **Args**:
  - Name of format to render rules to (e.g., "foo", "continue")
- **Flags**:
  - `--all`: Renders to all formats specified in config
- **Behavior**:
  - Creates `.{format}/rules/` directory (e.g., `.foo/rules/`)
  - Copies all rules from the default location (`.rules/`) to the target format location
  - Preserves directory structure of rule sets
  - Can perform format-specific transformations if needed
  - Does not modify the original rule files

### `rules install`

Synchronizes the `.rules` directory with the contents of `rules.json`.

```bash
rules install
rules install --force  # Skip confirmation prompts
```

- **Flags**:
  - `--force`: Skip confirmation prompts
- **Behavior**:
  - Performs a clean installation by:
    - Removing all existing rule files and directories first
    - Re-downloading and installing all rules specified in `rules.json`
  - Provides confirmation prompt before deleting existing rules
  - Ensures the `.rules` directory exactly matches what's defined in `rules.json`
  - Reports on installation progress and any errors encountered

### `rules publish`

Publishes a rule file to the registry.

```bash
rules publish my-rule.md --slug my-org/my-rules
rules publish my-rule.md --slug my-org/my-rules --visibility private
```

- **Args**:
  - Path to the rule file to publish
- **Flags**:
  - `--slug`: The organization/ruleset slug to publish to (required)
  - `--visibility`: Set the visibility of the rule to "public" or "private" (default: "public")
- **Behavior**:
  - Reads the content of the specified rule file
  - Uses the registry API's POST endpoint to publish the rule
  - Requires user to be logged in (uses Bearer auth)
  - Sets the visibility of the published rule according to the flag
  - Returns a confirmation message with the published rule's details

### `rules login`

Starts the authorization flow using utilities defined in [the auth folder](../internal/auth/) and saves auth information.

### `rules logout`

Logs the user out by removing the auth file. Use utilities defined in [the auth folder](../internal/auth/).

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
