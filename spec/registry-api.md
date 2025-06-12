# Registry API

The following curl commands and example outputs describe the specification for the registry API.

## GET

Request:

```bash
curl https://api.continue.dev/registry/v1/<owner-slug>/<rule-slug>/latest
```

Response:

```bash
{"content":"<content of the rule>"}
```

## POST

Request:

```bash
curl -X POST https://api.continue.dev/packages/<owner-slug>/<rule-slug>/versions/new -d '{
  "visibility": "public",
  "content": "this is the body of the rule"
}'
```

Reponse:

```bash

```

## Authorization

The registry API uses Bearer auth, with a header like `Authorization: Bearer xxx`. It should only be included if the user is logged in.
