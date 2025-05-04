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
│   ├── root.go         # Main command definition
│   ├── init.go         # Init command
│   ├── create.go       # Create command
│   ├── add.go          # Add command
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
    "workos/authkit-nextjs": "0.0.1"
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

### 3. Rule Importing (rules add)

```bash
rules add vercel/nextjs
```

- Adds the rule to rules.json "rules" object
- Downloads rule files from the registry to appropriate folder (e.g. `.rules/vercel/nextjs/`)

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

### `rules add`

- **Args**:
  - Name of ruleset to add
- **Behavior**:
  - Fetches ruleset from registry
  - Adds to rules.json "rules" object
  - Validates ruleset exists

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
