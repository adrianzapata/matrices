'use strict';

const { test, before, after } = require('node:test');
const assert = require('node:assert/strict');
const jwt = require('jsonwebtoken');

process.env.JWT_SECRET = 'test-secret';

const { createApp } = require('../src/index');

let server;
let baseUrl;

before(async () => {
  const app = createApp();
  server = app.listen(0);
  await new Promise((resolve) => server.once('listening', resolve));
  const { port } = server.address();
  baseUrl = `http://127.0.0.1:${port}`;
});

after(() => {
  server.close();
});

test('GET /health returns ok without auth', async () => {
  const res = await fetch(`${baseUrl}/health`);
  assert.equal(res.status, 200);
  assert.deepEqual(await res.json(), { status: 'ok' });
});

test('POST /api/stats without a token is rejected', async () => {
  const res = await fetch(`${baseUrl}/api/stats`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ matrices: { a: [[1]] } }),
  });
  assert.equal(res.status, 401);
});

test('POST /api/stats with a valid token returns computed statistics', async () => {
  const token = jwt.sign({ sub: 'test' }, process.env.JWT_SECRET);
  const res = await fetch(`${baseUrl}/api/stats`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify({
      matrices: {
        rotated: [
          [4, 1],
          [5, 2],
        ],
        r: [
          [2, 0],
          [0, 3],
        ],
      },
    }),
  });

  assert.equal(res.status, 200);
  const body = await res.json();
  assert.equal(body.max, 5);
  assert.equal(body.min, 0);
  assert.equal(body.sum, 4 + 1 + 5 + 2 + 2 + 0 + 0 + 3);
  assert.deepEqual(body.diagonal, { rotated: false, r: true });
});

test('POST /api/stats with an invalid matrix payload returns 422', async () => {
  const token = jwt.sign({ sub: 'test' }, process.env.JWT_SECRET);
  const res = await fetch(`${baseUrl}/api/stats`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify({ matrices: { bad: [[1, 2], [3]] } }),
  });
  assert.equal(res.status, 422);
});
