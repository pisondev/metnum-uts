package interpolation

import (
	"fmt"
	"math"
	"strings"
)

type NewtonStep struct {
	Order int
	Index int
	Value float64
}

func SolveNewtonInterpolation(points []Point, targets []float64) (string, error) {
	if err := ValidatePoints(points); err != nil {
		return "", err
	}

	coefficients, table, err := buildDividedDifferences(points)
	if err != nil {
		return "", err
	}

	var report strings.Builder
	report.WriteString("# Laporan 01. Newton Interpolation\n\n")
	WriteReportInputSection(&report, points, targets)
	report.WriteString("\n\n## Koefisien Divided Differences\n")
	for order := 0; order < len(table); order++ {
		for index := 0; index < len(table[order]); index++ {
			report.WriteString(fmt.Sprintf("- f[x_%d", index))
			for k := 1; k <= order; k++ {
				report.WriteString(fmt.Sprintf(", x_%d", index+k))
			}
			report.WriteString(fmt.Sprintf("] = %s\n", FormatReal(table[order][index])))
		}
	}

	report.WriteString("\n## Bentuk Polinom Newton\n")
	report.WriteString(fmt.Sprintf("P(x) = %s\n", buildNewtonPolynomialString(points, coefficients)))

	for i, target := range targets {
		value, products, err := evaluateNewton(points, coefficients, target)
		if err != nil {
			return "", err
		}

		report.WriteString(fmt.Sprintf("\n\n## Evaluasi Target %d\n", i+1))
		report.WriteString(fmt.Sprintf("- x = %s\n", FormatReal(target)))
		for k := 0; k < len(coefficients); k++ {
			if k == 0 {
				report.WriteString(fmt.Sprintf("- Suku 0: %s\n", FormatReal(coefficients[k])))
				continue
			}
			report.WriteString(fmt.Sprintf("- Suku %d: %s * %s = %s\n",
				k,
				FormatReal(coefficients[k]),
				products[k-1],
				FormatReal(coefficients[k]*productsValue(points, target, k)),
			))
		}
		report.WriteString(fmt.Sprintf("- Hasil akhir: P(%s) = %s\n", FormatReal(target), FormatReal(value)))
	}

	if len(targets) == 0 {
		report.WriteString("\n\n## Hasil Akhir\n- Polinom berhasil dibentuk, tetapi tidak ada target evaluasi.\n")
	}

	return report.String(), nil
}

func SolveLagrangeInterpolation(points []Point, targets []float64) (string, error) {
	if err := ValidatePoints(points); err != nil {
		return "", err
	}

	var report strings.Builder
	report.WriteString("# Laporan 02. Lagrange Interpolation\n\n")
	WriteReportInputSection(&report, points, targets)
	report.WriteString("\n\n## Basis Lagrange\n")
	for i := range points {
		report.WriteString(fmt.Sprintf("- L_%d(x) = %s\n", i, buildLagrangeBasisString(points, i)))
	}

	report.WriteString("\n## Bentuk Polinom Lagrange\n")
	report.WriteString(fmt.Sprintf("P(x) = %s\n", buildLagrangePolynomialString(points)))

	for i, target := range targets {
		value, basisValues, contributions, err := evaluateLagrange(points, target)
		if err != nil {
			return "", err
		}

		report.WriteString(fmt.Sprintf("\n\n## Evaluasi Target %d\n", i+1))
		report.WriteString(fmt.Sprintf("- x = %s\n", FormatReal(target)))
		for j := range points {
			report.WriteString(fmt.Sprintf("- y_%d * L_%d(x) = %s * %s = %s\n",
				j,
				j,
				FormatReal(points[j].Y),
				FormatReal(basisValues[j]),
				FormatReal(contributions[j]),
			))
		}
		report.WriteString(fmt.Sprintf("- Hasil akhir: P(%s) = %s\n", FormatReal(target), FormatReal(value)))
	}

	if len(targets) == 0 {
		report.WriteString("\n\n## Hasil Akhir\n- Polinom berhasil dibentuk, tetapi tidak ada target evaluasi.\n")
	}

	return report.String(), nil
}

