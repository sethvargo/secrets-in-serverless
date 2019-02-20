# Encrypted Environment Variables

This example shows how to encrypt secrets with [Google Cloud KMS][gcp-kms],
store them in environment variables, and then decrypt them inside a serverless
[Cloud Function][gcp-func].


## Setup

[IAM for Google Cloud Functions][gcf-iam] is not yet generally available. You
may need to [request access][gcf-iam-eap] if you are not part of the EAP.

If you have not previously used cloud functions or KMS, enable the APIs:

```text
$ gcloud services enable \
    cloudfunctions.googleapis.com \
    cloudkms.googleapis.com
```


## Encrypt Values

You need to encrypt all your values with a Google KMS key. You provide KMS the
plaintext, and KMS encrypts the value and returns ciphertext (encrypted text).

If you don't have a KMS key, you can create one first:

```text
$ gcloud kms keyrings create serverless-secrets \
    --location global

$ gcloud kms keys create app1 \
    --location global \
    --keyring serverless-secrets \
    --purpose encryption
```

Encrypt your plaintext values with this key. If you are using a different
location, keyring, or cryptokey, please use the appropriate values.

```text
$ echo "s3cr3t" | gcloud kms encrypt \
    --location=global \
    --keyring=serverless-secrets \
    --key=app1 \
    --ciphertext-file=- \
    --plaintext-file=- \
    | base64
```

Repeat this for all your plaintext values and save the ciphertext.


## Grant IAM Permissions

You need to grant IAM permission to the Google Cloud function to decrypt these
values.

Create a new service account.

```text
$ gcloud iam service-accounts create app1-kms-decrypter
```

Grant the most minimal set of permissions to decrypt data using the KMS key
created above. Be sure to replace `GOOGLE_CLOUD_PROJECT` with your project name.

```text
$ gcloud kms keys add-iam-policy-binding app1 \
    --location global \
    --keyring serverless-secrets \
    --member "serviceAccount:app1-kms-decrypter@${GOOGLE_CLOUD_PROJECT}.iam.gserviceaccount.com" \
    --role roles/cloudkms.cryptoKeyDecrypter
```


## Deploy

Deploy the function your language of choice. Be sure to populate the environment
with the encrypted values and attach the corresponding service account.

### Python

```text
$ gcloud beta functions deploy encrypted-envvars \
    --source ./python \
    --runtime python37 \
    --entry-point F \
    --service-account app1-kms-decrypter@${GOOGLE_CLOUD_PROJECT}.iam.gserviceaccount.com \
    --set-env-vars KMS_CRYPTO_KEY_ID=projects/${GOOGLE_CLOUD_PROJECT}/locations/global/keyRings/serverless-secrets/cryptoKeys/app1,DB_USER=CiQAePa3VEjcuknRhLX...,DB_PASS=CiQAePa3VEpDBjS2ac... \
    --trigger-http
```

### Node

```text
$ gcloud beta functions deploy encrypted-envvars \
    --source ./node \
    --runtime nodejs8 \
    --entry-point F \
    --service-account app1-kms-decrypter@${GOOGLE_CLOUD_PROJECT}.iam.gserviceaccount.com \
    --set-env-vars KMS_CRYPTO_KEY_ID=projects/${GOOGLE_CLOUD_PROJECT}/locations/global/keyRings/serverless-secrets/cryptoKeys/app1,DB_USER=CiQAePa3VEjcuknRhLX...,DB_PASS=CiQAePa3VEpDBjS2ac... \
    --trigger-http
```

### Go

```text
$ gcloud beta functions deploy encrypted-envvars \
    --source ./go \
    --runtime go111 \
    --entry-point F \
    --service-account app1-kms-decrypter@${GOOGLE_CLOUD_PROJECT}.iam.gserviceaccount.com \
    --set-env-vars KMS_CRYPTO_KEY_ID=projects/${GOOGLE_CLOUD_PROJECT}/locations/global/keyRings/serverless-secrets/cryptoKeys/app1,DB_USER=CiQAePa3VEjcuknRhLX...,DB_PASS=CiQAePa3VEpDBjS2ac... \
    --trigger-http
```


## Invoke

Invoke the cloud function at its invoke endpoint:

```text
$ gcloud functions call encrypted-envvars
my-user:s3cr3t
```

[gcp-kms]: https://cloud.google.com/kms
[gcp-func]: https://cloud.google.com/functions/
[gcf-iam-eap]: https://bit.ly/gcf-iam-alpha
[gcf-iam]: https://cloud.google.com/functions/docs/securing/managing-access
