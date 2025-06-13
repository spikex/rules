# Private Homebrew Tap for Continue Dev Tools

This is a private Homebrew tap for internal Continue Dev tools.

## Installation Options

### Option 1: Using NPM (Cross-platform)

You can install and run the rules CLI using npm:

```bash
# Install globally
npm install -g rules-cli

# Run the CLI
rules-cli

# Or run directly without installing
npx rules-cli
```

The npm package provides binaries for:

- macOS (x64, arm64)
- Linux (x64, arm64)
- Windows (x64)

### Option 2: Using Homebrew (macOS/Linux)

#### Prerequisites

- You must have [Homebrew](https://brew.sh/) installed
- You need Git SSH access to the private repositories at Continue Dev

#### 1. Tap the Repository

```bash
# Using SSH (recommended for private repos)
brew tap continuedev/tap git@github.com:continuedev/homebrew-tap.git
```

#### 2. Install the Rules CLI Tool

```bash
brew install rules
```

#### Updating

To get the latest version:

```bash
brew update
brew upgrade rules
```

## Available Tools

- `rules`: Internal rules CLI tool for Continue Dev

## Troubleshooting

If you encounter issues:

- Ensure you have SSH access to the `continuedev/rules-cli` repository
- Make sure your SSH key is loaded (`ssh-add -l` to check)
- Try running `brew doctor` to check for general Homebrew issues
- If the formula fails to install, try with verbose output: `brew install -v rules`
