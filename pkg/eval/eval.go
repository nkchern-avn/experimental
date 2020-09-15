package eval

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strconv"
)

var (
	binOps = map[token.Token]func(float64, float64) float64{
		token.ADD: func(x float64, y float64) float64 {
			return x + y
		},
		token.SUB: func(x float64, y float64) float64 {
			return x - y
		},
		token.MUL: func(x float64, y float64) float64 {
			return x * y
		},

		token.QUO: func(x float64, y float64) float64 {
			return x / y
		},
	}
)

func Eval(expr string) (float64, error) {
	parsed, err := parser.ParseExpr(expr)
	if err != nil {
		return 0, err
	}
	//	ast.Print(token.NewFileSet(), parsed)
	return evaluate(parsed)
}

func applyOp(op token.Token, x float64, y float64) (float64, error) {
	fn, found := binOps[op]
	if !found {
		return 0, fmt.Errorf("Unsupported operation: %s", op)
	}
	return fn(x, y), nil
}

func evaluate(expr ast.Expr) (float64, error) {
	switch ex := expr.(type) {
	case *ast.BasicLit:
		return strconv.ParseFloat(ex.Value, 64)
	case *ast.BinaryExpr:
		x, err := evaluate(ex.X)
		if err != nil {
			return 0, err
		}

		y, err := evaluate(ex.Y)
		if err != nil {
			return 0, err
		}

		return applyOp(ex.Op, x, y)
	case *ast.ParenExpr:
		return evaluate(ex.X)
	default:
		return 0, fmt.Errorf("Unsupported expression: %T", expr)
	}
}
