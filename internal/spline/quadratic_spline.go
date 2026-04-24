package spline

import (
	"fmt"
	"strings"
)

func SolveQuadraticSpline(points []Point, targets []float64) (string, error) {
	if err := ValidatePoints(points); err != nil {
		return "", err
	}

	var report strings.Builder
	n := len(points) - 1
	h := make([]float64, n)
	a := make([]float64, n)
	b := make([]float64, n)
	c := make([]float64, n)

	report.WriteString("# Laporan 02. Quadratic Spline Interpolation\n\n")
	WriteReportInputSection(&report, points, targets)

	PrintStageTitle("Tahap 1: Pemetaan Jarak (h)")
	report.WriteString("\n## Tahap 1: Pemetaan Jarak (h)\n")
	report.WriteString("Rumus umum: $h_i = x_{i+1} - x_i$\n\n")

	for i := 0; i < n; i++ {
		h[i] = points[i+1].X - points[i].X
		line := fmt.Sprintf("- $h_%d = %s - %s = %s$\n", i, FormatNum(points[i+1].X), FormatNum(points[i].X), FormatNum(h[i]))
		fmt.Print(line)
		report.WriteString(line)
	}

	PrintStageTitle("Tahap 2: Bentuk Umum dan Syarat Batas")
	report.WriteString("\n## Tahap 2: Bentuk Umum dan Syarat Batas\n")
	report.WriteString("Model tiap segmen: $S_i(x) = a_i + b_i(x - x_i) + c_i(x - x_i)^2$.\n\n")
	report.WriteString("Syarat interpolasi tiap segmen:\n")
	report.WriteString("- $S_i(x_i) = y_i$, sehingga $a_i = y_i$.\n")
	report.WriteString("- $S_i(x_{i+1}) = y_{i+1}$, sehingga $a_i + b_i h_i + c_i h_i^2 = y_{i+1}$.\n")
	report.WriteString("- Kontinuitas turunan pertama di titik sambung: $S'_{i-1}(x_i) = S'_i(x_i)$, sehingga $b_i = b_{i-1} + 2c_{i-1}h_{i-1}$.\n\n")
	report.WriteString("Syarat batas tambahan dipilih di kiri agar sistem kuadratik tunggal, yaitu $S_0''(x_0) = 0$ sehingga $c_0 = 0$.\n")
	boundary := "- $c_0 = 0$\n"
	fmt.Print(boundary)
	report.WriteString(boundary)

	PrintStageTitle("Tahap 3: Menentukan Koefisien Segmen")
	report.WriteString("\n## Tahap 3: Menentukan Koefisien Segmen\n")
	report.WriteString("Koefisien dihitung berurutan dari kiri ke kanan dengan kontinuitas turunan pertama.\n\n")

	for i := 0; i < n; i++ {
		a[i] = points[i].Y
	}

	b[0] = (points[1].Y - points[0].Y) / h[0]
	c[0] = 0
	seg0 := fmt.Sprintf(
		"- Segmen 0: $a_0 = %s$, $c_0 = 0$, dan $b_0 = \\frac{%s - %s}{%s} = %s$\n",
		FormatNum(a[0]),
		FormatNum(points[1].Y),
		FormatNum(points[0].Y),
		FormatNum(h[0]),
		FormatNum(b[0]),
	)
	fmt.Print(seg0)
	report.WriteString(seg0)

	for i := 1; i < n; i++ {
		b[i] = b[i-1] + 2*c[i-1]*h[i-1]
		c[i] = (points[i+1].Y - points[i].Y - b[i]*h[i]) / (h[i] * h[i])

		line := fmt.Sprintf(
			"- Segmen %d: $a_%d = %s$, $b_%d = b_%d + 2c_%d h_%d = %s$, $c_%d = \\frac{%s - %s - (%s)(%s)}{%s^2} = %s$\n",
			i,
			i,
			FormatNum(a[i]),
			i,
			i-1,
			i-1,
			i-1,
			FormatNum(b[i]),
			i,
			FormatNum(points[i+1].Y),
			FormatNum(points[i].Y),
			FormatNum(b[i]),
			FormatNum(h[i]),
			FormatNum(h[i]),
			FormatNum(c[i]),
		)
		fmt.Print(line)
		report.WriteString(line)
	}

	PrintStageTitle("Tahap 4: Persamaan Setiap Segmen")
	report.WriteString("\n## Tahap 4: Persamaan Setiap Segmen\n")
	for i := 0; i < n; i++ {
		line := fmt.Sprintf(
			"- Segmen %d pada interval $[%s, %s]$: $S_%d(x) = %s + (%s)(x - %s) + (%s)(x - %s)^2$\n",
			i,
			FormatNum(points[i].X),
			FormatNum(points[i+1].X),
			i,
			FormatFractionOnly(a[i]),
			FormatFractionOnly(b[i]),
			FormatFractionOnly(points[i].X),
			FormatFractionOnly(c[i]),
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
		dx := targetX - points[i].X
		yTarget := a[i] + b[i]*dx + c[i]*dx*dx

		PrintStageTitle(fmt.Sprintf("Tahap 5: Pilih Segmen %d", i))
		segmentLine := segmentDescription + "\n"
		fmt.Print(segmentLine)
		report.WriteString("\n### Tahap 5: Pemilihan Segmen\n")
		report.WriteString("- " + strings.TrimSpace(segmentLine) + "\n")

		PrintStageTitle(fmt.Sprintf("Tahap 6: Substitusi Segmen %d", i))
		fmt.Printf("a_%d = %s\nb_%d = %s\nc_%d = %s\n", i, FormatNum(a[i]), i, FormatNum(b[i]), i, FormatNum(c[i]))
		report.WriteString(fmt.Sprintf("\n### Tahap 6: Koefisien Segmen %d\n", i))
		report.WriteString(fmt.Sprintf("- $a_%d = %s$\n- $b_%d = %s$\n- $c_%d = %s$\n", i, FormatNum(a[i]), i, FormatNum(b[i]), i, FormatNum(c[i])))

		PrintStageTitle("Tahap 7: Bentuk Persamaan")
		equation := fmt.Sprintf(
			"S(x) = %s + (%s)(x - %s) + (%s)(x - %s)^2",
			FormatFractionOnly(a[i]),
			FormatFractionOnly(b[i]),
			FormatFractionOnly(points[i].X),
			FormatFractionOnly(c[i]),
			FormatFractionOnly(points[i].X),
		)
		fmt.Println(equation)
		report.WriteString("\n### Tahap 7: Persamaan\n")
		report.WriteString(fmt.Sprintf("$%s$\n", equation))

		PrintStageTitle("Tahap 8: Evaluasi Nilai")
		evalLine := fmt.Sprintf(
			"S(%s) = %s + (%s)(%s - %s) + (%s)(%s - %s)^2 = %s\n",
			FormatNum(targetX),
			FormatNum(a[i]),
			FormatNum(b[i]),
			FormatNum(targetX),
			FormatNum(points[i].X),
			FormatNum(c[i]),
			FormatNum(targetX),
			FormatNum(points[i].X),
			FormatNum(yTarget),
		)
		fmt.Print(evalLine)
		fmt.Printf("\nHasil taksiran y: %s\n", FormatNum(yTarget))
		report.WriteString("\n### Tahap 8: Evaluasi\n")
		report.WriteString("- " + strings.TrimSpace(evalLine) + "\n")
		report.WriteString(fmt.Sprintf("\n**Hasil Akhir: $y = %s$**\n", FormatNum(yTarget)))
	}

	return report.String(), nil
}
