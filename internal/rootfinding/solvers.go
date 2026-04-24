package rootfinding

import (
	"fmt"
	"math"
	"strings"
)

type RootSolution struct {
	Root          float64
	FunctionValue float64
	ApproxError   float64
	Iterations    int
	StopReason    string
}

type BisectionIteration struct {
	Iteration     int
	A             float64
	B             float64
	C             float64
	FC            float64
	IntervalWidth float64
}

type SecantIteration struct {
	Iteration   int
	XPrev       float64
	XCurr       float64
	XNext       float64
	FXNext      float64
	ApproxError float64
}

type NewtonIteration struct {
	Iteration   int
	XCurrent    float64
	FXCurrent   float64
	DFXCurrent  float64
	XNext       float64
	ApproxError float64
}

func SolveBisection(problem Problem) (string, error) {
	solution, iterations, err := runBisection(problem)
	if err != nil {
		return "", err
	}
	return buildBisectionReport(problem, solution, iterations), nil
}

func SolveSecant(problem Problem) (string, error) {
	solution, iterations, err := runSecant(problem)
	if err != nil {
		return "", err
	}
	return buildSecantReport(problem, solution, iterations), nil
}

func SolveNewton(problem Problem) (string, error) {
	solution, iterations, err := runNewton(problem)
	if err != nil {
		return "", err
	}
	return buildNewtonReport(problem, solution, iterations), nil
}

func runBisection(problem Problem) (RootSolution, []BisectionIteration, error) {
	if err := validateCommonProblem(problem); err != nil {
		return RootSolution{}, nil, err
	}
	if problem.Lower >= problem.Upper {
		return RootSolution{}, nil, fmt.Errorf("interval awal harus memenuhi a < b")
	}

	fa, err := problem.Function.Eval(problem.Lower)
	if err != nil {
		return RootSolution{}, nil, fmt.Errorf("gagal mengevaluasi f(a): %w", err)
	}
	fb, err := problem.Function.Eval(problem.Upper)
	if err != nil {
		return RootSolution{}, nil, fmt.Errorf("gagal mengevaluasi f(b): %w", err)
	}

	if math.Abs(fa) <= problem.Tolerance {
		return RootSolution{
			Root:          problem.Lower,
			FunctionValue: fa,
			ApproxError:   0,
			Iterations:    0,
			StopReason:    "batas bawah awal sudah memenuhi toleransi",
		}, nil, nil
	}
	if math.Abs(fb) <= problem.Tolerance {
		return RootSolution{
			Root:          problem.Upper,
			FunctionValue: fb,
			ApproxError:   0,
			Iterations:    0,
			StopReason:    "batas atas awal sudah memenuhi toleransi",
		}, nil, nil
	}
	if fa*fb > 0 {
		return RootSolution{}, nil, fmt.Errorf("f(a) dan f(b) harus memiliki tanda berbeda")
	}

	a := problem.Lower
	b := problem.Upper
	iterations := make([]BisectionIteration, 0, problem.MaxIterations)
	lastSolution := RootSolution{Root: a, FunctionValue: fa, ApproxError: math.Abs(b - a), Iterations: 0}

	for i := 1; i <= problem.MaxIterations; i++ {
		c := (a + b) / 2
		fc, err := problem.Function.Eval(c)
		if err != nil {
			return RootSolution{}, nil, fmt.Errorf("gagal mengevaluasi f(c) pada iterasi %d: %w", i, err)
		}

		intervalWidth := math.Abs(b-a) / 2
		iterations = append(iterations, BisectionIteration{
			Iteration:     i,
			A:             a,
			B:             b,
			C:             c,
			FC:            fc,
			IntervalWidth: intervalWidth,
		})

		lastSolution = RootSolution{
			Root:          c,
			FunctionValue: fc,
			ApproxError:   intervalWidth,
			Iterations:    i,
		}

		if math.Abs(fc) <= problem.Tolerance {
			lastSolution.StopReason = "nilai |f(c)| sudah memenuhi toleransi"
			return lastSolution, iterations, nil
		}
		if intervalWidth <= problem.Tolerance {
			lastSolution.StopReason = "lebar interval sudah memenuhi toleransi"
			return lastSolution, iterations, nil
		}

		if fa*fc < 0 {
			b = c
			fb = fc
		} else {
			a = c
			fa = fc
		}
	}

	lastSolution.StopReason = "batas iterasi maksimum tercapai"
	_ = fb
	return lastSolution, iterations, nil
}

