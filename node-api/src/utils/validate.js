'use strict';

/**
 * @param {unknown} matrix
 * @returns {string | null} an error message, or null if the matrix is a
 *   well-formed, non-empty, rectangular array of arrays of finite numbers.
 */
function validateMatrix(matrix) {
  if (!Array.isArray(matrix) || matrix.length === 0) {
    return 'matrix must be a non-empty array of arrays';
  }
  if (!Array.isArray(matrix[0]) || matrix[0].length === 0) {
    return 'matrix must be a non-empty array of arrays';
  }

  const cols = matrix[0].length;
  for (const row of matrix) {
    if (!Array.isArray(row) || row.length !== cols) {
      return 'matrix rows must all have the same length';
    }
    for (const value of row) {
      if (typeof value !== 'number' || !Number.isFinite(value)) {
        return 'matrix values must all be finite numbers';
      }
    }
  }

  return null;
}

/**
 * @param {unknown} matrices
 * @returns {string | null} an error message, or null if matrices is a
 *   non-empty object mapping names to valid matrices.
 */
function validateMatrices(matrices) {
  if (typeof matrices !== 'object' || matrices === null || Array.isArray(matrices)) {
    return 'matrices must be an object mapping names to matrices';
  }
  const names = Object.keys(matrices);
  if (names.length === 0) {
    return 'matrices must contain at least one matrix';
  }
  for (const name of names) {
    const err = validateMatrix(matrices[name]);
    if (err) return `matrices.${name}: ${err}`;
  }
  return null;
}

module.exports = { validateMatrix, validateMatrices };
