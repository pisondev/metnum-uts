package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

const (
	ColorBlue  = "\033[1;34m"
	ColorReset = "\033[0m"
)

// GenerateUniqueFilename mencari nama file report-n.md yang belum ada
func GenerateUniqueFilename() string {
	i := 1
	for {
		filename := fmt.Sprintf("report-%d.md", i)
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			return filename
		}
		i++
	}
}

func ParseInput(input string) float64 {
	if strings.Contains(input, "/") {
		parts := strings.Split(input, "/")
		if len(parts) == 2 {
			num, _ := strconv.ParseFloat(parts[0], 64)
			den, _ := strconv.ParseFloat(parts[1], 64)
			if den != 0 {
				return num / den
			}
		}
	}
	val, _ := strconv.ParseFloat(input, 64)
	return val
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
