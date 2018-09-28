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
