# Authentication

This CLI program is authenticated using the Device Code Flow with the WorkOS Go SDK. The flow looks like this:

```
# 1. Device requests codes
POST /device_authorization HTTP/1.1
client_id=abc123&scope=example_scope

# 2. Server response
{
  device_code: "abcxyz...",
  user_code: "ABCD-EFGH",
  verification_uri: "https://example.com/device",
  expires_in: 900,
  interval: 5
}

# 3. Show user prompt
"Visit https://example.com/device and enter code ABCD-EFGH"

# 4. Poll for token
POST /oauth/token
grant_type=device_code
device_code=abcxyz...
```
