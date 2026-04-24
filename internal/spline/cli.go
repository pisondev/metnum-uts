package spline

import (
	"fmt"
	"os"
)

type Point struct {
	X float64
	Y float64
}

func RunCLI() {
	fmt.Println("=== PROGRAM INTERPOLASI SPLINE ===")
	fmt.Println("------------------------------------------------------------------")
	fmt.Println("Pilih Metode:")
	fmt.Println("(01) Linear Spline")
	fmt.Println("(02) Quadratic Spline")
	fmt.Println("(03) Natural Cubic Spline")

	var choice int
	fmt.Print("Input: ")
	if _, err := fmt.Scan(&choice); err != nil {
		fmt.Printf("\nInput pilihan metode tidak valid: %v\n", err)
		return
	}

	method, err := GetSplineMethod(choice)
	if err != nil {
		fmt.Printf("\nPilihan metode tidak valid: %v\n", err)
		return
	}

	fmt.Printf("\nMetode terpilih: %s\n", method.DisplayName)

	var nPoints int
	fmt.Print("\nMasukkan jumlah pasangan titik data: ")
	if _, err := fmt.Scan(&nPoints); err != nil {
		fmt.Printf("\nInput jumlah titik tidak valid: %v\n", err)
		return
	}
	if nPoints < 2 {
		fmt.Println("\nJumlah pasangan titik data minimal 2.")
		return
	}

	points := make([]Point, nPoints)
	fmt.Println("Masukkan koordinat x dan y (contoh: 2 1/2):")
	for i := 0; i < nPoints; i++ {
		var strX, strY string
		fmt.Printf("Titik ke-%d: ", i)
		if _, err := fmt.Scan(&strX, &strY); err != nil {
			fmt.Printf("\nInput titik ke-%d tidak valid: %v\n", i, err)
			return
		}

		x, err := ParseInput(strX)
		if err != nil {
			fmt.Printf("\nNilai x pada titik ke-%d tidak valid: %v\n", i, err)
			return
		}

		y, err := ParseInput(strY)
		if err != nil {
			fmt.Printf("\nNilai y pada titik ke-%d tidak valid: %v\n", i, err)
			return
		}

		points[i] = Point{X: x, Y: y}
	}

	if err := ValidatePoints(points); err != nil {
		fmt.Printf("\nData titik tidak valid: %v\n", err)
		return
	}

	var nTargets int
	fmt.Print("\nMasukkan jumlah titik x yang ditanyakan: ")
	if _, err := fmt.Scan(&nTargets); err != nil {
		fmt.Printf("\nInput jumlah target tidak valid: %v\n", err)
		return
	}
	if nTargets < 0 {
		fmt.Println("\nJumlah target tidak boleh negatif.")
		return
	}

	targets := make([]float64, nTargets)
	for i := 0; i < nTargets; i++ {
		var strTgt string
		fmt.Printf("Target x-%d: ", i+1)
		if _, err := fmt.Scan(&strTgt); err != nil {
			fmt.Printf("\nInput target x-%d tidak valid: %v\n", i+1, err)
			return
		}

		target, err := ParseInput(strTgt)
		if err != nil {
			fmt.Printf("\nNilai target x-%d tidak valid: %v\n", i+1, err)
			return
		}

		targets[i] = target
	}

	fmt.Println("\n================= MEMULAI PROSES PERHITUNGAN =================")

	// Jalankan algoritma dan tangkap string laporan
	markdownReport, err := method.Solve(points, targets)
	if err != nil {
		fmt.Printf("\nPerhitungan gagal: %v\n", err)
		return
	}

	// Simpan ke file markdown yang unik
	filename, err := GenerateUniqueFilename(method.ResultsSlug)
	if err != nil {
		fmt.Printf("\nGagal menyiapkan file laporan: %v\n", err)
		return
	}

	err = os.WriteFile(filename, []byte(markdownReport), 0644)

	if err != nil {
		fmt.Printf("\nGagal menyimpan file: %v\n", err)
	} else {
		fmt.Printf("\nLaporan berhasil disimpan ke file: %s%s%s\n", ColorBlue, filename, ColorReset)
	}
}
