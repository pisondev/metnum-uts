package spline

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGetSplineMethod(t *testing.T) {
	method, err := GetSplineMethod(2)
	if err != nil {
		t.Fatalf("GetSplineMethod returned error: %v", err)
	}

	if method.Name != "Quadratic Spline" {
		t.Fatalf("unexpected method name: %q", method.Name)
	}

	if method.DisplayName != "02. Quadratic Spline" {
		t.Fatalf("unexpected method display name: %q", method.DisplayName)
	}

	if method.ResultsSlug != "02-quadratic-spline" {
		t.Fatalf("unexpected results slug: %q", method.ResultsSlug)
	}
}

func TestSolveLinearSpline(t *testing.T) {
	points := []Point{{X: 0, Y: 0}, {X: 2, Y: 4}, {X: 5, Y: 7}}
	targets := []float64{3}

	report, err := SolveLinearSpline(points, targets)
	if err != nil {
		t.Fatalf("SolveLinearSpline returned error: %v", err)
	}

	if !strings.Contains(report, "**Hasil Akhir: $y = 5$**") {
		t.Fatalf("expected linear spline result in report, got:\n%s", report)
	}

	if !strings.Contains(report, "Segmen 1") {
		t.Fatalf("expected segment details in report, got:\n%s", report)
	}
}

func TestSolveQuadraticSpline(t *testing.T) {
	points := []Point{{X: 0, Y: 0}, {X: 1, Y: 1}, {X: 2, Y: 4}}
	targets := []float64{1.5}

	report, err := SolveQuadraticSpline(points, targets)
	if err != nil {
		t.Fatalf("SolveQuadraticSpline returned error: %v", err)
	}

	if !strings.Contains(report, "**Hasil Akhir: $y = 2$**") {
		t.Fatalf("expected quadratic spline result in report, got:\n%s", report)
	}

	if !strings.Contains(report, "$c_0 = 0$") {
		t.Fatalf("expected quadratic boundary condition in report, got:\n%s", report)
	}
}

func TestSolveCubicSplineTwoPointsLinear(t *testing.T) {
	points := []Point{{X: 0, Y: 0}, {X: 2, Y: 4}}
	targets := []float64{1}

	report, err := SolveCubicSpline(points, targets)
	if err != nil {
		t.Fatalf("SolveCubicSpline returned error: %v", err)
	}

	if !strings.Contains(report, "**Hasil Akhir: $y = 2$**") {
		t.Fatalf("expected linear interpolation result in report, got:\n%s", report)
	}
}

func TestSolveCubicSplineCollinearConstantH(t *testing.T) {
	points := []Point{{X: 0, Y: 0}, {X: 1, Y: 1}, {X: 2, Y: 2}}
	targets := []float64{1.5}

	report, err := SolveCubicSpline(points, targets)
	if err != nil {
		t.Fatalf("SolveCubicSpline returned error: %v", err)
	}

	if !strings.Contains(report, "**Hasil Akhir: $y = 3/2 (≈ 1.5000)$**") {
		t.Fatalf("expected constant-h collinear interpolation result in report, got:\n%s", report)
	}
}

func TestSolveCubicSplineCollinearVariableH(t *testing.T) {
	points := []Point{{X: 0, Y: 0}, {X: 2, Y: 2}, {X: 5, Y: 5}}
	targets := []float64{3}

	report, err := SolveCubicSpline(points, targets)
	if err != nil {
		t.Fatalf("SolveCubicSpline returned error: %v", err)
	}

	if !strings.Contains(report, "**Hasil Akhir: $y = 3$**") {
		t.Fatalf("expected variable-h collinear interpolation result in report, got:\n%s", report)
	}
}

func TestSolveCubicSplineRejectsNonIncreasingX(t *testing.T) {
	points := []Point{{X: 1, Y: 1}, {X: 1, Y: 2}, {X: 2, Y: 3}}

	_, err := SolveCubicSpline(points, nil)
	if err == nil {
		t.Fatal("expected validation error for duplicate x values")
	}
}

func TestParseInputRejectsInvalidFraction(t *testing.T) {
	_, err := ParseInput("1/0")
	if err == nil {
		t.Fatal("expected division-by-zero validation error")
	}
}

func TestParseInputRejectsNonFiniteValue(t *testing.T) {
	_, err := ParseInput("NaN")
	if err == nil {
		t.Fatal("expected non-finite validation error")
	}
}

func TestGenerateUniqueFilenameUsesResultsDir(t *testing.T) {
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd failed: %v", err)
	}

	tempDir := t.TempDir()
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Chdir to temp dir failed: %v", err)
	}
	defer func() {
		_ = os.Chdir(oldWd)
	}()

	filename, err := GenerateUniqueFilename("01-linear-spline")
	if err != nil {
		t.Fatalf("GenerateUniqueFilename returned error: %v", err)
	}

	expectedDir := filepath.Join("results", "01-linear-spline")
	if filepath.Dir(filename) != expectedDir {
		t.Fatalf("expected file inside results directory, got %q", filename)
	}

	if _, err := os.Stat(expectedDir); err != nil {
		t.Fatalf("expected results directory to exist: %v", err)
	}
}
