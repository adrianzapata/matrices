'use strict';

const jwt = require('jsonwebtoken');

/**
 * Express middleware requiring a valid "Bearer <token>" Authorization
 * header, signed with `secret` using HS256. Mirrors the Go API's JWT
 * middleware so both services enforce the same auth contract.
 *
 * @param {string} secret
 */
function jwtAuth(secret) {
  return (req, res, next) => {
    const header = req.get('Authorization') || '';
    if (!header.startsWith('Bearer ')) {
      return res.status(401).json({
        error: "missing or malformed Authorization header, expected 'Bearer <token>'",
      });
    }

    const token = header.slice('Bearer '.length);
    try {
      req.auth = jwt.verify(token, secret, { algorithms: ['HS256'] });
      return next();
    } catch (err) {
      return res.status(401).json({ error: 'invalid or expired token' });
    }
  };
}

module.exports = { jwtAuth };
