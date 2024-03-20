package matrix

import (
	"reflect"
	"testing"
)

func TestMatrix_Get(t *testing.T) {
	type args struct {
		x int
		y int
		d int
	}
	tests := []struct {
		name string
		m    Matrix
		args args
		want int
	}{
		{
			name: "should return default value",
			m: Matrix(
				[][]int{
					{1, 2, 3},
					{4, 5, 6},
					{7, 8, 9},
				},
			),
			args: args{
				x: 1,
				y: 1,
				d: -1,
			},
			want: 5,
		},
		{
			name: "should return default value",
			m: Matrix(
				[][]int{
					{1, 2, 3},
					{4, 5, 6},
					{7, 8, 9},
				},
			),
			args: args{
				x: 3,
				y: 3,
				d: -1,
			},
			want: -1,
		},
		{
			name: "x is out of bounds",
			m: Matrix(
				[][]int{
					{1, 2, 3},
					{4, 5, 6},
					{7, 8, 9},
				},
			),
			args: args{
				x: -1,
				y: 0,
				d: -1,
			},
			want: -1,
		},
		{
			name: "y is out of bounds",
			m: Matrix(
				[][]int{
					{1, 2, 3},
					{4, 5, 6},
					{7, 8, 9},
				},
			),
			args: args{
				x: 0,
				y: -1,
				d: -1,
			},
			want: -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.Get(tt.args.x, tt.args.y, tt.args.d); got != tt.want {
				t.Errorf("Matrix.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMatrix_Set(t *testing.T) {
	t.Run("should set value", func(t *testing.T) {
		m := Matrix(
			[][]int{
				{1, 2, 3},
				{4, 5, 6},
				{7, 8, 9},
			},
		)
		m.Set(1, 1, 10)
		if m[1][1] != 10 {
			t.Errorf("Matrix.Set() = %v, want %v", m[1][1], 10)
		}
	})

	t.Run("should ignore out of bounds", func(t *testing.T) {
		m := Matrix([][]int{
			{1, 2, 3},
			{4, 5, 6},
			{7, 8, 9},
		})
		m.Set(3, 3, 10)
	})
}

func TestMatrix_Contains(t *testing.T) {
	type args struct {
		x int
		y int
	}
	tests := []struct {
		name string
		m    Matrix
		args args
		want bool
	}{
		{
			name: "should return true",
			m: Matrix([][]int{
				{1, 2, 3},
				{4, 5, 6},
				{7, 8, 9},
			}),
			args: args{
				x: 1,
				y: 1,
			},
			want: true,
		},
		{
			name: "should return false for x out of bounds",
			m: Matrix([][]int{
				{1, 2, 3},
				{4, 5, 6},
				{7, 8, 9},
			}),
			args: args{
				x: 3,
				y: 1,
			},
			want: false,
		},
		{
			name: "should return false for y out of bounds",
			m: Matrix([][]int{
				{1, 2, 3},
				{4, 5, 6},
				{7, 8, 9},
			}),
			args: args{
				x: 1,
				y: 3,
			},
			want: false,
		},
		{
			name: "should return false for x and y out of bounds",
			m: Matrix([][]int{
				{1, 2, 3},
				{4, 5, 6},
				{7, 8, 9},
			}),
			args: args{
				x: 3,
				y: 3,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.Contains(tt.args.x, tt.args.y); got != tt.want {
				t.Errorf("Matrix.Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIterate(t *testing.T) {
	t.Run("should iterate over all values", func(t *testing.T) {
		m := Matrix([][]int{
			{1, 2, 3},
			{4, 5, 6},
			{7, 8, 9},
		})
		var values []int
		Iterate(m, func(x, y, v int) {
			values = append(values, v)
		})
		if !reflect.DeepEqual(values, []int{1, 2, 3, 4, 5, 6, 7, 8, 9}) {
			t.Errorf("Iterate() = %v, want %v", values, []int{1, 2, 3, 4, 5, 6, 7, 8, 9})
		}
	})
}

func TestNew(t *testing.T) {
	t.Run("should create matrix", func(t *testing.T) {
		m := New(3, 3, 0)
		e := Matrix([][]int{
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
		})
		if !reflect.DeepEqual(m, e) {
			t.Errorf("New() = %v, want %v", m, e)
		}
	})
}
