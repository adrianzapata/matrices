package matrix

import (
	"reflect"
	"testing"
)

func TestValidate(t *testing.T) {
	cases := []struct {
		name    string
		m       Matrix
		wantErr error
	}{
		{"empty rows", Matrix{}, ErrEmptyMatrix},
		{"empty first row", Matrix{{}}, ErrEmptyMatrix},
		{"ragged", Matrix{{1, 2}, {1}}, ErrNotRectangular},
		{"ok square", Matrix{{1, 2}, {3, 4}}, nil},
		{"ok rectangular", Matrix{{1, 2, 3}, {4, 5, 6}}, nil},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if err := Validate(c.m); err != c.wantErr {
				t.Fatalf("Validate() = %v, want %v", err, c.wantErr)
			}
		})
	}
}

func TestRotateClockwise90_Square(t *testing.T) {
	in := Matrix{
		{1, 2},
		{3, 4},
	}
	want := Matrix{
		{3, 1},
		{4, 2},
	}
	got := RotateClockwise90(in)
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("RotateClockwise90() = %v, want %v", got, want)
	}
}

func TestRotateClockwise90_Rectangular(t *testing.T) {
	// 2x3 -> 3x2
	in := Matrix{
		{1, 2, 3},
		{4, 5, 6},
	}
	want := Matrix{
		{4, 1},
		{5, 2},
		{6, 3},
	}
	got := RotateClockwise90(in)
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("RotateClockwise90() = %v, want %v", got, want)
	}
}

func TestRotateClockwise90_Idempotent4x(t *testing.T) {
	in := Matrix{
		{1, 2, 3},
		{4, 5, 6},
	}
	got := in
	for i := 0; i < 4; i++ {
		got = RotateClockwise90(got)
	}
	if !reflect.DeepEqual(got, in) {
		t.Fatalf("rotating 4 times should return to original: got %v, want %v", got, in)
	}
}
