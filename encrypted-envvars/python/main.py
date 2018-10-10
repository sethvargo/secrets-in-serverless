# Copyright 2018 Seth Vargo
# Copyright 2018 Google, Inc.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

import base64
import os
import googleapiclient.discovery

crypto_key_id = os.environ['KMS_CRYPTO_KEY_ID']

def decrypt(client, s):
    if not s:
        raise ValueError('ciphertext is missing')

    response = client \
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
