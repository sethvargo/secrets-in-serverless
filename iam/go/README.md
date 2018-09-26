# IAM

This example shows using IAM to authenticate a cloud function directly without
exchanging credentials.

## Create Service Account

Create the service account

```text
$ gcloud iam service-accounts create my-instance-readonly
```

Grant permission to access IAM roles to the service account

```text
$ gcloud bigtable instances add-iam-policy-binding my-instance \
    --member serviceAccount:my-instance-readonly@${GOOGLE_CLOUD_PROJECT}.iam.gserviceaccount.com \
    --role roles/bigtable.reader
```

## Deploy

Deploy the function with the service account attached

```text
$ gcloud alpha functions deploy envvars \
    --runtime go111 \
    --entry-point F \
    --service-account my-instance-readonly@${GOOGLE_CLOUD_PROJECT}.iam.gserviceaccount.com \
    --trigger-http
```
