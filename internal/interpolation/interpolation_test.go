package interpolation

import (
	"strings"
	"testing"
)

func TestGetInterpolationMethod(t *testing.T) {
	method, err := GetInterpolationMethod(1)
	if err != nil {
		t.Fatalf("GetInterpolationMethod returned error: %v", err)
	}

	if method.Name != "Newton Interpolation" {
		t.Fatalf("unexpected method name: %q", method.Name)
	}

	if method.ResultsSlug != "polynomial-interpolation/01-newton-interpolation" {
		t.Fatalf("unexpected results slug: %q", method.ResultsSlug)
	}
}

func TestSolveNewtonInterpolation(t *testing.T) {
	points := []Point{{X: 1, Y: 1}, {X: 2, Y: 4}, {X: 3, Y: 9}}
	report, err := SolveNewtonInterpolation(points, []float64{2.5})
	if err != nil {
		t.Fatalf("SolveNewtonInterpolation returned error: %v", err)
	}

	if !strings.Contains(report, "P(2.5) = 6.25") {
		t.Fatalf("expected interpolation value in report, got:\n%s", report)
	}

	if !strings.Contains(report, "Koefisien Divided Differences") {
		t.Fatalf("expected divided differences section, got:\n%s", report)
	}
}

func TestSolveLagrangeInterpolation(t *testing.T) {
	points := []Point{{X: 1, Y: 1}, {X: 2, Y: 4}, {X: 3, Y: 9}}
	report, err := SolveLagrangeInterpolation(points, []float64{2.5})
	if err != nil {
		t.Fatalf("SolveLagrangeInterpolation returned error: %v", err)
	}

	if !strings.Contains(report, "P(2.5) = 6.25") {
		t.Fatalf("expected interpolation value in report, got:\n%s", report)
	}

	if !strings.Contains(report, "Basis Lagrange") {
		t.Fatalf("expected basis section, got:\n%s", report)
	}
}

func TestValidatePointsRejectsNonIncreasingX(t *testing.T) {
	err := ValidatePoints([]Point{{X: 1, Y: 2}, {X: 1, Y: 3}})
	if err == nil {
		t.Fatal("expected validation error for duplicate x values")
	}
}
