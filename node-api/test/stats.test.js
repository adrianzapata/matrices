'use strict';

const { test } = require('node:test');
const assert = require('node:assert/strict');

const { computeStats, isDiagonal } = require('../src/utils/stats');

test('isDiagonal returns true for a diagonal square matrix', () => {
  assert.equal(
    isDiagonal([
      [1, 0, 0],
      [0, 5, 0],
      [0, 0, -2],
    ]),
    true
  );
});

test('isDiagonal returns false when an off-diagonal entry is non-zero', () => {
  assert.equal(
    isDiagonal([
      [1, 0],
      [1, 5],
    ]),
    false
  );
});

test('isDiagonal returns false for non-square matrices', () => {
  assert.equal(isDiagonal([[1, 0, 0], [0, 1, 0]]), false);
});

test('isDiagonal treats the zero matrix as diagonal', () => {
  assert.equal(
    isDiagonal([
      [0, 0],
      [0, 0],
    ]),
    true
  );
});

test('computeStats aggregates max, min, sum and average across all matrices', () => {
  const result = computeStats({
    a: [
      [1, 2],
      [3, 4],
    ],
    b: [[10, -5]],
  });

  assert.equal(result.max, 10);
  assert.equal(result.min, -5);
  assert.equal(result.sum, 1 + 2 + 3 + 4 + 10 - 5);
  assert.equal(result.average, (1 + 2 + 3 + 4 + 10 - 5) / 6);
});

test('computeStats reports diagonal status per named matrix', () => {
  const result = computeStats({
    diag: [
      [2, 0],
      [0, 3],
    ],
    notDiag: [
      [2, 1],
      [0, 3],
    ],
  });

  assert.deepEqual(result.diagonal, { diag: true, notDiag: false });
});

test('computeStats throws on an empty matrices object', () => {
  assert.throws(() => computeStats({}));
});
