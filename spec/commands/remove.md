# `rules remove`

Removes a rule from the project.

## Usage

```bash
rules remove vercel/nextjs
rules remove gh:owner/repo
```

## Args

- Name of ruleset to remove (including GitHub-sourced rules)

## Behavior

- Removes rule reference from rules.json
- Deletes rule files from `.rules` folder