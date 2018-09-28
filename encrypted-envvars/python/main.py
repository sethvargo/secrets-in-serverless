import base64
import os
import googleapiclient.discovery

crypto_key_id = os.environ['KMS_CRYPTO_KEY_ID']

def decrypt(client, s):
    if not s:
        raise ValueError('ciphertext is missing')

    response = kms_client \
        .projects() \
        .locations() \
        .keyRings() \
        .cryptoKeys() \
        .decrypt(name=crypto_key_id, body={"ciphertext":s}) \
        .execute()

    return base64.b64decode(response['plaintext']).decode('utf-8').strip()


kms_client = googleapiclient.discovery.build('cloudkms', 'v1')

username = decrypt(kms_client, os.environ['DB_USER'])
password = decrypt(kms_client, os.environ['DB_PASS'])

def F(request):
    return f'{username}:{password}'
