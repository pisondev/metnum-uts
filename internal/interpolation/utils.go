package interpolation

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

func ValidatePoints(points []Point) error {
	if len(points) < 2 {
		return fmt.Errorf("minimal diperlukan 2 titik data")
	}

	for i, point := range points {
		if !IsFinite(point.X) || !IsFinite(point.Y) {
			return fmt.Errorf("titik ke-%d mengandung nilai yang tidak hingga", i)
		}
	}

	for i := 1; i < len(points); i++ {
		if points[i].X <= points[i-1].X {
			return fmt.Errorf(
				"nilai x harus naik ketat: x_%d=%s tidak lebih besar dari x_%d=%s",
				i,
				FormatReal(points[i].X),
				i-1,
				FormatReal(points[i-1].X),
			)
		}
	}

	return nil
}

func IsFinite(value float64) bool {
	return !math.IsNaN(value) && !math.IsInf(value, 0)
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

func WriteReportInputSection(report *strings.Builder, points []Point, targets []float64) {
	report.WriteString("## Data Input\n")
	for i, point := range points {
		report.WriteString(fmt.Sprintf("- Titik %d: (%s, %s)\n", i, FormatReal(point.X), FormatReal(point.Y)))
	}

	if len(targets) == 0 {
		report.WriteString("- Tidak ada target x yang diminta.\n")
		return
	}

	report.WriteString("\n## Target Evaluasi\n")
	for i, target := range targets {
		report.WriteString(fmt.Sprintf("- Target %d: x = %s\n", i+1, FormatReal(target)))
	}
}
