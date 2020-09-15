package eval

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEvalWithBasicExpr(t *testing.T) {
	var tests = []struct {
		name     string
		expected float64
		given    string
	}{
		{"trivial single number",
			3.14, "3.14"},
		{"simple sum",
			5, "2+3"},
		{"simple subtraction",
			-2, "5-7"},
		{"simple multiplication",
			23000, "1000*23"},
		{"simple division",
			10, "230/23"},
		{"division and multiplication",
			15, "35/7*3"},
		{"division and multiplication with sum and sub",
			47, "35/7*3+33-1"},
		{"with parenthesis",
			1070, "107*(23-13)"},
		{"with 2 parenthesis",
			-6, "(10-13)*(15-13)"},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			actual, err := Eval(tt.given)
			assert.Nil(t, err)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
