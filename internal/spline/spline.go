package spline

import (
	"fmt"
	"math"
	"strings"
)

func SolveCubicSpline(points []Point, targets []float64) (string, error) {
	var report strings.Builder
	if err := ValidatePoints(points); err != nil {
		return "", err
	}

	n := len(points) - 1

	report.WriteString("# Laporan 03. Natural Cubic Spline Interpolation\n\n")
	WriteReportInputSection(&report, points, targets)

	// TAHAP 1
	PrintStageTitle("Tahap 1: Pemetaan Jarak (h)")
	report.WriteString("\n## Tahap 1: Pemetaan Jarak (h)\n")
	report.WriteString("Rumus Umum: $h_i = x_{i+1} - x_i$\n\n")

	h := make([]float64, n)
	for i := 0; i < n; i++ {
		h[i] = points[i+1].X - points[i].X
		res := fmt.Sprintf("- $h_%d = %v - %v = %s$\n", i, FormatNum(points[i+1].X), FormatNum(points[i].X), FormatNum(h[i]))
		fmt.Print(res)
		report.WriteString(res)
	}

	// TAHAP 2
	PrintStageTitle("Tahap 2: Syarat Batas (Natural Boundary)")
	report.WriteString("\n## Tahap 2: Syarat Batas\n")
	boundary := fmt.Sprintf("- $f''(x_0) = 0$\n- $f''(x_%d) = 0$\n", n)
	fmt.Print(boundary)
	report.WriteString(boundary)

	M := make([]float64, len(points))

	// TAHAP 3
	PrintStageTitle("Tahap 3: Sistem Tridiagonal untuk M_i")
	report.WriteString("\n## Tahap 3: Sistem Tridiagonal untuk M_i\n")
	report.WriteString("Dengan notasi $M_i = f''(x_i)$ dan syarat natural $M_0 = M_n = 0$, untuk tiap titik interior berlaku:\n\n")
	report.WriteString("$h_{i-1}M_{i-1} + 2(h_{i-1}+h_i)M_i + h_iM_{i+1} = 6\\left(\\frac{y_{i+1}-y_i}{h_i} - \\frac{y_i-y_{i-1}}{h_{i-1}}\\right)$\n\n")

	size := n - 1
	A, B, C, D := make([]float64, size), make([]float64, size), make([]float64, size), make([]float64, size)

	if size == 0 {
		msg := "- Tidak ada titik interior, sehingga tidak perlu menyusun sistem tridiagonal tambahan.\n"
		fmt.Print(msg)
		report.WriteString(msg)
	} else {
		report.WriteString("Koefisien sistem selalu disusun dengan rumus umum agar tetap konsisten untuk jarak konstan maupun bervariasi.\n\n")
		for i := 1; i < n; i++ {
			idx := i - 1
			A[idx], B[idx], C[idx] = h[i-1], 2.0*(h[i-1]+h[i]), h[i]
			term1 := (points[i+1].Y - points[i].Y) / h[i]
			term2 := (points[i].Y - points[i-1].Y) / h[i-1]
			D[idx] = 6.0 * (term1 - term2)
			line := fmt.Sprintf("- $i=%d: %sM_%d + %sM_%d + %sM_%d = %s$\n", i, FormatNum(A[idx]), i-1, FormatNum(B[idx]), i, FormatNum(C[idx]), i+1, FormatNum(D[idx]))
			fmt.Print(line)
			report.WriteString(line)
		}
	}

	if size > 0 {
		interiorM, err := ThomasAlgorithm(A, B, C, D)
		if err != nil {
			return "", err
		}

		for i := 1; i < n; i++ {
			M[i] = interiorM[i-1]
		}
	}

	report.WriteString("\n### Hasil Nilai $M_i = f''(x_i)$:\n")
	for i := 0; i <= n; i++ {
		res := fmt.Sprintf("- $M_%d = f''(x_%d) = %s$\n", i, i, FormatNum(M[i]))
		fmt.Print(res)
		report.WriteString(res)
	}

	for _, targetX := range targets {
		// Output untuk Terminal
		fmt.Printf("\n===================================================\n")
		fmt.Printf("EVALUASI UNTUK X DITANYA = %v\n", FormatNum(targetX))

		report.WriteString(fmt.Sprintf("\n---\n## Evaluasi untuk $x = %v$\n", FormatNum(targetX)))

		segIdx, segmentDescription := ResolveTargetSegment(points, targetX)

		if segIdx != -1 {
			i := segIdx
			h_i := h[i]
			a := points[i].Y
			b := (points[i+1].Y-points[i].Y)/h_i - (h_i/6.0)*(2.0*M[i]+M[i+1])
			c := M[i] / 2.0
			d := (M[i+1] - M[i]) / (6.0 * h_i)

			PrintStageTitle(fmt.Sprintf("Tahap 4: Pilih Segmen %d", i))
			segmentLine := segmentDescription + "\n"
			fmt.Print(segmentLine)
			report.WriteString("\n### Tahap 4: Pemilihan Segmen\n")
			report.WriteString("- " + strings.TrimSpace(segmentLine) + "\n")

			PrintStageTitle(fmt.Sprintf("Tahap 5: Substitusi Koef Segmen %d", i))
			fmt.Printf("a_%d = %s\nb_%d = %s\nc_%d = %s\nd_%d = %s\n", i, FormatNum(a), i, FormatNum(b), i, FormatNum(c), i, FormatNum(d))

			report.WriteString(fmt.Sprintf("\n### Tahap 5: Koefisien Segmen %d\n", i))
			report.WriteString(fmt.Sprintf("- $a_%d = %s$\n- $b_%d = %s$\n- $c_%d = %s$\n- $d_%d = %s$\n", i, FormatNum(a), i, FormatNum(b), i, FormatNum(c), i, FormatNum(d)))

			PrintStageTitle("Tahap 6: Bentuk Persamaan")
			equation := fmt.Sprintf(
				"S(x) = %s + (%s)(x - %s) + (%s)(x - %s)^2 + (%s)(x - %s)^3",
				FormatFractionOnly(a),
				FormatFractionOnly(b),
				FormatFractionOnly(points[i].X),
				FormatFractionOnly(c),
				FormatFractionOnly(points[i].X),
				FormatFractionOnly(d),
				FormatFractionOnly(points[i].X),
			)
			fmt.Println(equation)

			report.WriteString("\n### Tahap 6: Persamaan\n")
			report.WriteString(fmt.Sprintf("$%s$\n", equation))

			dx := targetX - points[i].X
			yTarget := a + b*dx + c*(dx*dx) + d*(dx*dx*dx)

			PrintStageTitle("Tahap 7: Evaluasi Nilai")
			evalLine := fmt.Sprintf(
				"S(%s) = %s + (%s)(%s - %s) + (%s)(%s - %s)^2 + (%s)(%s - %s)^3 = %s\n",
				FormatNum(targetX),
				FormatNum(a),
				FormatNum(b),
				FormatNum(targetX),
				FormatNum(points[i].X),
				FormatNum(c),
				FormatNum(targetX),
				FormatNum(points[i].X),
				FormatNum(d),
				FormatNum(targetX),
				FormatNum(points[i].X),
				FormatNum(yTarget),
			)
			fmt.Print(evalLine)
			fmt.Printf("\nHasil taksiran y: %s\n", FormatNum(yTarget))
			report.WriteString("\n### Tahap 7: Evaluasi\n")
			report.WriteString("- " + strings.TrimSpace(evalLine) + "\n")
			report.WriteString(fmt.Sprintf("\n**Hasil Akhir: $y = %s$**\n", FormatNum(yTarget)))
		} else {
			errMsg := fmt.Sprintf("Error: x=%v di luar rentang data.\n", FormatNum(targetX))
			fmt.Print(errMsg)
			report.WriteString(errMsg)
		}
	}

	return report.String(), nil
}

func ThomasAlgorithm(a, b, c, d []float64) ([]float64, error) {
	n := len(d)
	if n == 0 {
		return nil, nil
	}

	cp, dp, x := make([]float64, n), make([]float64, n), make([]float64, n)
	if math.Abs(b[0]) < 1e-12 {
		return nil, fmt.Errorf("sistem spline tidak valid: pivot nol pada baris 0")
	}

	cp[0], dp[0] = c[0]/b[0], d[0]/b[0]
	for i := 1; i < n; i++ {
		w := b[i] - a[i]*cp[i-1]
		if math.Abs(w) < 1e-12 {
			return nil, fmt.Errorf("sistem spline tidak valid: pivot nol pada baris %d", i)
		}

		if i < n-1 {
			cp[i] = c[i] / w
		}
		dp[i] = (d[i] - a[i]*dp[i-1]) / w
	}
	x[n-1] = dp[n-1]
	for i := n - 2; i >= 0; i-- {
		x[i] = dp[i] - cp[i]*x[i+1]
	}
	return x, nil
}
