# Authentication

By default, the CLI uses a device code flow for manual authentication with the `rules login`, `rules logout`, and `rules whoami` commands. But there is also an option to supply a `CONTINUE_API_KEY` environment variable that will override the authentication details that would otherwise be stored locally after running `rules login`.
