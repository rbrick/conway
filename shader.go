// Shader is used to represent and work with GL Shaders in Go.
// Provides things like Uniforms, Shader models,
package main

import (
	"github.com/go-gl/gl/v2.1/gl"
	"gonum.org/v1/gonum/mat"
)

type Uniform struct {
	Location int32
}

func (u *Uniform) Matrix4f(m mat.Matrix, transpose bool) {
	x, y := m.Dims()

	matrix := make([]float64, x*y)

	// store the matrix in the array
	matrixToArray(m, matrix)

	if x == y {
		switch x {
		case 2:
			gl.UniformMatrix2dv(u.Location, 4, transpose, &matrix[0])
		case 3:
			gl.UniformMatrix3dv(u.Location, 9, transpose, &matrix[0])
		case 4:
			gl.UniformMatrix4dv(u.Location, 16, transpose, &matrix[0])
		}
	} else {
		// TODO: support for non square matrices
	}
}

func matrixToArray(m mat.Matrix, array []float64) {
	r, c := m.Dims()

	for i := 0; i < r; i++ {
		row := mat.Row(nil, i, m)
		for j := 0; j < c; j++ {
			array[i+j*c] = row[j]
		}
	}
}
