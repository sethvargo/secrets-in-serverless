# Environment Variables

This example shows how to store secrets in environment variables in a serverless
[Cloud Function][gcp-func].

**To be absolutely clear, you should not do this under any circumstances. This
is for illustration purposes only.**


## Setup

If you have not previously used cloud functions, enable the API:

```text
$ gcloud services enable cloudfunctions.googleapis.com
```


## Deploy

```text
$ gcloud alpha functions deploy envvars \
    --runtime go111 \
    --entry-point F \
    --set-env-vars DB_USER=my-user,DB_PASS=s3cr3t \
    --trigger-http
```


## Invoke

Invoke the cloud function at its invoke endpoint:

```text
$ open $(gcloud alpha functions describe envvars --format='value(httpsTrigger.url)')
```

[gcp-func]: https://cloud.google.com/functions/
