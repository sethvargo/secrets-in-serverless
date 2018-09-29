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

const cryptoKeyID = process.env.KMS_CRYPTO_KEY_ID;

const kms = require('@google-cloud/kms');
const client = new kms.v1.KeyManagementServiceClient();

let username;
client.decrypt({
  name: cryptoKeyID,
  ciphertext: process.env.DB_USER,
}).then(res => {
  username = res[0].plaintext.toString().trim();
}).catch(err => {
  console.error(err);
});

let password;
client.decrypt({
  name: cryptoKeyID,
  ciphertext: process.env.DB_PASS,
}).then(res => {
  password = res[0].plaintext.toString().trim();
}).catch(err => {
  console.error(err);
});

exports.F = (req, res) => {
  res.send(`${username}:${password}`)
}
