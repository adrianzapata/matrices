'use strict';

const express = require('express');
const cors = require('cors');

const { jwtAuth } = require('./middleware/auth');
const statsRouter = require('./routes/stats');

const PORT = process.env.PORT || 4000;
const JWT_SECRET = process.env.JWT_SECRET || 'dev-secret-change-me';
const CORS_ORIGINS = process.env.CORS_ORIGINS || '*';

function createApp() {
  const app = express();

  app.use(cors({ origin: CORS_ORIGINS }));
  app.use(express.json({ limit: '5mb' }));

  app.get('/health', (req, res) => res.json({ status: 'ok' }));

  app.use('/api', jwtAuth(JWT_SECRET), statsRouter);

  // Centralized error handler, e.g. for malformed JSON bodies raised by
  // express.json().
  app.use((err, req, res, next) => {
    if (err) {
      return res.status(400).json({ error: 'invalid JSON body' });
    }
    return next();
  });

  return app;
}

if (require.main === module) {
  const app = createApp();
  app.listen(PORT, () => {
    console.log(`matrices-node-api listening on :${PORT}`);
  });
}

module.exports = { createApp };
