# `rules whoami`

Displays information about the currently authenticated user.

## Usage

```bash
rules whoami
```

## Behavior

- Checks if the user is currently logged in
- If logged in, displays user information (username, email, organization)
- If not logged in, displays a message indicating the user is not authenticated
- Uses the authentication information stored in the auth file
- May make a request to the registry API to fetch the latest user information