'use strict';

const {Storage} = require('@google-cloud/storage');
const storage = new Storage();

let username;
let password;
storage.bucket(process.env.STORAGE_BUCKET).file('app1').download()
  .then(data => {
    console.log('here');
    const j = JSON.parse(data);
    username = j.username;
    password = j.password;
  })
  .catch(err => {
    console.error(err)
  });

exports.F = (req, res) => {
  res.send(`${username}:${password}`)
}
