package rootfinding

import (
	"fmt"
	"os"
)

func RunCLI() {
	fmt.Println("\n=== MODUL PENCARIAN AKAR ===")
	fmt.Println("------------------------------------------------------------------")
	fmt.Println("Pilih Metode:")
	fmt.Println("(01) Bisection Method")
	fmt.Println("(02) Secant Method")
	fmt.Println("(03) Newton Method")

	var choice int
	fmt.Print("Input: ")
	if _, err := fmt.Scan(&choice); err != nil {
		fmt.Printf("\nInput pilihan metode tidak valid: %v\n", err)
		return
	}

	method, err := GetRootMethod(choice)
	if err != nil {
		fmt.Printf("\nPilihan metode tidak valid: %v\n", err)
		return
	}

	problem, err := readProblemInput(method)
	if err != nil {
		fmt.Printf("\nInput tidak valid: %v\n", err)
		return
	}

	report, err := method.Solve(problem)
	if err != nil {
		fmt.Printf("\nPerhitungan gagal: %v\n", err)
		return
	}

	filename, err := GenerateUniqueFilename(method.ResultsSlug)
	if err != nil {
		fmt.Printf("\nGagal menyiapkan file laporan: %v\n", err)
		return
	}

	if err := os.WriteFile(filename, []byte(report), 0644); err != nil {
		fmt.Printf("\nGagal menyimpan file: %v\n", err)
		return
	}

	fmt.Printf("\nLaporan berhasil disimpan ke file: %s%s%s\n", ColorBlue, filename, ColorReset)
}

func readProblemInput(method RootMethod) (Problem, error) {
	problem := Problem{}

	fmt.Println("\nGunakan operator +, -, *, /, ^ dan tanpa spasi.")
	fmt.Println("Contoh fungsi: x^3-x-2, sin(x)-0.5, exp(-x)-x")

	var functionInput string
	fmt.Print("f(x) = ")
	if _, err := fmt.Scan(&functionInput); err != nil {
		return Problem{}, err
	}

	function, err := CompileExpression(functionInput)
	if err != nil {
		return Problem{}, fmt.Errorf("f(x) tidak valid: %w", err)
	}

	problem.FunctionRaw = functionInput
	problem.Function = function

	if method.Choice == 3 {
		fmt.Println("Contoh turunan: 3*x^2-1, cos(x)+1, -exp(-x)-1")
		var derivativeInput string
		fmt.Print("f'(x) = ")
		if _, err := fmt.Scan(&derivativeInput); err != nil {
			return Problem{}, err
		}

		derivative, err := CompileExpression(derivativeInput)
		if err != nil {
			return Problem{}, fmt.Errorf("f'(x) tidak valid: %w", err)
		}

		problem.DerivativeRaw = derivativeInput
		problem.Derivative = &derivative
	}

	switch method.Choice {
	case 1:
		lower, upper, err := readIntervalInput()
		if err != nil {
			return Problem{}, err
		}
		problem.Lower = lower
		problem.Upper = upper
	case 2:
		x0, x1, err := readTwoInitialGuesses()
		if err != nil {
			return Problem{}, err
		}
		problem.X0 = x0
		problem.X1 = x1
	case 3:
		x0, err := readSingleInitialGuess()
		if err != nil {
			return Problem{}, err
		}
		problem.X0 = x0
	}

	tolerance, maxIterations, err := readStoppingCriteria()
	if err != nil {
		return Problem{}, err
	}
	problem.Tolerance = tolerance
	problem.MaxIterations = maxIterations

	return problem, nil
}

func readIntervalInput() (float64, float64, error) {
	var lowerInput string
	var upperInput string

	fmt.Print("Batas bawah a = ")
	if _, err := fmt.Scan(&lowerInput); err != nil {
		return 0, 0, err
	}
	fmt.Print("Batas atas b = ")
	if _, err := fmt.Scan(&upperInput); err != nil {
		return 0, 0, err
	}

	lower, err := ParseNumberInput(lowerInput)
	if err != nil {
		return 0, 0, fmt.Errorf("batas bawah tidak valid: %w", err)
	}
	upper, err := ParseNumberInput(upperInput)
	if err != nil {
		return 0, 0, fmt.Errorf("batas atas tidak valid: %w", err)
	}

	return lower, upper, nil
}

func readTwoInitialGuesses() (float64, float64, error) {
	var x0Input string
	var x1Input string

	fmt.Print("Tebakan awal x0 = ")
	if _, err := fmt.Scan(&x0Input); err != nil {
		return 0, 0, err
	}
	fmt.Print("Tebakan awal x1 = ")
	if _, err := fmt.Scan(&x1Input); err != nil {
		return 0, 0, err
	}

	x0, err := ParseNumberInput(x0Input)
	if err != nil {
		return 0, 0, fmt.Errorf("x0 tidak valid: %w", err)
	}
	x1, err := ParseNumberInput(x1Input)
	if err != nil {
		return 0, 0, fmt.Errorf("x1 tidak valid: %w", err)
	}

	return x0, x1, nil
}

func readSingleInitialGuess() (float64, error) {
	var x0Input string
	fmt.Print("Tebakan awal x0 = ")
	if _, err := fmt.Scan(&x0Input); err != nil {
		return 0, err
	}

	x0, err := ParseNumberInput(x0Input)
	if err != nil {
		return 0, fmt.Errorf("x0 tidak valid: %w", err)
	}

	return x0, nil
}

func readStoppingCriteria() (float64, int, error) {
	var toleranceInput string
	var maxIterations int

	fmt.Print("Toleransi = ")
	if _, err := fmt.Scan(&toleranceInput); err != nil {
		return 0, 0, err
	}
	fmt.Print("Maksimum iterasi = ")
	if _, err := fmt.Scan(&maxIterations); err != nil {
		return 0, 0, err
	}

	tolerance, err := ParseNumberInput(toleranceInput)
	if err != nil {
		return 0, 0, fmt.Errorf("toleransi tidak valid: %w", err)
	}

	return tolerance, maxIterations, nil
}
