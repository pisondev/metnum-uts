package spline

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

// GenerateUniqueFilename mencari nama file report-n.md yang belum ada di folder results/<metode>.
func GenerateUniqueFilename(methodSlug string) (string, error) {
	methodSlug = strings.TrimSpace(methodSlug)
	if methodSlug == "" {
		return "", fmt.Errorf("slug metode tidak boleh kosong")
	}

	resultsDir := filepath.Join("results", methodSlug)

	if err := os.MkdirAll(resultsDir, 0755); err != nil {
		return "", err
	}

	i := 1
	for {
		filename := filepath.Join(resultsDir, fmt.Sprintf("report-%d.md", i))
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			return filename, nil
		}
		i++
	}
}

func ParseInput(input string) (float64, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return 0, fmt.Errorf("input kosong")
	}

	if strings.Contains(input, "/") {
		parts := strings.Split(input, "/")
		if len(parts) == 2 {
			num, err := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
			if err != nil {
				return 0, fmt.Errorf("pembilang tidak valid: %q", parts[0])
			}

			den, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
			if err != nil {
				return 0, fmt.Errorf("penyebut tidak valid: %q", parts[1])
			}

			if den != 0 {
				val := num / den
				if !IsFinite(val) {
					return 0, fmt.Errorf("hasil input tidak terhingga: %q", input)
				}

				return val, nil
			}

			return 0, fmt.Errorf("penyebut tidak boleh nol")
		}

		return 0, fmt.Errorf("format pecahan tidak valid: %q", input)
	}

	val, err := strconv.ParseFloat(input, 64)
	if err != nil {
		return 0, fmt.Errorf("angka tidak valid: %q", input)
	}
	if !IsFinite(val) {
		return 0, fmt.Errorf("angka harus bernilai hingga: %q", input)
	}

	return val, nil
}

func IsFinite(val float64) bool {
	return !math.IsNaN(val) && !math.IsInf(val, 0)
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
				FormatNum(points[i].X),
				i-1,
				FormatNum(points[i-1].X),
			)
		}
	}

	return nil
}

func FindSegment(points []Point, targetX float64) int {
	for i := 0; i < len(points)-1; i++ {
		if targetX >= points[i].X && targetX <= points[i+1].X {
			return i
		}
	}

	return -1
}

func NearlyEqual(a, b float64) bool {
	const eps = 1e-9
	diff := math.Abs(a - b)
	scale := math.Max(1.0, math.Max(math.Abs(a), math.Abs(b)))
	return diff <= eps*scale
}

func ResolveTargetSegment(points []Point, targetX float64) (int, string) {
	for i, point := range points {
		if NearlyEqual(targetX, point.X) {
			switch {
			case i == 0:
				return 0, fmt.Sprintf(
					"x = %s tepat pada titik data awal x_%d, sehingga digunakan segmen 0 pada interval [%s, %s].",
					FormatNum(targetX),
					i,
					FormatNum(points[0].X),
					FormatNum(points[1].X),
				)
			case i == len(points)-1:
				return len(points) - 2, fmt.Sprintf(
					"x = %s tepat pada titik data akhir x_%d, sehingga digunakan segmen %d pada interval [%s, %s].",
					FormatNum(targetX),
					i,
					len(points)-2,
					FormatNum(points[len(points)-2].X),
					FormatNum(points[len(points)-1].X),
				)
			default:
				return i - 1, fmt.Sprintf(
					"x = %s tepat pada titik data x_%d. Karena berada di sambungan dua segmen, digunakan segmen %d sebagai konvensi kiri.",
					FormatNum(targetX),
					i,
					i-1,
				)
			}
		}
	}

	segIdx := FindSegment(points, targetX)
	if segIdx == -1 {
		return -1, ""
	}

	return segIdx, fmt.Sprintf(
		"x = %s berada pada interval [%s, %s], sehingga digunakan segmen %d.",
		FormatNum(targetX),
		FormatNum(points[segIdx].X),
		FormatNum(points[segIdx+1].X),
		segIdx,
	)
}

func FormatNum(val float64) string {
	if math.Abs(val) < 1e-10 {
		return "0"
	}
	frac := FloatToFraction(val)
	if frac == fmt.Sprintf("%g", val) {
		return frac
	}
	return fmt.Sprintf("%s (≈ %.4f)", frac, val)
}

func FormatFractionOnly(val float64) string {
	if math.Abs(val) < 1e-10 {
		return "0"
	}
	return FloatToFraction(val)
}

func FloatToFraction(val float64) string {
	sign := 1.0
	if val < 0 {
		sign = -1.0
		val = -val
	}
	tolerance := 1.0e-6
	h1, h2, k1, k2 := 1, 0, 0, 1
	b := val
	for i := 0; i < 50; i++ {
		a := math.Floor(b)
		h, k := int(a)*h1+h2, int(a)*k1+k2
		h2, h1, k2, k1 = h1, h, k1, k
		if math.Abs(val-float64(h)/float64(k)) < tolerance {
			break
		}
		if b-a < 1e-10 {
			break
		}
		b = 1.0 / (b - a)
	}
	num, den := int(sign)*h1, k1
	if den == 1 {
		return fmt.Sprintf("%d", num)
	}
	return fmt.Sprintf("%d/%d", num, den)
}

func PrintStageTitle(title string) {
	fmt.Printf("\n%s%s%s\n", ColorBlue, title, ColorReset)
}
