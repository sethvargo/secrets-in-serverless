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

import os
import requests

vault_addr = os.environ['VAULT_ADDR']

jwt = requests.get('http://metadata/computeMetadata/v1/instance/service-accounts/default/identity',
    headers={'Metadata-Flavor':'Google'},
    params={'audience':'http://vault/socialmedia', 'format':'full'})

auth = requests.post(vault_addr + '/v1/auth/gcp/login',
    json={'role':'socialmedia', 'jwt':jwt.text})
token = auth.json()['auth']['client_token']

r = requests.get(vault_addr + '/v1/secret/apikeys/twitter',
    headers={'x-vault-token': token})

apikey = r.json()['data']['value']

def F(request):
    return f'{apikey}'
