# Environment Variables

This example shows how to store secrets in environment variables in a serverless
[Cloud Function][gcp-func].

**To be absolutely clear, you should not do this under any circumstances. This
is for illustration purposes only.**


## Setup

[IAM for Google Cloud Functions][gcf-iam] is not yet generally available. You
may need to [request access][gcf-iam-eap] if you are not part of the EAP.

If you have not previously used cloud functions, enable the API:

```text
$ gcloud services enable cloudfunctions.googleapis.com
```


## Deploy

Deploy the function your language of choice.

### Python

```text
$ gcloud alpha functions deploy envvars \
    --source ./python \
    --runtime python37 \
    --entry-point F \
    --set-env-vars DB_USER=my-user,DB_PASS=s3cr3t \
    --trigger-http
```

### Node

```text
$ gcloud alpha functions deploy envvars \
    --source ./node \
    --runtime nodejs8 \
    --entry-point F \
    --set-env-vars DB_USER=my-user,DB_PASS=s3cr3t \
    --trigger-http
```

### Go

```text
$ gcloud alpha functions deploy envvars \
    --source ./go \
    --runtime go111 \
    --entry-point F \
    --set-env-vars DB_USER=my-user,DB_PASS=s3cr3t \
    --trigger-http
```


## Authorize

Make the function publicly accessible. This step is optional, but you will be
unable to invoke the Cloud Function without adding this or a similar IAM policy.

```text
$ gcloud alpha functions add-iam-policy-binding envvars \
    --member allUsers \
    --role roles/cloudfunctions.invoker
```


## Invoke

Invoke the cloud function at its invoke endpoint:

```text
$ curl $(gcloud alpha functions describe envvars --format='value(httpsTrigger.url)')
my-user:s3cr3t
```

[gcp-func]: https://cloud.google.com/functions/
[gcf-iam-eap]: https://bit.ly/gcf-iam-alpha
[gcf-iam]: https://cloud.google.com/functions/docs/securing/managing-access
