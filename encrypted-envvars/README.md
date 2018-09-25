# Encrypted Environment Variables

This example shows using encrypted environment variables in a serverless function.


## Encrypt Values

First, you'll need to encrypt all your values with a KMS key:

```text
$ echo "s3cr3t" | gcloud kms encrypt \
    --location=global \
    --keyring=serverless-apps \
    --key=app1 \
    --ciphertext-file=- \
    --plaintext-file=- \
    | base64
```

Do this for all your plaintext values and save the ciphertext.


## Grant IAM Permissions

Next, you'll need to grant IAM permission to the Google Cloud function to
decrypt these values. Be sure to replace `GOOGLE_CLOUD_PROJECT` with your
project name.

```text
$ gcloud kms keys add-iam-policy-binding app1 \
    --location global \
    --keyring serverless-apps \
    --member "serviceAccount:${GOOGLE_CLOUD_PROJECT}@appspot.gserviceaccount.com" \
    --role roles/cloudkms.cryptoKeyDecrypter
```


## Deploy

Finally, deploy the function. Populate all the environment variables with the
encrypted values you received earlier.

```text
$ gcloud beta functions deploy encrypted-envvars \
    --runtime go111 \
    --entry-point F \
    --set-env-vars KMS_CRYPTO_KEY_ID=projects/my-project/locations/global/keyRings/serverless-apps/cryptoKeys/app1,DB_USER=CiQAePa3VEjcuknRhLX...,DB_PASS=CiQAePa3VEpDBjS2ac... \
    --trigger-http
```
