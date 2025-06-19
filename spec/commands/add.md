# `rules add`

Adds a rule to the project.

## Usage

```bash
rules add vercel/nextjs
rules add gh:owner/repo
```

## Args

- Name of ruleset to add (with optional `gh:` prefix for GitHub repositories)

## Behavior

- If rules.json doesn't exist, creates it with default structure and adds the rule
- Downloads rule files from the registry to appropriate folder (e.g. `.rules/vercel/nextjs/`) using the [registry API GET endpoint](../registry-api.md#get)
- Adds the rule to rules.json "rules" object with the literal version that was downloaded
- For GitHub repos (`gh:` prefix):
  - Downloads all files in the repository
  - Uses the main branch by default
  - Looks for rules.json in the downloaded files to find the version, just like with the normal `add` command
- When rules.json doesn't exist:
  - Checks for any top-level folder of the structure ".{folder-name}/rules"
  - If one exists, suggests to the user to run `rules render {folder-name}`