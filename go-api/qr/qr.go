// Package qr implements QR factorization of a real rectangular matrix using
// Householder reflections. Unlike classical/modified Gram-Schmidt, Householder
// reflections are numerically stable and work for any m x n matrix (m >= n,
// m == n or m < n), which is why they were chosen here.
//
// For an input A (m x n), Decompose returns Q (m x m, orthogonal) and
// R (m x n, upper triangular/trapezoidal) such that A = Q * R.
package qr

import (
	"math"

	"github.com/interseguro/matrices-go-api/matrix"
)

// Decompose computes the QR factorization of a using Householder reflections.
func Decompose(a matrix.Matrix) (q, r matrix.Matrix, err error) {
	if err := matrix.Validate(a); err != nil {
		return nil, nil, err
	}
	rows, cols := matrix.Dims(a)

	r = cloneMatrix(a)
	q = identity(rows)

	steps := rows
	if cols < steps {
		steps = cols
	}
	// If rows == 1 there is nothing to reflect (a single row is already
	// upper triangular), steps would be min(1, cols) = 1 but a 1xN vector
	// has no sub-diagonal entries to zero out.
	if rows == 1 {
		steps = 0
	}

	for k := 0; k < steps; k++ {
		length := rows - k
		x := make([]float64, length)
		for i := 0; i < length; i++ {
			x[i] = r[k+i][k]
		}

		normX := euclideanNorm(x)
		if normX == 0 {
			continue // column already zero below the diagonal, no reflection needed
		}

		alpha := normX
		if x[0] >= 0 {
			alpha = -normX
		}

		v := make([]float64, length)
		copy(v, x)
		v[0] -= alpha

		normV := euclideanNorm(v)
		if normV == 0 {
			continue // x was already a multiple of e1, no reflection needed
		}
		for i := range v {
			v[i] /= normV
		}

		applyHouseholderLeft(r, v, k, rows, cols)
		applyHouseholderRight(q, v, k, rows)
	}

	return q, r, nil
}

// applyHouseholderLeft updates r in place as r := H * r, where H = I - 2vv^T
// is embedded in the identity at rows/cols [k, rows).
func applyHouseholderLeft(r matrix.Matrix, v []float64, k, rows, cols int) {
	for j := k; j < cols; j++ {
		var dot float64
		for i := 0; i < len(v); i++ {
			dot += v[i] * r[k+i][j]
		}
		for i := 0; i < len(v); i++ {
			r[k+i][j] -= 2 * v[i] * dot
		}
	}
}

// applyHouseholderRight updates q in place as q := q * H, where H = I - 2vv^T
// is embedded in the identity at rows/cols [k, rows).
func applyHouseholderRight(q matrix.Matrix, v []float64, k, rows int) {
	for rIdx := 0; rIdx < rows; rIdx++ {
		var dot float64
		for i := 0; i < len(v); i++ {
			dot += v[i] * q[rIdx][k+i]
		}
		for i := 0; i < len(v); i++ {
			q[rIdx][k+i] -= 2 * v[i] * dot
		}
	}
}

func euclideanNorm(v []float64) float64 {
	var sum float64
	for _, x := range v {
		sum += x * x
	}
	return math.Sqrt(sum)
}

func identity(n int) matrix.Matrix {
	m := make(matrix.Matrix, n)
	for i := range m {
		m[i] = make([]float64, n)
		m[i][i] = 1
	}
	return m
}

func cloneMatrix(a matrix.Matrix) matrix.Matrix {
	out := make(matrix.Matrix, len(a))
	for i, row := range a {
		out[i] = append([]float64(nil), row...)
	}
	return out
}
