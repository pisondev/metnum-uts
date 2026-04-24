package main

import (
	"fmt"
	"os"
)

type Point struct {
	X float64
	Y float64
}

func main() {
	fmt.Println("=== PROGRAM NATURAL CUBIC SPLINE INTERPOLATION ===")
	fmt.Println("------------------------------------------------------------------")

	var nPoints int
	fmt.Print("Masukkan jumlah pasangan titik data: ")
	fmt.Scan(&nPoints)

	points := make([]Point, nPoints)
	fmt.Println("Masukkan koordinat x dan y (contoh: 2 1/2):")
	for i := 0; i < nPoints; i++ {
		var strX, strY string
		fmt.Printf("Titik ke-%d: ", i)
		fmt.Scan(&strX, &strY)
		points[i].X, points[i].Y = ParseInput(strX), ParseInput(strY)
	}

	var nTargets int
	fmt.Print("\nMasukkan jumlah titik x yang ditanyakan: ")
	fmt.Scan(&nTargets)

	targets := make([]float64, nTargets)
	for i := 0; i < nTargets; i++ {
		var strTgt string
		fmt.Printf("Target x-%d: ", i+1)
		fmt.Scan(&strTgt)
		targets[i] = ParseInput(strTgt)
	}

	fmt.Println("\n================= MEMULAI PROSES PERHITUNGAN =================")

	// Jalankan algoritma dan tangkap string laporan
	markdownReport := SolveCubicSpline(points, targets)

	// Simpan ke file markdown yang unik
	filename := GenerateUniqueFilename()
	err := os.WriteFile(filename, []byte(markdownReport), 0644)

	if err != nil {
		fmt.Printf("\nGagal menyimpan file: %v\n", err)
	} else {
		fmt.Printf("\nLaporan berhasil disimpan ke file: %s%s%s\n", ColorBlue, filename, ColorReset)
	}
}
