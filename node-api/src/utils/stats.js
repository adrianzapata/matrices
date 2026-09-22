'use strict';

/**
 * Pure statistics helpers over one or more numeric matrices, per the
 * "operación adicional" section of the challenge: max, min, average, sum
 * across all values, and a per-matrix diagonal check.
 */

/**
 * @param {number[][]} matrix
 * @returns {boolean} true if matrix is square and every off-diagonal entry
 *   is zero. A non-square matrix can never be diagonal, so it returns false.
 */
function isDiagonal(matrix) {
  const rows = matrix.length;
  if (rows === 0) return false;
  const cols = matrix[0].length;
  if (rows !== cols) return false;

  for (let i = 0; i < rows; i++) {
    for (let j = 0; j < cols; j++) {
      if (i !== j && matrix[i][j] !== 0) {
        return false;
      }
    }
  }
  return true;
}

/**
 * @param {Record<string, number[][]>} matrices named matrices to aggregate
 * @returns {{max: number, min: number, average: number, sum: number, diagonal: Record<string, boolean>}}
 */
function computeStats(matrices) {
  const names = Object.keys(matrices);
  if (names.length === 0) {
    throw new Error('at least one matrix is required');
  }

  let max = -Infinity;
  let min = Infinity;
  let sum = 0;
  let count = 0;
  const diagonal = {};

  for (const name of names) {
    const matrix = matrices[name];
    for (const row of matrix) {
      for (const value of row) {
        if (value > max) max = value;
        if (value < min) min = value;
        sum += value;
        count++;
      }
    }
    diagonal[name] = isDiagonal(matrix);
  }

  return {
    max,
    min,
    average: count === 0 ? 0 : sum / count,
    sum,
    diagonal,
  };
}

module.exports = { computeStats, isDiagonal };
