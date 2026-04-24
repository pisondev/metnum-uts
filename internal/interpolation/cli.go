package interpolation

import (
	"fmt"
	"os"
)

type Point struct {
	X float64
	Y float64
}

func RunCLI() {
	fmt.Println("\n=== MODUL INTERPOLASI POLINOM ===")
	fmt.Println("------------------------------------------------------------------")
	fmt.Println("Pilih Metode:")
	fmt.Println("(01) Newton Interpolation")
	fmt.Println("(02) Lagrange Interpolation")

	var choice int
	fmt.Print("Input: ")
	if _, err := fmt.Scan(&choice); err != nil {
		fmt.Printf("\nInput pilihan metode tidak valid: %v\n", err)
		return
	}

	method, err := GetInterpolationMethod(choice)
	if err != nil {
		fmt.Printf("\nPilihan metode tidak valid: %v\n", err)
		return
	}

	fmt.Printf("\nMetode terpilih: %s\n", method.DisplayName)

	points, err := readPoints()
	if err != nil {
		fmt.Printf("\nInput titik tidak valid: %v\n", err)
		return
	}

	targets, err := readTargets()
	if err != nil {
		fmt.Printf("\nInput target tidak valid: %v\n", err)
		return
	}

	report, err := method.Solve(points, targets)
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

func readPoints() ([]Point, error) {
	var nPoints int
	fmt.Print("\nMasukkan jumlah pasangan titik data: ")
	if _, err := fmt.Scan(&nPoints); err != nil {
		return nil, err
	}
	if nPoints < 2 {
		return nil, fmt.Errorf("jumlah pasangan titik data minimal 2")
	}

	points := make([]Point, nPoints)
	fmt.Println("Masukkan koordinat x dan y (contoh: 2 1/2):")
	for i := 0; i < nPoints; i++ {
		var strX string
		var strY string

		fmt.Printf("Titik ke-%d: ", i)
		if _, err := fmt.Scan(&strX, &strY); err != nil {
			return nil, fmt.Errorf("gagal membaca titik ke-%d: %w", i, err)
		}

		x, err := ParseNumberInput(strX)
		if err != nil {
			return nil, fmt.Errorf("nilai x pada titik ke-%d tidak valid: %w", i, err)
		}
		y, err := ParseNumberInput(strY)
		if err != nil {
			return nil, fmt.Errorf("nilai y pada titik ke-%d tidak valid: %w", i, err)
		}

		points[i] = Point{X: x, Y: y}
	}

	if err := ValidatePoints(points); err != nil {
		return nil, err
	}

	return points, nil
}

func readTargets() ([]float64, error) {
	var nTargets int
	fmt.Print("\nMasukkan jumlah titik x yang ditanyakan: ")
	if _, err := fmt.Scan(&nTargets); err != nil {
		return nil, err
	}
	if nTargets < 0 {
		return nil, fmt.Errorf("jumlah target tidak boleh negatif")
	}

	targets := make([]float64, nTargets)
	for i := 0; i < nTargets; i++ {
		var strTarget string
		fmt.Printf("Target x-%d: ", i+1)
		if _, err := fmt.Scan(&strTarget); err != nil {
			return nil, fmt.Errorf("gagal membaca target ke-%d: %w", i+1, err)
		}

		target, err := ParseNumberInput(strTarget)
		if err != nil {
			return nil, fmt.Errorf("nilai target x-%d tidak valid: %w", i+1, err)
		}
		targets[i] = target
	}

	return targets, nil
}
