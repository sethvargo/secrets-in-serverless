# HashiCorp

This example shows how to access secrets in [HashiCorp Vault][hashicorp-vault],
and then access them inside a serverless [Cloud Function][gcp-func].


## Setup

[IAM for Google Cloud Functions][gcf-iam] is not yet generally available. You
may need to [request access][gcf-iam-eap] if you are not part of the EAP.

If you have not previously used cloud functions or cloud storage, enable the
APIs:

```text
$ gcloud services enable \
    cloudfunctions.googleapis.com \
    iam.googleapis.com
```

## Create Vault IAM User

Vault itself needs the ability to communicate with the Google Cloud API to
validate identity. Create a dedicated service account for Vault:

```text
$ gcloud iam service-accounts create vault-verifier
```

Grant the service account the ability to verify service accounts:

```text
$ gcloud projects add-iam-policy-binding ${GOOGLE_CLOUD_PROJECT} \
    --member=serviceAccount:vault-verifier@${GOOGLE_CLOUD_PROJECT}.iam.gserviceaccount.com \
    --role=roles/iam.serviceAccountUser
```

This IAM user can be attached directly to the GCE/GKE instances on which Vault
is running, or it can be provided to Vault as a configuration parameter. For
simplicity, this guide supplies the credentials directly to Vault.


## Deploy Vault

Deploying Vault is out of scope for this guide. Please see
[sethvargo/vault-on-gke](https://github.com/sethvargo/vault-on-gke) or
[kelseyhightower/vault-on-google-kubernetes-engine](https://github.com/kelseyhightower/vault-on-google-kubernetes-engine)
for detailed steps.

For the purpose of this exercise, the Cloud Function needs to access a
credential at `secret/apikeys/twitter` which contains the API key for
communicating with the Twitter API.

```text
$ vault write secret/apikeys/twitter value=abcd1234...
```


## Enable Authentication

Enable Google Cloud authentication to HashiCorp Vault. This enables Google Cloud
entities, including Cloud Functions, to authentication to Vault using their
instance metadata.

```text
$ vault auth enable \
    -default-lease-ttl=1m \
    -max-lease-ttl=10m \
    -token-type=batch \
    gcp
```

- `-default-lease-ttl` is set to 1 minute because that is the default function
  timeout.

- `-max-lease-ttl` is set to 10 minutes because that is the maximum allowed
  function execution time

- `-token-type` is set to batch, which allows for much more scalability than
  service tokens

Finally, configure the authenticate with permission to validate other service
accounts (which is how our Cloud Function will authenticate).

```text
$ vault write auth/gcp/config \
    credentials="$(gcloud iam service-accounts keys create - --iam-account=vault-verifier@${GOOGLE_CLOUD_PROJECT}.iam.gserviceaccount.com)"
```


## Create Vault Policy

Create a policy in Vault that permits retrieving the Twitter API key value
created above. Vault will assign this policy to the Cloud Function's
authentication, allowing it to retrieve the value.

```text
$ vault policy write apikey-twitter -<<EOF
path "secret/apikeys/twitter" {
  capabilities = ["read"]
}
EOF
```


## Create Service Account and Vault Role

Create a new service account which will be attached to the Cloud Function at
boot.

```text
$ gcloud iam service-accounts create app1-vault-auther
```

Create a role that permits this service account to authenticate to Vault. Upon
success, Vault will assign the policy just created to the resulting token.

```text
$ vault write auth/gcp/role/socialmedia \
    type=iam \
    project_id=${GOOGLE_CLOUD_PROJECT} \
    policies=apikey-twitter \
    bound_service_accounts=app1-vault-auther@${GOOGLE_CLOUD_PROJECT}.iam.gserviceaccount.com \
    max_jwt_exp=60m
```


## Deploy

Deploy the function your language of choice. Be sure to attach the service
account.

### Python

```text
$ gcloud beta functions deploy vault-batch \
    --source ./python \
    --runtime python37 \
    --entry-point F \
    --service-account app1-vault-auther@${GOOGLE_CLOUD_PROJECT}.iam.gserviceaccount.com \
    --set-env-vars VAULT_ADDR=... \
    --trigger-http
```

### Node

```text
$ gcloud beta functions deploy vault-batch \
    --source ./node \
    --runtime nodejs8 \
    --entry-point F \
    --service-account app1-vault-auther@${GOOGLE_CLOUD_PROJECT}.iam.gserviceaccount.com \
    --set-env-vars VAULT_ADDR=... \
    --trigger-http
```

### Go

```text
$ gcloud beta functions deploy vault-batch \
    --source ./go \
    --runtime go111 \
    --entry-point F \
    --service-account app1-vault-auther@${GOOGLE_CLOUD_PROJECT}.iam.gserviceaccount.com \
    --set-env-vars VAULT_ADDR=... \
    --trigger-http
```


## Invoke

Invoke the cloud function at its invoke endpoint:

```text
$ gcloud functions call vault-batch
abcd1234...
```

[gcp-func]: https://cloud.google.com/functions/
[gcf-iam-eap]: https://bit.ly/gcf-iam-alpha
[gcf-iam]: https://cloud.google.com/functions/docs/securing/managing-access
[hashicorp-vault]: https://www.vaultproject.io
