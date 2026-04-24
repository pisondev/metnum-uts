package rootfinding

import (
	"math"
	"strings"
	"testing"
)

func TestGetRootMethod(t *testing.T) {
	method, err := GetRootMethod(2)
	if err != nil {
		t.Fatalf("GetRootMethod returned error: %v", err)
	}

	if method.Name != "Secant Method" {
		t.Fatalf("unexpected method name: %q", method.Name)
	}

	if method.ResultsSlug != "root-finding/02-secant-method" {
		t.Fatalf("unexpected results slug: %q", method.ResultsSlug)
	}
}

func TestCompileExpressionEvaluatesPolynomial(t *testing.T) {
	expression, err := CompileExpression("x^3-x-2")
	if err != nil {
		t.Fatalf("CompileExpression returned error: %v", err)
	}

	value, err := expression.Eval(2)
	if err != nil {
		t.Fatalf("Eval returned error: %v", err)
	}

	if !almostEqual(value, 4, 1e-12) {
		t.Fatalf("unexpected expression value: got %v want %v", value, 4.0)
	}
}

func TestRunBisectionFindsRoot(t *testing.T) {
	function, err := CompileExpression("x^2-4")
	if err != nil {
		t.Fatalf("CompileExpression returned error: %v", err)
	}

	solution, iterations, err := runBisection(Problem{
		FunctionRaw:   "x^2-4",
		Function:      function,
		Lower:         1,
		Upper:         3,
		Tolerance:     1e-8,
		MaxIterations: 100,
	})
	if err != nil {
		t.Fatalf("runBisection returned error: %v", err)
	}

	if len(iterations) == 0 {
		t.Fatal("expected bisection iterations to be recorded")
	}

	if !almostEqual(solution.Root, 2, 1e-6) {
		t.Fatalf("unexpected root approximation: got %v want near 2", solution.Root)
	}
}

func TestRunSecantFindsRoot(t *testing.T) {
	function, err := CompileExpression("x^2-4")
	if err != nil {
		t.Fatalf("CompileExpression returned error: %v", err)
	}

	solution, iterations, err := runSecant(Problem{
		FunctionRaw:   "x^2-4",
		Function:      function,
		X0:            1,
		X1:            3,
		Tolerance:     1e-8,
		MaxIterations: 100,
	})
	if err != nil {
		t.Fatalf("runSecant returned error: %v", err)
	}

	if len(iterations) == 0 {
		t.Fatal("expected secant iterations to be recorded")
	}

	if !almostEqual(solution.Root, 2, 1e-6) {
		t.Fatalf("unexpected root approximation: got %v want near 2", solution.Root)
	}
}

func TestRunNewtonFindsRoot(t *testing.T) {
	function, err := CompileExpression("x^2-4")
	if err != nil {
		t.Fatalf("CompileExpression returned error: %v", err)
	}
	derivative, err := CompileExpression("2*x")
	if err != nil {
		t.Fatalf("CompileExpression derivative returned error: %v", err)
	}

	solution, iterations, err := runNewton(Problem{
		FunctionRaw:   "x^2-4",
		Function:      function,
		DerivativeRaw: "2*x",
		Derivative:    &derivative,
		X0:            3,
		Tolerance:     1e-8,
		MaxIterations: 100,
	})
	if err != nil {
		t.Fatalf("runNewton returned error: %v", err)
	}

	if len(iterations) == 0 {
		t.Fatal("expected newton iterations to be recorded")
	}

	if !almostEqual(solution.Root, 2, 1e-6) {
		t.Fatalf("unexpected root approximation: got %v want near 2", solution.Root)
	}
}

func TestSolveBisectionBuildsReport(t *testing.T) {
	function, err := CompileExpression("x^2-4")
	if err != nil {
		t.Fatalf("CompileExpression returned error: %v", err)
	}

	report, err := SolveBisection(Problem{
		FunctionRaw:   "x^2-4",
		Function:      function,
		Lower:         1,
		Upper:         3,
		Tolerance:     1e-6,
		MaxIterations: 100,
	})
	if err != nil {
		t.Fatalf("SolveBisection returned error: %v", err)
	}

	if !strings.Contains(report, "# Laporan 01. Bisection Method") {
		t.Fatalf("expected report title, got:\n%s", report)
	}

	if !strings.Contains(report, "Akar hampiran") {
		t.Fatalf("expected final result section, got:\n%s", report)
	}
}

func almostEqual(a, b, tolerance float64) bool {
	return math.Abs(a-b) <= tolerance
}
