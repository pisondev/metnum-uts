package rootfinding

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	ColorBlue  = "\033[1;34m"
	ColorReset = "\033[0m"
)

func GenerateUniqueFilename(methodSlug string) (string, error) {
	methodSlug = strings.TrimSpace(methodSlug)
	if methodSlug == "" {
		return "", fmt.Errorf("slug metode tidak boleh kosong")
	}

	resultsDir := filepath.Join("results", methodSlug)
	if err := os.MkdirAll(resultsDir, 0755); err != nil {
		return "", err
	}

	for i := 1; ; i++ {
		filename := filepath.Join(resultsDir, fmt.Sprintf("report-%d.md", i))
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			return filename, nil
		}
	}
}

func ParseNumberInput(input string) (float64, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return 0, fmt.Errorf("input kosong")
	}

	if strings.Contains(input, "/") {
		parts := strings.Split(input, "/")
		if len(parts) != 2 {
			return 0, fmt.Errorf("format pecahan tidak valid: %q", input)
		}

		num, err := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
		if err != nil {
			return 0, fmt.Errorf("pembilang tidak valid: %q", parts[0])
		}

		den, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
		if err != nil {
			return 0, fmt.Errorf("penyebut tidak valid: %q", parts[1])
		}

		if den == 0 {
			return 0, fmt.Errorf("penyebut tidak boleh nol")
		}

		value := num / den
		if !IsFinite(value) {
			return 0, fmt.Errorf("hasil input tidak terhingga: %q", input)
		}

		return value, nil
	}

	value, err := strconv.ParseFloat(input, 64)
	if err != nil {
		return 0, fmt.Errorf("angka tidak valid: %q", input)
	}
	if !IsFinite(value) {
		return 0, fmt.Errorf("angka harus bernilai hingga: %q", input)
	}

	return value, nil
}

func IsFinite(value float64) bool {
	return !math.IsNaN(value) && !math.IsInf(value, 0)
}

func NearlyEqual(a, b float64) bool {
	const eps = 1e-12
	diff := math.Abs(a - b)
	scale := math.Max(1, math.Max(math.Abs(a), math.Abs(b)))
	return diff <= eps*scale
}

func FormatReal(value float64) string {
	if !IsFinite(value) {
		return "tak hingga"
	}
	if math.Abs(value) < 1e-12 {
		return "0"
	}
	return strconv.FormatFloat(value, 'g', 12, 64)
}
