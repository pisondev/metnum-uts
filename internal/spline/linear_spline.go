package spline

import (
	"fmt"
	"strings"
)

func SolveLinearSpline(points []Point, targets []float64) (string, error) {
	if err := ValidatePoints(points); err != nil {
		return "", err
	}

	var report strings.Builder
	n := len(points) - 1
	h := make([]float64, n)
	m := make([]float64, n)

	report.WriteString("# Laporan 01. Linear Spline Interpolation\n\n")
	WriteReportInputSection(&report, points, targets)

	PrintStageTitle("Tahap 1: Gradien Tiap Segmen")
	report.WriteString("\n## Tahap 1: Gradien Tiap Segmen\n")
	report.WriteString("Rumus umum: $m_i = \\frac{y_{i+1} - y_i}{x_{i+1} - x_i}$\n\n")

	for i := 0; i < n; i++ {
		h[i] = points[i+1].X - points[i].X
		m[i] = (points[i+1].Y - points[i].Y) / h[i]
		line := fmt.Sprintf(
			"- Segmen %d: $m_%d = \\frac{%s - %s}{%s - %s} = %s$\n",
			i,
			i,
			FormatNum(points[i+1].Y),
			FormatNum(points[i].Y),
			FormatNum(points[i+1].X),
			FormatNum(points[i].X),
			FormatNum(m[i]),
		)
		fmt.Print(line)
		report.WriteString(line)
	}

	PrintStageTitle("Tahap 2: Persamaan Setiap Segmen")
	report.WriteString("\n## Tahap 2: Persamaan Setiap Segmen\n")
	report.WriteString("Model: $S_i(x) = a_i + b_i(x - x_i)$ dengan $a_i = y_i$ dan $b_i = m_i$.\n\n")

	for i := 0; i < n; i++ {
		a := points[i].Y
		b := m[i]
		line := fmt.Sprintf(
			"- Segmen %d pada interval $[%s, %s]$: $S_%d(x) = %s + (%s)(x - %s)$\n",
			i,
			FormatNum(points[i].X),
			FormatNum(points[i+1].X),
			i,
			FormatFractionOnly(a),
			FormatFractionOnly(b),
			FormatFractionOnly(points[i].X),
		)
		fmt.Print(line)
		report.WriteString(line)
	}

	for _, targetX := range targets {
		fmt.Printf("\n===================================================\n")
		fmt.Printf("EVALUASI UNTUK X DITANYA = %v\n", FormatNum(targetX))
		report.WriteString(fmt.Sprintf("\n---\n## Evaluasi untuk $x = %s$\n", FormatNum(targetX)))

		segIdx, segmentDescription := ResolveTargetSegment(points, targetX)
		if segIdx == -1 {
			errMsg := fmt.Sprintf("Error: x=%v di luar rentang data.\n", FormatNum(targetX))
			fmt.Print(errMsg)
			report.WriteString(errMsg)
			continue
		}

		i := segIdx
		a := points[i].Y
		b := m[i]
		dx := targetX - points[i].X
		yTarget := a + b*dx

		PrintStageTitle(fmt.Sprintf("Tahap 3: Pilih Segmen %d", i))
		segmentLine := segmentDescription + "\n"
		fmt.Print(segmentLine)
		report.WriteString("\n### Tahap 3: Pemilihan Segmen\n")
		report.WriteString("- " + strings.TrimSpace(segmentLine) + "\n")

		PrintStageTitle(fmt.Sprintf("Tahap 4: Substitusi Segmen %d", i))
		fmt.Printf("a_%d = %s\nb_%d = %s\n", i, FormatNum(a), i, FormatNum(b))
		report.WriteString(fmt.Sprintf("\n### Tahap 4: Koefisien Segmen %d\n", i))
		report.WriteString(fmt.Sprintf("- $a_%d = %s$\n- $b_%d = %s$\n", i, FormatNum(a), i, FormatNum(b)))

		PrintStageTitle("Tahap 5: Bentuk Persamaan")
		equation := fmt.Sprintf("S(x) = %s + (%s)(x - %s)", FormatFractionOnly(a), FormatFractionOnly(b), FormatFractionOnly(points[i].X))
		fmt.Println(equation)
		report.WriteString("\n### Tahap 5: Persamaan\n")
		report.WriteString(fmt.Sprintf("$%s$\n", equation))

		PrintStageTitle("Tahap 6: Evaluasi Nilai")
		evalLine := fmt.Sprintf(
			"S(%s) = %s + (%s)(%s - %s) = %s\n",
			FormatNum(targetX),
			FormatNum(a),
			FormatNum(b),
			FormatNum(targetX),
			FormatNum(points[i].X),
			FormatNum(yTarget),
		)
		fmt.Print(evalLine)
		fmt.Printf("\nHasil taksiran y: %s\n", FormatNum(yTarget))
		report.WriteString("\n### Tahap 6: Evaluasi\n")
		report.WriteString("- " + strings.TrimSpace(evalLine) + "\n")
		report.WriteString(fmt.Sprintf("\n**Hasil Akhir: $y = %s$**\n", FormatNum(yTarget)))
	}

	return report.String(), nil
}
