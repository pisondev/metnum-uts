package rootfinding

import "fmt"

type Problem struct {
	FunctionRaw    string
	Function       Expression
	DerivativeRaw  string
	Derivative     *Expression
	Lower          float64
	Upper          float64
	X0             float64
	X1             float64
	Tolerance      float64
	MaxIterations  int
}

type RootMethod struct {
	Choice      int
	Code        string
	Name        string
	DisplayName string
	ResultsSlug string
	Solve       func(Problem) (string, error)
}

func GetRootMethod(choice int) (RootMethod, error) {
	switch choice {
	case 1:
		return RootMethod{
			Choice:      1,
			Code:        "01",
			Name:        "Bisection Method",
			DisplayName: "01. Bisection Method",
			ResultsSlug: "root-finding/01-bisection-method",
			Solve:       SolveBisection,
		}, nil
	case 2:
		return RootMethod{
			Choice:      2,
			Code:        "02",
			Name:        "Secant Method",
			DisplayName: "02. Secant Method",
			ResultsSlug: "root-finding/02-secant-method",
			Solve:       SolveSecant,
		}, nil
	case 3:
		return RootMethod{
			Choice:      3,
			Code:        "03",
			Name:        "Newton Method",
			DisplayName: "03. Newton Method",
			ResultsSlug: "root-finding/03-newton-method",
			Solve:       SolveNewton,
		}, nil
	default:
		return RootMethod{}, fmt.Errorf("pilihan metode harus 1, 2, atau 3")
	}
}
