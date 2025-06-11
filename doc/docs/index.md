---
title: rules
---

# rules

`rules` is a CLI built for managing rules across any AI developer tool. Rules are markdown files that encode workflows, preferences, tech stack details, and more in plain natural language so you can get better help from LLMs.

## Install `rules`

To install the `rules` CLI, run the following command on Mac. If you don't have `brew` installed yet, you can find the command to do so [here](https://brew.sh/).

```bash
brew install rules
```

## Add rules

To download rules, just run (for example)

```bash
rules add vercel/nextjs
```

This will add them to your project in a local `.rules` folder. If you are already using rules with an AI developer tool like Continue, Cursor, Windsurf, Cline, etc. then it will also automatically render them in the correct folder (e.g. `.continue/rules`, `.cursor/rules`, etc.).

You can also download rules from a GitHub repository:

```bash
rules add gh:my-username/my-repository
```

## Publish rules

The simplest thing you can publish is a single markdown file

```bash
rules publish my-rule.md
```

You can also publish a folder of markdown files

```bash
rules publish ./my-rules
```

The slug of the rules are defined in a `rules.json` file.

## Integration with AI IDEs

If you are building a developer tool and want to optimize how AI IDEs use your tool, `rules` makes it easy to give your users the best experience.

1. Publish rules
2. Add `rules` to your documentation