func buildDividedDifferences(points []Point) ([]float64, [][]float64, error) {
	n := len(points)
	table := make([][]float64, n)
	table[0] = make([]float64, n)
	for i := range points {
		table[0][i] = points[i].Y
	}

	for order := 1; order < n; order++ {
		table[order] = make([]float64, n-order)
		for i := 0; i < n-order; i++ {
			denominator := points[i+order].X - points[i].X
			if math.Abs(denominator) < 1e-12 {
				return nil, nil, fmt.Errorf("selisih x menghasilkan pembagian nol pada divided differences")
			}
			table[order][i] = (table[order-1][i+1] - table[order-1][i]) / denominator
		}
	}

	coefficients := make([]float64, n)
	for i := 0; i < n; i++ {
		coefficients[i] = table[i][0]
	}

	return coefficients, table, nil
}

func evaluateNewton(points []Point, coefficients []float64, x float64) (float64, []string, error) {
	if len(points) != len(coefficients) {
		return 0, nil, fmt.Errorf("jumlah titik dan koefisien tidak konsisten")
	}

	value := coefficients[0]
	products := make([]string, 0, len(points)-1)
	product := 1.0
	for i := 1; i < len(coefficients); i++ {
		product *= x - points[i-1].X
		products = append(products, buildNewtonProductString(points, x, i))
		value += coefficients[i] * product
	}

	return value, products, nil
}

func productsValue(points []Point, x float64, terms int) float64 {
	product := 1.0
	for i := 0; i < terms; i++ {
		product *= x - points[i].X
	}
	return product
}

func evaluateLagrange(points []Point, x float64) (float64, []float64, []float64, error) {
	basisValues := make([]float64, len(points))
	contributions := make([]float64, len(points))
	value := 0.0

	for i := range points {
		basis := 1.0
		for j := range points {
			if i == j {
				continue
			}
			denominator := points[i].X - points[j].X
			if math.Abs(denominator) < 1e-12 {
				return 0, nil, nil, fmt.Errorf("selisih x menghasilkan pembagian nol pada basis lagrange")
			}
			basis *= (x - points[j].X) / denominator
		}

		basisValues[i] = basis
		contributions[i] = points[i].Y * basis
		value += contributions[i]
	}

	return value, basisValues, contributions, nil
}

func buildNewtonPolynomialString(points []Point, coefficients []float64) string {
	parts := make([]string, 0, len(coefficients))
	for i := range coefficients {
		if i == 0 {
			parts = append(parts, FormatReal(coefficients[i]))
			continue
		}
		parts = append(parts, fmt.Sprintf("(%s)(%s)", FormatReal(coefficients[i]), buildNewtonFactorString(points, i)))
	}
	return strings.Join(parts, " + ")
}

func buildNewtonFactorString(points []Point, order int) string {
	factors := make([]string, 0, order)
	for i := 0; i < order; i++ {
		factors = append(factors, fmt.Sprintf("(x-%s)", FormatReal(points[i].X)))
	}
	return strings.Join(factors, "")
}

func buildNewtonProductString(points []Point, x float64, order int) string {
	parts := make([]string, 0, order)
	for i := 0; i < order; i++ {
		parts = append(parts, fmt.Sprintf("(%s-%s)", FormatReal(x), FormatReal(points[i].X)))
	}
	return strings.Join(parts, "")
}

func buildLagrangeBasisString(points []Point, basisIndex int) string {
	parts := make([]string, 0, len(points)-1)
	for j := range points {
		if j == basisIndex {
			continue
		}
		parts = append(parts, fmt.Sprintf("((x-%s)/(%s-%s))", FormatReal(points[j].X), FormatReal(points[basisIndex].X), FormatReal(points[j].X)))
	}
	return strings.Join(parts, "")
}

func buildLagrangePolynomialString(points []Point) string {
	parts := make([]string, 0, len(points))
	for i := range points {
		parts = append(parts, fmt.Sprintf("(%s)L_%d(x)", FormatReal(points[i].Y), i))
	}
	return strings.Join(parts, " + ")
}