func runSecant(problem Problem) (RootSolution, []SecantIteration, error) {
	if err := validateCommonProblem(problem); err != nil {
		return RootSolution{}, nil, err
	}
	if NearlyEqual(problem.X0, problem.X1) {
		return RootSolution{}, nil, fmt.Errorf("x0 dan x1 tidak boleh sama")
	}

	xPrev := problem.X0
	xCurr := problem.X1
	fPrev, err := problem.Function.Eval(xPrev)
	if err != nil {
		return RootSolution{}, nil, fmt.Errorf("gagal mengevaluasi f(x0): %w", err)
	}
	fCurr, err := problem.Function.Eval(xCurr)
	if err != nil {
		return RootSolution{}, nil, fmt.Errorf("gagal mengevaluasi f(x1): %w", err)
	}

	if math.Abs(fPrev) <= problem.Tolerance {
		return RootSolution{Root: xPrev, FunctionValue: fPrev, ApproxError: 0, Iterations: 0, StopReason: "x0 awal sudah memenuhi toleransi"}, nil, nil
	}
	if math.Abs(fCurr) <= problem.Tolerance {
		return RootSolution{Root: xCurr, FunctionValue: fCurr, ApproxError: 0, Iterations: 0, StopReason: "x1 awal sudah memenuhi toleransi"}, nil, nil
	}

	iterations := make([]SecantIteration, 0, problem.MaxIterations)
	lastSolution := RootSolution{Root: xCurr, FunctionValue: fCurr, ApproxError: math.Abs(xCurr - xPrev), Iterations: 0}

	for i := 1; i <= problem.MaxIterations; i++ {
		denominator := fCurr - fPrev
		if NearlyEqual(denominator, 0) {
			return RootSolution{}, nil, fmt.Errorf("penyebut rumus secant bernilai nol pada iterasi %d", i)
		}

		xNext := xCurr - fCurr*(xCurr-xPrev)/denominator
		fNext, err := problem.Function.Eval(xNext)
		if err != nil {
			return RootSolution{}, nil, fmt.Errorf("gagal mengevaluasi f(x_{k+1}) pada iterasi %d: %w", i, err)
		}

		approxError := math.Abs(xNext - xCurr)
		iterations = append(iterations, SecantIteration{
			Iteration:   i,
			XPrev:       xPrev,
			XCurr:       xCurr,
			XNext:       xNext,
			FXNext:      fNext,
			ApproxError: approxError,
		})

		lastSolution = RootSolution{
			Root:          xNext,
			FunctionValue: fNext,
			ApproxError:   approxError,
			Iterations:    i,
		}

		if math.Abs(fNext) <= problem.Tolerance {
			lastSolution.StopReason = "nilai |f(x)| sudah memenuhi toleransi"
			return lastSolution, iterations, nil
		}
		if approxError <= problem.Tolerance {
			lastSolution.StopReason = "selisih dua hampiran terakhir sudah memenuhi toleransi"
			return lastSolution, iterations, nil
		}

		xPrev = xCurr
		fPrev = fCurr
		xCurr = xNext
		fCurr = fNext
	}

	lastSolution.StopReason = "batas iterasi maksimum tercapai"
	return lastSolution, iterations, nil
}

