package spline

import "fmt"

type SplineMethod struct {
	Choice      int
	Code        string
	Name        string
	DisplayName string
	ResultsSlug string
	Solve       func([]Point, []float64) (string, error)
}

func GetSplineMethod(choice int) (SplineMethod, error) {
	switch choice {
	case 1:
		return SplineMethod{
			Choice:      1,
			Code:        "01",
			Name:        "Linear Spline",
			DisplayName: "01. Linear Spline",
			ResultsSlug: "01-linear-spline",
			Solve:       SolveLinearSpline,
		}, nil
	case 2:
		return SplineMethod{
			Choice:      2,
			Code:        "02",
			Name:        "Quadratic Spline",
			DisplayName: "02. Quadratic Spline",
			ResultsSlug: "02-quadratic-spline",
			Solve:       SolveQuadraticSpline,
		}, nil
	case 3:
		return SplineMethod{
			Choice:      3,
			Code:        "03",
			Name:        "Natural Cubic Spline",
			DisplayName: "03. Natural Cubic Spline",
			ResultsSlug: "03-natural-cubic-spline",
			Solve:       SolveCubicSpline,
		}, nil
	default:
		return SplineMethod{}, fmt.Errorf("pilihan metode harus 1, 2, atau 3")
	}
}
