// Shader is used to represent and work with GL Shaders in Go.
// Provides things like Uniforms, Shader models,
package main

import (
	"github.com/go-gl/gl/v2.1/gl"
	"gonum.org/v1/gonum/mat"
)

//Program represents a GL Shader Program
type Program struct {

	// The GL ID for this program
	Id uint32
	// The Uniforms for this program
	Uniforms map[string]*Uniform

	// The shaders linked to this program
	Shaders []*Shader
}

func (*Program) Attach(s *Shader) {

}

func (p *Program) Link() {
	gl.LinkProgram(p.Id)
}

//NewProgram creates a new shader program
func NewProgram() *Program {
	return &Program{
		Id:       gl.CreateProgram(),
		Uniforms: map[string]*Uniform{},
		Shaders:  []*Shader{},
	}
}

type Shader struct {
}

func (*Shader) Compile() {
}

type Uniform struct {
	Location int32
}

func (u *Uniform) Int(i int32) {
	gl.Uniform1i(u.Location, i)
}

func (u *Uniform) Float(f float32) {
	gl.Uniform1f(u.Location, f)
}

func (u *Uniform) Double(f float64) {
	gl.Uniform1d(u.Location, f)
}

//Matrix puts a matrix as an uniform for a given shader
func (u *Uniform) Matrix(m mat.Matrix, transpose bool) {
	x, y := m.Dims()

	matrix := make([]float64, x*y)

	// store the matrix in the array
	MatrixToArray(m, matrix)

	if x == y {
		switch x {
		case 2:
			gl.UniformMatrix2dv(u.Location, 4, transpose, &matrix[0])
		case 3:
			gl.UniformMatrix3dv(u.Location, 9, transpose, &matrix[0])
		case 4:
			gl.UniformMatrix4dv(u.Location, 16, transpose, &matrix[0])
		}
	}
}
