---
title: rules
---

# rules

:::tip

**tl;dr:** `brew install rules` then `rules add vercel/nextjs`

:::

`rules` is a CLI built for managing rules across any AI developer tool. Rules are markdown files that encode workflows, preferences, tech stack details, and more in plain natural language so you can get better help from LLMs.

## Install `rules`

:::info

`rules` is currently only available privately through the `continuedev` tap. Soon it will be available publicly with `brew install rules`, and on other operating systems.

:::

To install `rules` on Mac, you can run the following command. If you don't have `brew` installed yet, you can find the command to do so [here](https://brew.sh/).

```bash
brew install continuedev/tap/rules
```

## Add rules

To download rules to your repository you can use `rules add`. For example:

```bash
rules add vercel/nextjs
```

This will add them to your project in a local `.rules` folder.

You can also download from GitHub rather than the rules registry:

```bash
rules add gh:continuedev/continue-internal-rules
```

## Render rules

To use rules with your AI code assistant of choice, you can "render" them to the necessary format and location using `rules render`. For example,

```bash
rules render continue
```

will copy all of the `.rules/` into a `.continue/rules/` folder.

## Publish rules

To make your rule available to others, you can publish a markdown file using `rules publish`:

```bash
rules publish my-rule.md --slug my-username/my-rule
```

This would make your rule available to download with `rules add my-username/my-rule`.

<!--
You can also publish a folder of markdown files:

```bash
rules publish ./my-rules --slug my-username/my-rules
``` -->

## Helping users use your rules

If you are building a developer tool and want to optimize how AI IDEs work with your tool, `rules` makes it easy to give your users the best experience.

1. Make your account on the [registry](https://hub.continue.dev/signup) and create an organization
2. [Publish your rules](index.md#publish-rules)
3. Mention the corresponding `rules add <my-rules>` command in your documentation
