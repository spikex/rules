---
title: rules CLI
---

# rules

:::tip

**tl;dr:** `npm i -g rules-cli` then `rules add starter/nextjs-rules`

:::

`rules` is a CLI built for managing rules across any AI IDE. Rules are markdown files that encode workflows, preferences, tech stack details, and more in plain natural language so you can get better help from LLMs.

## Install `rules`

The `rules` CLI can be installed using NPM:

```bash
npm i -g rules-cli
```

## Add rules

To download rules to your repository you can use `rules add`. For example:

```bash
rules add starter/nextjs-rules
```

This will add them to your project in a local `.rules` folder.

You can also download from GitHub rather than the rules registry:

```bash
rules add gh:continuedev/rules-template
```

## Render rules

To use rules with your AI IDE of choice, you can "render" them to the necessary format and location using `rules render`. For example,

```bash
rules render cursor
```

will copy all of the `.rules/` into a `.cursor/rules/` folder. `rules` currently supports the following formats: cursor, continue, windsurf, claude, copilot, codex, cline, cody, and amp.

## Publish rules

To make your rules available to others, you can publish using `rules publish`:

```bash
rules login
rules publish
```

This would make your rule available to download with `rules add <name-of-rules>`.

The command automatically determines the slug from your `rules.json` file. To make sure you have a `rules.json` file in your current directory, use `rules init`.

## Helping users use your rules

If you are building a developer tool and want to optimize how AI IDEs work with your tool, `rules` makes it easy to give your users the best experience.

1. Make your account on the [registry](https://hub.continue.dev/signup) and create an organization
2. [Publish your rules](index.md#publish-rules)
3. Mention the corresponding `rules add <name-of-rules>` command in your documentation
