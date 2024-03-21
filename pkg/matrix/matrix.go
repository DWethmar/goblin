// Package matrix provides a simple 2D matrix implementation.
package matrix

type Matrix [][]int

// Get returns the value at the given position.
// If the position is out of bounds, the default value is returned.
func (m Matrix) Get(x, y, d int) int {
	if m.Contains(x, y) {
		return m[y][x]
	}
	return d
}

// Set sets the value at the given position.
func (m Matrix) Set(x, y int, v int) {
	if m.Contains(x, y) {
		m[y][x] = v
	}
}

func (m Matrix) Contains(x, y int) bool {
	return y >= 0 && y < len(m) && x >= 0 && x < len(m[y])
}

// Iterate iterates over all values, from top to bottom, left to right.
func Iterate(m Matrix, fn func(x, y, v int)) {
	for y := range m {
		for x, v := range m[y] {
			fn(x, x, v)
		}
	}
}

func New(width, height, v int) Matrix {
	m := make([][]int, height)
	for i := range m {
		m[i] = make([]int, width)
		for j := range m[i] {
			m[i][j] = v
		}
	}
	return m
}
