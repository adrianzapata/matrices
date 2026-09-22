// Package matrix provides basic operations on rectangular float64 matrices:
// validation and 90-degree rotation.
package matrix

import "errors"

// Matrix is a rectangular grid of float64 values stored row-major.
type Matrix [][]float64

var (
	ErrEmptyMatrix    = errors.New("matrix must contain at least one row and one column")
	ErrNotRectangular = errors.New("matrix rows must all have the same length")
)

// Validate checks that m is non-empty and rectangular (every row has the
// same number of columns as the first row).
func Validate(m Matrix) error {
	if len(m) == 0 || len(m[0]) == 0 {
		return ErrEmptyMatrix
	}
	cols := len(m[0])
	for _, row := range m {
		if len(row) != cols {
			return ErrNotRectangular
		}
	}
	return nil
}

// Dims returns the number of rows and columns of m. Callers must ensure m
// has already passed Validate.
func Dims(m Matrix) (rows, cols int) {
	return len(m), len(m[0])
}

// RotateClockwise90 returns a new matrix equal to m rotated 90 degrees
// clockwise. For an r x c input the output is c x r, where
// out[j][r-1-i] = m[i][j].
func RotateClockwise90(m Matrix) Matrix {
	rows, cols := Dims(m)
	out := make(Matrix, cols)
	for j := 0; j < cols; j++ {
		out[j] = make([]float64, rows)
		for i := 0; i < rows; i++ {
			out[j][rows-1-i] = m[i][j]
		}
	}
	return out
}
