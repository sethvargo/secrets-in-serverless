# Wrapped Secrets

This example shows using wrapped secrets for exactly-once decryption.


## Encrypt Values

First, you'll need to store or encrypt the data in Vault:

```text
$ vault kv put secret/dbcreds username=my-user password=s3cr3t
```

Then create a wrapped secret with the response:

```text
$ vault kv get -wrap-ttl=5m secret/dbcreds
Key               Value
---               -----
wrapping_token:   7xTpP9LMUMUysBEa7pdZNsFP
# ...
```

## Deploy

Finally, deploy the function. Populate all the environment variables with the
encrypted values you received earlier.

```text
$ gcloud beta functions deploy encrypted-envvars \
    --runtime go111 \
    --entry-point Secrets \
    --set-env-vars VAULT_ADDR=my.vault.server,VAULT_TOKEN=7xTpP9LMUMUysBEa7pdZNsFP \
    --trigger-http
```
