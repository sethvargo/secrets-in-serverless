'use strict';

const username = process.env.DB_USER;
const password = process.env.DB_PASS;

exports.F = (req, res) => {
  res.send(`${username}:${password}`);
};
