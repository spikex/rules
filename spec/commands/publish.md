# `rules publish`

Publishes a rule file to the registry.

## Usage

```bash
rules publish                    # Publish from current directory
rules publish ./my-rules         # Publish from specified directory
rules publish --visibility private
```

## Args

- Optional path to directory containing rules.json (defaults to current directory)

## Flags

- `--visibility`: Set the visibility of the rule to "public" or "private" (default: "public")

## Behavior

- Reads the slug from rules.json in the current directory or specified path
- The slug is constructed as `{organization}/{ruleset-name}` where:
  - `organization` is determined from the authenticated user's organization slug, username, or email prefix
  - `ruleset-name` is the "name" field from rules.json
- Automatically finds the main rule file to publish (index.md or first .md file found)
- Uses the registry API's POST endpoint to publish the rule
- Requires user to be logged in (uses Bearer auth)
- Sets the visibility of the published rule according to the flag
- Returns a confirmation message with the published rule's details, including the URL where the rule is available