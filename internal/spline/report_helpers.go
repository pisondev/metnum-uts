package spline

import (
	"fmt"
	"strings"
)

func WriteReportInputSection(report *strings.Builder, points []Point, targets []float64) {
	report.WriteString("## Data Input\n")
	for i, p := range points {
		report.WriteString(fmt.Sprintf("- Titik %d: (%v, %v)\n", i, FormatNum(p.X), FormatNum(p.Y)))
	}

	if len(targets) == 0 {
		report.WriteString("- Tidak ada target x yang diminta.\n")
		return
	}

	report.WriteString("\n## Target Evaluasi\n")
	for i, target := range targets {
		report.WriteString(fmt.Sprintf("- Target %d: $x = %s$\n", i+1, FormatNum(target)))
	}
}
