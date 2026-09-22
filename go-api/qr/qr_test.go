package qr

import (
	"math"
	"testing"

	"github.com/interseguro/matrices-go-api/matrix"
)

const epsilon = 1e-9

func multiply(a, b matrix.Matrix) matrix.Matrix {
	rows, inner := len(a), len(a[0])
	cols := len(b[0])
	out := make(matrix.Matrix, rows)
	for i := 0; i < rows; i++ {
		out[i] = make([]float64, cols)
		for j := 0; j < cols; j++ {
			var sum float64
			for k := 0; k < inner; k++ {
				sum += a[i][k] * b[k][j]
			}
			out[i][j] = sum
		}
	}
	return out
}

func transpose(a matrix.Matrix) matrix.Matrix {
	rows, cols := len(a), len(a[0])
	out := make(matrix.Matrix, cols)
	for j := 0; j < cols; j++ {
		out[j] = make([]float64, rows)
		for i := 0; i < rows; i++ {
			out[j][i] = a[i][j]
		}
	}
	return out
}

func almostEqual(a, b matrix.Matrix, eps float64) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if len(a[i]) != len(b[i]) {
			return false
		}
		for j := range a[i] {
			if math.Abs(a[i][j]-b[i][j]) > eps {
				return false
			}
		}
	}
	return true
}

func assertReconstructs(t *testing.T, a matrix.Matrix) {
	t.Helper()
	q, r, err := Decompose(a)
	if err != nil {
		t.Fatalf("Decompose() error = %v", err)
	}

	// A = Q * R
	product := multiply(q, r)
	if !almostEqual(product, a, epsilon) {
		t.Fatalf("Q*R = %v, want %v", product, a)
	}

	// Q is orthogonal: Q^T * Q = I
	qtq := multiply(transpose(q), q)
	rows := len(q)
	for i := 0; i < rows; i++ {
		for j := 0; j < rows; j++ {
			want := 0.0
			if i == j {
				want = 1.0
			}
			if math.Abs(qtq[i][j]-want) > epsilon {
				t.Fatalf("Q not orthogonal: (Q^T Q)[%d][%d] = %v, want %v", i, j, qtq[i][j], want)
			}
		}
	}

	// R is upper triangular (zero strictly below the diagonal)
	for i := 0; i < len(r); i++ {
		for j := 0; j < i && j < len(r[i]); j++ {
			if math.Abs(r[i][j]) > epsilon {
				t.Fatalf("R not upper triangular: R[%d][%d] = %v", i, j, r[i][j])
			}
		}
	}
}

func TestDecompose_Square(t *testing.T) {
	assertReconstructs(t, matrix.Matrix{
		{12, -51, 4},
		{6, 167, -68},
		{-4, 24, -41},
	})
}

func TestDecompose_TallRectangular(t *testing.T) {
	assertReconstructs(t, matrix.Matrix{
		{1, 2},
		{3, 4},
		{5, 6},
	})
}

func TestDecompose_WideRectangular(t *testing.T) {
	assertReconstructs(t, matrix.Matrix{
		{1, 2, 3},
		{4, 5, 6},
	})
}

func TestDecompose_SingleRow(t *testing.T) {
	assertReconstructs(t, matrix.Matrix{
		{1, 2, 3},
	})
}

func TestDecompose_SingleColumn(t *testing.T) {
	assertReconstructs(t, matrix.Matrix{
		{3},
		{4},
	})
}

func TestDecompose_ZeroColumn(t *testing.T) {
	// First column already zero below the diagonal, exercises the
	// normX == 0 skip branch.
	assertReconstructs(t, matrix.Matrix{
		{0, 2},
		{0, 4},
	})
}

func TestDecompose_InvalidMatrix(t *testing.T) {
	_, _, err := Decompose(matrix.Matrix{{1, 2}, {3}})
	if err != matrix.ErrNotRectangular {
		t.Fatalf("Decompose() error = %v, want %v", err, matrix.ErrNotRectangular)
	}
}
