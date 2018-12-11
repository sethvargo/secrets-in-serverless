// Copyright 2018 Seth Vargo
// Copyright 2018 Google, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

'use strict';

const fetch = require('node-fetch');

const vaultAddr = process.env.VAULT_ADDR;

let jwt;
const getJwt = async () => {
  if (!jwt) {
    const url = 'http://metadata/computeMetadata/v1/instance/service-accounts/default/identity'
    const resp = await fetch(url + '?audience=http://vault/socialmedia&format=full', {
      headers: { 'Metadata-Flavor': 'Google' },
    });
    jwt = await resp.text();
  }
  return jwt;
}

let token;
const getToken = async (jwt) => {
  if (!token) {
    const resp = await fetch(vaultAddr + '/v1/auth/gcp/login', {
      method: 'post',
      body: JSON.stringify({
        'role': 'socialmedia',
        'jwt': jwt,
      }),
    });
    const j = await resp.json();
    token = j['auth']['client_token'];
  }
  return token;
}

let apikey;
const getApikey = async (token) => {
  if (!apikey) {
    const resp = await fetch(vaultAddr + '/v1/secret/apikeys/twitter', {
      headers: { 'x-vault-token': token },
    });
    const j = await resp.json();
    apikey = j['data']['value'];
  }
  return apikey;
}

exports.F = async (req, res) => {
  const apikey = await getApikey(await getToken(await getJwt()))
  res.send(`${apikey}`)
}
