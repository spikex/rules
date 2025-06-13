# Render Formats

This document provides information about the formats supported by the `rules render` command.

## Standard Format

Our standard format is a superset of all capabilities. When users download rules, they will always be in this format. When users run `rules render <format-id>`, we help them translate all of their rules to the desired format.

Most formats require just placing the markdown files in a new folder, maintaining directory structure, and making slight adjustments to the frontmatter of each markdown file. However, for formats that support only a single file (like AGENT.md), we concatenate all rules that have `alwaysApply: true` into a single file at the root of the repository, and ignore all other rules.

## Supported Formats

### Continue

- **ID**: "continue"
- **Folder**: `.continue/rules`
- **File Extension**: `.md`
- **Format**:

  ```md
  ---
  alwaysApply: name
  description: This is a description
  globs: "**/*.tsx"
  ---

  # Markdown content
  ```

### Cursor

- **ID**: "cursor"
- **Folder**: `.cursor/rules`
- **File Extension**: `.mdc`
- **Format**: Same format as Continue

### Windsurf

- **ID**: "windsurf"
- **Folder**: `.windsurf/rules`
- **File Extension**: `.md`
- **Format**:

  ```md
  ---
  trigger: [manual, always_on]
  description: This is a description
  globs: "**/*.tsx"
  ---

  # Markdown content
  ```

### Claude Code

- **ID**: "claude"
- **File**: `CLAUDE.md`
- **Format**: Markdown only (single file)

### Copilot

- **ID**: "copilot"
- **Folder**: `.github/instructions`
- **File Extension**: `.instructions.md`
- **Format**:

  ```md
  ---
  applyTo: "**"
  ---

  # Markdown content
  ```

### Codex

- **ID**: "codex"
- **File**: `AGENT.md`
- **Format**: Markdown only (single file)

### Cline

- **ID**: "cline"
- **Folder**: `.clinerules`
- **File Extension**: `.md`
- **Format**: Markdown only

### Cody

- **ID**: "cody"
- **Folder**: `.sourcegraph`
- **File Extension**: `.rule.md`
- **Format**: Markdown only

### Amp

- **ID**: "amp"
- **File**: `AGENT.md`
- **Format**: Markdown only (single file)
