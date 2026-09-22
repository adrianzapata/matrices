'use strict';

const { test } = require('node:test');
const assert = require('node:assert/strict');

const { validateMatrix, validateMatrices } = require('../src/utils/validate');

test('validateMatrix accepts a well-formed rectangular matrix', () => {
  assert.equal(
    validateMatrix([
      [1, 2],
      [3, 4],
    ]),
    null
  );
});

test('validateMatrix rejects an empty array', () => {
  assert.match(validateMatrix([]), /non-empty/);
});

test('validateMatrix rejects ragged rows', () => {
  assert.match(validateMatrix([[1, 2], [3]]), /same length/);
});

test('validateMatrix rejects non-numeric values', () => {
  assert.match(validateMatrix([[1, 'x']]), /finite numbers/);
});

test('validateMatrix rejects non-finite values', () => {
  assert.match(validateMatrix([[1, Infinity]]), /finite numbers/);
});

test('validateMatrices requires an object of matrices', () => {
  assert.match(validateMatrices([1, 2, 3]), /object mapping names/);
});

test('validateMatrices requires at least one matrix', () => {
  assert.match(validateMatrices({}), /at least one matrix/);
});

test('validateMatrices prefixes errors with the matrix name', () => {
  assert.match(validateMatrices({ q: [[1]], bad: [[1, 2], [3]] }), /^matrices\.bad:/);
});
