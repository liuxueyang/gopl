package eval

import (
	"fmt"
	"math"
	"testing"
)

type Expr interface {
	Eval(env Env) float64
}

type Var string

func (v Var) Eval(env Env) float64 {
	return env[v]
}

type literal float64

func (v literal) Eval(_ Env) float64 {
	return float64(v)
}

type unary struct {
	op rune
	x  Expr
}

func (v unary) Eval(env Env) float64 {
	switch v.op {
	case '+':
		return +v.x.Eval(env)
	case '-':
		return -v.x.Eval(env)
	}
	panic(fmt.Sprintf("unsupported unary operator: %q", v.op))
}

type binary struct {
	op   rune
	x, y Expr
}

func (v binary) Eval(env Env) float64 {
	switch v.op {
	case '+':
		return v.x.Eval(env) + v.y.Eval(env)
	case '-':
		return v.x.Eval(env) - v.y.Eval(env)
	case '*':
		return v.x.Eval(env) * v.y.Eval(env)
	case '/':
		return v.x.Eval(env) / v.y.Eval(env)
	}
	panic(fmt.Sprintf("unsupported binary operator: %q", v.op))
}

type call struct {
	fn   string
	args []Expr
}

func (v call) Eval(env Env) float64 {
	switch v.fn {
	case "pow":
		return math.Pow(v.args[0].Eval(env), v.args[1].Eval(env))
	case "sin":
		return math.Sin(v.args[0].Eval(env))
	case "sqrt":
		return math.Sqrt(v.args[0].Eval(env))
	}
	panic(fmt.Sprintf("unsupported function call: %s", v.fn))
}

type Env map[Var]float64

func TestEval(t *testing.T) {

	tests := []struct {
		expr string
		env  Env
		want string
	}{
		{"sqrt(A/pi)", Env{"A": 87616, "pi": math.Pi}, "167"},
		{"pow(x, 3) + pow(y, 3)", Env{"x": 12, "y": 1}, "1729"},
		{"pow(x, 3) + pow(y, 3)", Env{"x": 9, "y": 10}, "1729"},
		{"5/9*(F-32)", Env{"F": -40}, "-40"},
		{"5/9*(F-32)", Env{"F": 32}, "0"},
		{"5/9*(F-32)", Env{"F": 212}, "100"},
	}

	var preExpr string
	for _, test := range tests {
		if test.expr != preExpr {
			fmt.Printf("\n%s\n", test.expr)
			preExpr = test.expr
		}
	}
}
