# Go CLI Specification: "Rules" Tool

## Overview

A command-line tool to create, manage, and convert rule sets for code guidance across different AI assistant platforms (Continue, Cursor, Windsurf, Copilot, etc.). The tool allows for creating rule sets in different formats and locations, managing rules, and publishing them to a central registry.

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
│   ├── list.go       # List command
│   ├── render.go       # Render command
│   ├── install.go      # Install command
│   ├── publish.go      # Publish command
│   ├── whoami.go       # Whoami command
├── internal/
│   ├── config/         # Configuration management
│   ├── formats/        # Format handling
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

### Core Commands

- [`rules init`](commands/init.md) - Creates initial rule directory structure
- [`rules create`](commands/create.md) - Creates a new rule file in the rules directory
- [`rules add`](commands/add.md) - Adds a rule to the project
- [`rules remove`](commands/remove.md) - Removes a rule from the project
- [`rules list`](commands/list.md) - Lists all rules currently installed in the project
- [`rules render`](commands/render.md) - Renders existing rules to a specified format
- [`rules install`](commands/install.md) - Synchronizes the `.rules` directory with the contents of `rules.json`

### Registry Commands

- [`rules publish`](commands/publish.md) - Publishes a rule file to the registry
- [`rules whoami`](commands/whoami.md) - Displays information about the currently authenticated user
- [`rules login`](commands/login.md) - Starts the authorization flow and saves auth information
- [`rules logout`](commands/logout.md) - Logs the user out by removing the auth file

## Error Handling

- Clear error messages with actionable advice
- Validation checks before operations
- Proper exit codes for different error scenarios

## Configuration

- Uses Viper for configuration management
- Supports environment variables
- Configuration file stored in user's config directory
