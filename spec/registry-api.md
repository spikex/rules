# Registry API

The following curl commands and example outputs describe the specification for the registry API.

## GET - Download Package

### Get Latest Version

Request:

```bash
curl https://api.continue.dev/v0/<owner-slug>/<rule-slug>/latest/download
```

### Get Specific Version

Request:

```bash
curl https://api.continue.dev/v0/<owner-slug>/<rule-slug>/<version>/download
```

Response:

Returns the zip file as binary data with appropriate headers:

```
Content-Type: application/zip
Content-Disposition: attachment; filename="<rule-slug>-<version>.zip"
Content-Length: 15420
```

## POST - Upload Package

Request:

```bash
curl -X POST https://api.continue.dev/v0/<owner-slug>/<rule-slug>/<version-slug> \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: multipart/form-data" \
  -F "file=@package.zip" \
  -F 'metadata={
    "visibility": "public",
  }'
```

## Authorization

The registry API uses Bearer auth, with a header like `Authorization: Bearer <token>`. It should only be included if the user is logged in.

- **GET** requests: No authorization required for public packages
- **POST/DELETE** requests: Authorization required
- **GET** requests for private packages: Authorization required
