package app

import (
	"fmt"

	"remed-uts/internal/interpolation"
	"remed-uts/internal/rootfinding"
	"remed-uts/internal/spline"
)

type Module struct {
	Choice      int
	DisplayName string
	Run         func()
}

func RunCLI() {
	fmt.Println("=== PROGRAM METODE NUMERIK ===")
	fmt.Println("------------------------------------------------------------------")
	fmt.Println("Pilih Modul:")
	fmt.Println("(01) Pencarian Akar")
	fmt.Println("(02) Interpolasi Polinom")
	fmt.Println("(03) Interpolasi Spline")

	var choice int
	fmt.Print("Input: ")
	if _, err := fmt.Scan(&choice); err != nil {
		fmt.Printf("\nInput pilihan modul tidak valid: %v\n", err)
		return
	}

	module, err := getModule(choice)
	if err != nil {
		fmt.Printf("\nPilihan modul tidak valid: %v\n", err)
		return
	}

	fmt.Printf("\nModul terpilih: %s\n", module.DisplayName)
	module.Run()
}

func getModule(choice int) (Module, error) {
	switch choice {
	case 1:
		return Module{
			Choice:      1,
			DisplayName: "01. Pencarian Akar",
			Run:         rootfinding.RunCLI,
		}, nil
	case 2:
		return Module{
			Choice:      2,
			DisplayName: "02. Interpolasi Polinom",
			Run:         interpolation.RunCLI,
		}, nil
	case 3:
		return Module{
			Choice:      3,
			DisplayName: "03. Interpolasi Spline",
			Run:         spline.RunCLI,
		}, nil
	default:
		return Module{}, fmt.Errorf("pilihan modul harus 1, 2, atau 3")
	}
}
