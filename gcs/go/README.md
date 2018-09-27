# Google Cloud Storage

This example shows how to encrypt secrets on [Google Cloud Storage][gcp-gcs],
and then access them inside a serverless [Cloud Function][gcp-func].


## Setup

If you have not previously used cloud functions or cloud storage, enable the
APIs:

```text
$ gcloud services enable \
    cloudfunctions.googleapis.com \
    storage-component.googleapis.com
```


## Create Bucket

Create a new bucket specifically for storing secrets.

```text
$ gsutil mb gs://${GOOGLE_CLOUD_PROJECT}-serverless-secrets
```

You can also optionally enable object versioning and retention policies, but
that is not specifically discussed here.

Next, change the default bucket permissions. By default, anyone with access
to the project has access to the data in the bucket. You must do this **before
storing any data in the bucket!**

```text
$ gsutil defacl set private gs://${GOOGLE_CLOUD_PROJECT}-serverless-secrets
$ gsutil acl set -r private gs://${GOOGLE_CLOUD_PROJECT}-serverless-secrets
```


## Store Secrets

You need to store some secrets in the Cloud Storage bucket for the application
to access. By default, the app looks for secrets in a file named "app1".

```text
$ gsutil -h 'Content-Type: application/json' cp - gs://${GOOGLE_CLOUD_PROJECT}-serverless-secrets/app1 <<< '{"username":"my-user", "password":"s3cr3t"}'
```


## Grant IAM Permissions

You need to grant IAM permission to the Google Cloud function to access these
values.

Create a new service account.

```text
$ gcloud iam service-accounts create app1-gcs-reader
```

Grant the most minimal set of permissions to access data. Be sure to replace
`GOOGLE_CLOUD_PROJECT` with your project name.

```text
$ gsutil iam ch serviceAccount:app1-gcs-reader@${GOOGLE_CLOUD_PROJECT}.iam.gserviceaccount.com:legacyObjectReader \
    gs://${GOOGLE_CLOUD_PROJECT}-serverless-secrets/app1
```


## Deploy

Deploy the function, with the attached service account. Populate the environment
variables with the name of the GCS bucket.

```text
$ gcloud alpha functions deploy gcs \
    --runtime go111 \
    --entry-point F \
    --service-account app1-gcs-reader@${GOOGLE_CLOUD_PROJECT}.iam.gserviceaccount.com \
    --set-env-vars STORAGE_BUCKET=${GOOGLE_CLOUD_PROJECT}-serverless-secrets \
    --trigger-http
```


## Invoke

Invoke the cloud function at its invoke endpoint:

```text
$ open $(gcloud alpha functions describe gcs --format='value(httpsTrigger.url)')
```

[gcp-gcs]: https://cloud.google.com/storage
[gcp-func]: https://cloud.google.com/functions/