func runNewton(problem Problem) (RootSolution, []NewtonIteration, error) {
	if err := validateCommonProblem(problem); err != nil {
		return RootSolution{}, nil, err
	}
	if problem.Derivative == nil {
		return RootSolution{}, nil, fmt.Errorf("turunan f'(x) wajib diberikan untuk metode newton")
	}

	xCurrent := problem.X0
	fCurrent, err := problem.Function.Eval(xCurrent)
	if err != nil {
		return RootSolution{}, nil, fmt.Errorf("gagal mengevaluasi f(x0): %w", err)
	}

	if math.Abs(fCurrent) <= problem.Tolerance {
		return RootSolution{Root: xCurrent, FunctionValue: fCurrent, ApproxError: 0, Iterations: 0, StopReason: "x0 awal sudah memenuhi toleransi"}, nil, nil
	}

	iterations := make([]NewtonIteration, 0, problem.MaxIterations)
	lastSolution := RootSolution{Root: xCurrent, FunctionValue: fCurrent, ApproxError: 0, Iterations: 0}

	for i := 1; i <= problem.MaxIterations; i++ {
		dfCurrent, err := problem.Derivative.Eval(xCurrent)
		if err != nil {
			return RootSolution{}, nil, fmt.Errorf("gagal mengevaluasi f'(x) pada iterasi %d: %w", i, err)
		}
		if NearlyEqual(dfCurrent, 0) {
			return RootSolution{}, nil, fmt.Errorf("nilai f'(x) bernilai nol pada iterasi %d", i)
		}

		xNext := xCurrent - fCurrent/dfCurrent
		fNext, err := problem.Function.Eval(xNext)
		if err != nil {
			return RootSolution{}, nil, fmt.Errorf("gagal mengevaluasi f(x_{k+1}) pada iterasi %d: %w", i, err)
		}

		approxError := math.Abs(xNext - xCurrent)
		iterations = append(iterations, NewtonIteration{
			Iteration:   i,
			XCurrent:    xCurrent,
			FXCurrent:   fCurrent,
			DFXCurrent:  dfCurrent,
			XNext:       xNext,
			ApproxError: approxError,
		})

		lastSolution = RootSolution{
			Root:          xNext,
			FunctionValue: fNext,
			ApproxError:   approxError,
			Iterations:    i,
		}

		if math.Abs(fNext) <= problem.Tolerance {
			lastSolution.StopReason = "nilai |f(x)| sudah memenuhi toleransi"
			return lastSolution, iterations, nil
		}
		if approxError <= problem.Tolerance {
			lastSolution.StopReason = "selisih dua hampiran terakhir sudah memenuhi toleransi"
			return lastSolution, iterations, nil
		}

		xCurrent = xNext
		fCurrent = fNext
	}

	lastSolution.StopReason = "batas iterasi maksimum tercapai"
	return lastSolution, iterations, nil
}

func validateCommonProblem(problem Problem) error {
	if problem.MaxIterations <= 0 {
		return fmt.Errorf("maksimum iterasi harus lebih besar dari nol")
	}
	if problem.Tolerance <= 0 {
		return fmt.Errorf("toleransi harus lebih besar dari nol")
	}
	return nil
}

func buildBisectionReport(problem Problem, solution RootSolution, iterations []BisectionIteration) string {
	var report strings.Builder
	report.WriteString("# Laporan 01. Bisection Method\n\n")
	writeCommonInputSection(&report, problem)
	report.WriteString(fmt.Sprintf("- Interval awal: [%s, %s]\n\n", FormatReal(problem.Lower), FormatReal(problem.Upper)))

	report.WriteString("## Iterasi\n")
	if len(iterations) == 0 {
		report.WriteString("- Tidak ada iterasi tambahan karena akar sudah ditemukan pada interval awal.\n\n")
	} else {
		report.WriteString("| Iterasi | a | b | c | f(c) | Lebar Interval |\n")
		report.WriteString("| --- | --- | --- | --- | --- | --- |\n")
		for _, iteration := range iterations {
			report.WriteString(fmt.Sprintf("| %d | %s | %s | %s | %s | %s |\n",
				iteration.Iteration,
				FormatReal(iteration.A),
				FormatReal(iteration.B),
				FormatReal(iteration.C),
				FormatReal(iteration.FC),
				FormatReal(iteration.IntervalWidth),
			))
		}
		report.WriteString("\n")
	}

	writeSolutionSection(&report, solution)
	return report.String()
}

