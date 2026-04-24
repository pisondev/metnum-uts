package main

import (
	"fmt"
	"math"
	"strings"
)

func SolveCubicSpline(points []Point, targets []float64) string {
	var report strings.Builder
	n := len(points) - 1

	report.WriteString("# Laporan Natural Cubic Spline Interpolation\n\n")
	report.WriteString("## Data Input\n")
	for i, p := range points {
		report.WriteString(fmt.Sprintf("- Titik %d: (%v, %v)\n", i, FormatNum(p.X), FormatNum(p.Y)))
	}

	// TAHAP 1
	PrintStageTitle("Tahap 1: Pemetaan Jarak (h)")
	report.WriteString("\n## Tahap 1: Pemetaan Jarak (h)\n")
	report.WriteString("Rumus Umum: $h_i = x_{i+1} - x_i$\n\n")

	h := make([]float64, n)
	isConstantH := true
	for i := 0; i < n; i++ {
		h[i] = points[i+1].X - points[i].X
		res := fmt.Sprintf("- $h_%d = %v - %v = %s$\n", i, FormatNum(points[i+1].X), FormatNum(points[i].X), FormatNum(h[i]))
		fmt.Print(res)
		report.WriteString(res)
		if i > 0 && math.Abs(h[i]-h[i-1]) > 1e-9 {
			isConstantH = false
		}
	}

	// TAHAP 2
	PrintStageTitle("Tahap 2: Syarat Batas (Natural Boundary)")
	report.WriteString("\n## Tahap 2: Syarat Batas\n")
	boundary := fmt.Sprintf("- $f''(x_0) = 0$\n- $f''(x_%d) = 0$\n", n)
	fmt.Print(boundary)
	report.WriteString(boundary)

	M := make([]float64, len(points))

	// TAHAP 3
	PrintStageTitle("Tahap 3: Eksekusi Batas (Mencari titik interior)")
	report.WriteString("\n## Tahap 3: Eksekusi Batas\n")

	size := n - 1
	A, B, C, D := make([]float64, size), make([]float64, size), make([]float64, size), make([]float64, size)

	if isConstantH {
		report.WriteString("Jarak h konstan, menggunakan rumus sederhana.\n\n")
		hVal := h[0]
		for i := 1; i < n; i++ {
			idx := i - 1
			A[idx], B[idx], C[idx] = 1.0, 4.0, 1.0
			D[idx] = (6.0 / (hVal * hVal)) * (points[i+1].Y - 2*points[i].Y + points[i-1].Y)
			line := fmt.Sprintf("- $i=%d: f''(x_%d) + 4f''(x_%d) + f''(x_%d) = %s$\n", i, i-1, i, i+1, FormatNum(D[idx]))
			fmt.Print(line)
			report.WriteString(line)
		}
	} else {
		report.WriteString("Jarak h bervariasi, menggunakan rumus umum.\n\n")
		for i := 1; i < n; i++ {
			idx := i - 1
			A[idx], B[idx], C[idx] = h[i-1], 2.0*(h[i-1]+h[i]), h[i]
			term1 := (points[i+1].Y - points[i].Y) / h[i]
			term2 := (points[i].Y - points[i-1].Y) / h[i-1]
			D[idx] = 6.0 * (term1 - term2)
			line := fmt.Sprintf("- $i=%d: %sf''(x_%d) + %sf''(x_%d) + %sf''(x_%d) = %s$\n", i, FormatNum(A[idx]), i-1, FormatNum(B[idx]), i, FormatNum(C[idx]), i+1, FormatNum(D[idx]))
			fmt.Print(line)
			report.WriteString(line)
		}
	}

	interiorM := ThomasAlgorithm(A, B, C, D)
	for i := 1; i < n; i++ {
		M[i] = interiorM[i-1]
	}

	report.WriteString("\n### Hasil Nilai $f''(x)$:\n")
	for i := 0; i <= n; i++ {
		res := fmt.Sprintf("- $f''(x_%d) = %s$\n", i, FormatNum(M[i]))
		fmt.Print(res)
		report.WriteString(res)
	}

	for _, targetX := range targets {
		// Output untuk Terminal
		fmt.Printf("\n===================================================\n")
		fmt.Printf("EVALUASI UNTUK X DITANYA = %v\n", FormatNum(targetX))

		// Output untuk Markdown
		report.WriteString(fmt.Sprintf("\n---\n## Evaluasi untuk $x = %v$\n", FormatNum(targetX)))

		segIdx := -1
		for i := 0; i < n; i++ {
			if targetX >= points[i].X && targetX <= points[i+1].X {
				segIdx = i
				break
			}
		}

		if segIdx != -1 {
			i := segIdx
			h_i := h[i]
			a := points[i].Y
			b := (points[i+1].Y-points[i].Y)/h_i - (h_i/6.0)*(2.0*M[i]+M[i+1])
			c := M[i] / 2.0
			d := (M[i+1] - M[i]) / (6.0 * h_i)

			// TAHAP 4
			PrintStageTitle(fmt.Sprintf("Tahap 4: Substitusi Koef Segmen %d", i))
			fmt.Printf("a_%d = %s\nb_%d = %s\nc_%d = %s\nd_%d = %s\n", i, FormatNum(a), i, FormatNum(b), i, FormatNum(c), i, FormatNum(d))

			report.WriteString(fmt.Sprintf("\n### Tahap 4: Koefisien Segmen %d\n", i))
			report.WriteString(fmt.Sprintf("- $a_%d = %s$\n- $b_%d = %s$\n- $c_%d = %s$\n- $d_%d = %s$\n", i, FormatNum(a), i, FormatNum(b), i, FormatNum(c), i, FormatNum(d)))

			// TAHAP 5
			PrintStageTitle("Tahap 5: Bentuk Persamaan")
			fmt.Printf("S(x) = %s + (%s)(x - %v) + (%s)(x - %v)^2 + (%s)(x - %v)^3\n", FormatFractionOnly(a), FormatFractionOnly(b), points[i].X, FormatFractionOnly(c), points[i].X, FormatFractionOnly(d), points[i].X)

			report.WriteString("\n### Tahap 5: Persamaan\n")
			eqMd := fmt.Sprintf("$S(x) = %s + (%s)(x - %v) + (%s)(x - %v)^2 + (%s)(x - %v)^3$\n", FormatFractionOnly(a), FormatFractionOnly(b), points[i].X, FormatFractionOnly(c), points[i].X, FormatFractionOnly(d), points[i].X)
			report.WriteString(eqMd)

			// HASIL AKHIR
			dx := targetX - points[i].X
			yTarget := a + b*dx + c*(dx*dx) + d*(dx*dx*dx)

			fmt.Printf("\nHasil taksiran y: %s\n", FormatNum(yTarget))
			report.WriteString(fmt.Sprintf("\n**Hasil Akhir: $y = %s$**\n", FormatNum(yTarget)))
		} else {
			errMsg := fmt.Sprintf("Error: x=%v di luar rentang data.\n", FormatNum(targetX))
			fmt.Print(errMsg)
			report.WriteString(errMsg)
		}
	}

	return report.String()
}

func ThomasAlgorithm(a, b, c, d []float64) []float64 {
	n := len(d)
	cp, dp, x := make([]float64, n), make([]float64, n), make([]float64, n)
	cp[0], dp[0] = c[0]/b[0], d[0]/b[0]
	for i := 1; i < n; i++ {
		w := b[i] - a[i]*cp[i-1]
		if i < n-1 {
			cp[i] = c[i] / w
		}
		dp[i] = (d[i] - a[i]*dp[i-1]) / w
	}
	x[n-1] = dp[n-1]
	for i := n - 2; i >= 0; i-- {
		x[i] = dp[i] - cp[i]*x[i+1]
	}
	return x
}
