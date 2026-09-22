'use strict';

const express = require('express');
const { computeStats } = require('../utils/stats');
const { validateMatrices } = require('../utils/validate');

const router = express.Router();

/**
 * POST /api/stats
 * Body: { "matrices": { "<name>": number[][], ... } }
 *
 * Receives the matrices produced by the Go API (rotated matrix, Q and R
 * from the QR factorization) and returns aggregate statistics over them,
 * per the "operación adicional" section of the challenge.
 */
router.post('/stats', (req, res) => {
  const { matrices } = req.body || {};

  const validationError = validateMatrices(matrices);
  if (validationError) {
    return res.status(422).json({ error: validationError });
  }

  const stats = computeStats(matrices);
  return res.json(stats);
});

module.exports = router;