func buildSecantReport(problem Problem, solution RootSolution, iterations []SecantIteration) string {
	var report strings.Builder
	report.WriteString("# Laporan 02. Secant Method\n\n")
	writeCommonInputSection(&report, problem)
	report.WriteString(fmt.Sprintf("- Tebakan awal: x0 = %s, x1 = %s\n\n", FormatReal(problem.X0), FormatReal(problem.X1)))

	report.WriteString("## Iterasi\n")
	if len(iterations) == 0 {
		report.WriteString("- Tidak ada iterasi tambahan karena akar sudah ditemukan pada tebakan awal.\n\n")
	} else {
		report.WriteString("| Iterasi | x_(k-1) | x_k | x_(k+1) | f(x_(k+1)) | Galat |\n")
		report.WriteString("| --- | --- | --- | --- | --- | --- |\n")
		for _, iteration := range iterations {
			report.WriteString(fmt.Sprintf("| %d | %s | %s | %s | %s | %s |\n",
				iteration.Iteration,
				FormatReal(iteration.XPrev),
				FormatReal(iteration.XCurr),
				FormatReal(iteration.XNext),
				FormatReal(iteration.FXNext),
				FormatReal(iteration.ApproxError),
			))
		}
		report.WriteString("\n")
	}

	writeSolutionSection(&report, solution)
	return report.String()
}

func buildNewtonReport(problem Problem, solution RootSolution, iterations []NewtonIteration) string {
	var report strings.Builder
	report.WriteString("# Laporan 03. Newton Method\n\n")
	writeCommonInputSection(&report, problem)
	report.WriteString(fmt.Sprintf("- Turunan: `$f'(x) = %s$`\n", problem.DerivativeRaw))
	report.WriteString(fmt.Sprintf("- Tebakan awal: x0 = %s\n\n", FormatReal(problem.X0)))

	report.WriteString("## Iterasi\n")
	if len(iterations) == 0 {
		report.WriteString("- Tidak ada iterasi tambahan karena akar sudah ditemukan pada tebakan awal.\n\n")
	} else {
		report.WriteString("| Iterasi | x_k | f(x_k) | f'(x_k) | x_(k+1) | Galat |\n")
		report.WriteString("| --- | --- | --- | --- | --- | --- |\n")
		for _, iteration := range iterations {
			report.WriteString(fmt.Sprintf("| %d | %s | %s | %s | %s | %s |\n",
				iteration.Iteration,
				FormatReal(iteration.XCurrent),
				FormatReal(iteration.FXCurrent),
				FormatReal(iteration.DFXCurrent),
				FormatReal(iteration.XNext),
				FormatReal(iteration.ApproxError),
			))
		}
		report.WriteString("\n")
	}

	writeSolutionSection(&report, solution)
	return report.String()
}

func writeCommonInputSection(report *strings.Builder, problem Problem) {
	report.WriteString("## Data Input\n")
	report.WriteString(fmt.Sprintf("- Fungsi: `$f(x) = %s$`\n", problem.FunctionRaw))
	report.WriteString(fmt.Sprintf("- Toleransi: %s\n", FormatReal(problem.Tolerance)))
	report.WriteString(fmt.Sprintf("- Maksimum iterasi: %d\n", problem.MaxIterations))
}

func writeSolutionSection(report *strings.Builder, solution RootSolution) {
	report.WriteString("## Hasil Akhir\n")
	report.WriteString(fmt.Sprintf("- Akar hampiran: %s\n", FormatReal(solution.Root)))
	report.WriteString(fmt.Sprintf("- Nilai fungsi: %s\n", FormatReal(solution.FunctionValue)))
	report.WriteString(fmt.Sprintf("- Galat hampiran: %s\n", FormatReal(solution.ApproxError)))
	report.WriteString(fmt.Sprintf("- Jumlah iterasi: %d\n", solution.Iterations))
	report.WriteString(fmt.Sprintf("- Alasan berhenti: %s\n", solution.StopReason))
}
