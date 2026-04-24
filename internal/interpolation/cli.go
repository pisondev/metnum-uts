package interpolation

import "fmt"

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

	switch choice {
	case 1:
		fmt.Println("\nNewton Interpolation belum diimplementasikan pada tahap ini.")
	case 2:
		fmt.Println("\nLagrange Interpolation belum diimplementasikan pada tahap ini.")
	default:
		fmt.Println("\nPilihan metode tidak valid.")
	}
}
