# Google Cloud Storage

This example shows how to encrypt secrets on [Google Cloud Storage][gcp-gcs],
and then access them inside a serverless [Cloud Function][gcp-func].


## Setup

[IAM for Google Cloud Functions][gcf-iam] is not yet generally available. You
may need to [request access][gcf-iam-eap] if you are not part of the EAP.

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

$ gsutil iam ch serviceAccount:app1-gcs-reader@${GOOGLE_CLOUD_PROJECT}.iam.gserviceaccount.com:legacyBucketReader \
    gs://${GOOGLE_CLOUD_PROJECT}-serverless-secrets
```


## Deploy

Deploy the function your language of choice. Be sure to populate the environment
with the encrypted values and attach the corresponding service account.

### Python

```text
$ gcloud beta functions deploy gcs \
    --source ./python \
    --runtime python37 \
    --entry-point F \
    --service-account app1-gcs-reader@${GOOGLE_CLOUD_PROJECT}.iam.gserviceaccount.com \
    --set-env-vars STORAGE_BUCKET=${GOOGLE_CLOUD_PROJECT}-serverless-secrets \
    --trigger-http
```

### Node

```text
$ gcloud beta functions deploy gcs \
    --source ./node \
    --runtime nodejs8 \
    --entry-point F \
    --service-account app1-gcs-reader@${GOOGLE_CLOUD_PROJECT}.iam.gserviceaccount.com \
    --set-env-vars STORAGE_BUCKET=${GOOGLE_CLOUD_PROJECT}-serverless-secrets \
    --trigger-http
```

### Go

```text
$ gcloud beta functions deploy gcs \
    --source ./go \
    --runtime go111 \
    --entry-point F \
    --service-account app1-gcs-reader@${GOOGLE_CLOUD_PROJECT}.iam.gserviceaccount.com \
    --set-env-vars STORAGE_BUCKET=${GOOGLE_CLOUD_PROJECT}-serverless-secrets \
    --trigger-http
```


## Invoke

Invoke the cloud function at its invoke endpoint:

```text
$ gcloud functions call gcs
my-user:s3cr3t
```

[gcp-gcs]: https://cloud.google.com/storage
[gcp-func]: https://cloud.google.com/functions/
[gcf-iam-eap]: https://bit.ly/gcf-iam-alpha
[gcf-iam]: https://cloud.google.com/functions/docs/securing/managing-access
