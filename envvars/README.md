# Environment Variables

This example shows using unencrypted environment variables in a serverless
function.

## Deploy

```text
$ gcloud beta functions deploy envvars \
    --runtime go111 \
    --entry-point Secrets \
    --set-env-vars DB_USER=my-user,DB_PASS=s3cr3t \
    --trigger-http
```
