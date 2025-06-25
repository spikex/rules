# Contributing to Rules CLI

## Development Setup

### Prerequisites

- Go 1.24.2+
- Node.js 12+ (for packaging)

### Local Development

```bash
# Clone the repository
git clone https://github.com/continuedev/rules.git
cd rules

# Install dependencies
go mod download

# Build the CLI
go build

# Test your changes
go test ./...
```

## Making Changes

1. **Fork the repository** and create a feature branch
2. **Make your changes** following these guidelines:
   - Keep commits focused and atomic
   - Add tests for new functionality
   - Update documentation if needed
3. **Test your changes** locally
4. **Include a clear description** in your pull request

## Release Process

Releases are automated via semantic-release. Use conventional commit messages:

- `fix:` for bug fixes
- `feat:` for new features
- `docs:` for documentation changes
- `BREAKING CHANGE:` for breaking changes

## Questions?

Open an issue for bugs, feature requests, or questions.
