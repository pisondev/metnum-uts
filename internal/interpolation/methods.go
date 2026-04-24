package interpolation

import "fmt"

type InterpolationMethod struct {
	Choice      int
	Code        string
	Name        string
	DisplayName string
	ResultsSlug string
	Solve       func([]Point, []float64) (string, error)
}

func GetInterpolationMethod(choice int) (InterpolationMethod, error) {
	switch choice {
	case 1:
		return InterpolationMethod{
			Choice:      1,
			Code:        "01",
			Name:        "Newton Interpolation",
			DisplayName: "01. Newton Interpolation",
			ResultsSlug: "polynomial-interpolation/01-newton-interpolation",
			Solve:       SolveNewtonInterpolation,
		}, nil
	case 2:
		return InterpolationMethod{
			Choice:      2,
			Code:        "02",
			Name:        "Lagrange Interpolation",
			DisplayName: "02. Lagrange Interpolation",
			ResultsSlug: "polynomial-interpolation/02-lagrange-interpolation",
			Solve:       SolveLagrangeInterpolation,
		}, nil
	default:
		return InterpolationMethod{}, fmt.Errorf("pilihan metode harus 1 atau 2")
	}
}
