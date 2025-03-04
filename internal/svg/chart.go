package svg

import (
	"fmt"
	"math"
	"strings"

	"nps-card-generator/internal/models"
)

// Chart handles the SVG generation for pie/donut charts
type Chart struct {
	Width  int
	Height int
}

// NewChart creates a new Chart with the specified dimensions
func NewChart(width, height int) *Chart {
	return &Chart{
		Width:  width,
		Height: height,
	}
}

// GenerateSVG creates an SVG string for the background card and donut arcs
func (c *Chart) GenerateSVG(userData models.UserPortfolio) string {
	return fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d">
    <!-- Background Card -->
    <rect width="%d" height="%d" fill="white" rx="10" ry="10"
          stroke="#e0e0e0" stroke-width="1"/>

    <!-- Donut Chart and Legend -->
    %s
</svg>`, c.Width, c.Height, c.Width, c.Height, c.Width, c.Height, c.generatePieChart(userData.Schemes))
}

// generatePieChart creates donut segments with white stroke boundaries and legend squares
func (c *Chart) generatePieChart(schemes []models.Scheme) string {
	// Donut center and radii
	cx, cy := float64(c.Width) * 0.25, float64(c.Height) * 0.67
	radius := float64(c.Height) * 0.23
	innerRadius := radius * 0.43

	var pieElements, legendElements strings.Builder
	startAngle := -math.Pi / 2

	for i, scheme := range schemes {
		sweepAngle := (float64(scheme.Allocation) / 100.0) * 2 * math.Pi
		endAngle := startAngle + sweepAngle

		// Outer arc points
		x1 := cx + radius*math.Cos(startAngle)
		y1 := cy + radius*math.Sin(startAngle)
		x2 := cx + radius*math.Cos(endAngle)
		y2 := cy + radius*math.Sin(endAngle)

		// Inner arc points
		x3 := cx + innerRadius*math.Cos(endAngle)
		y3 := cy + innerRadius*math.Sin(endAngle)
		x4 := cx + innerRadius*math.Cos(startAngle)
		y4 := cy + innerRadius*math.Sin(startAngle)

		// Large arc flag
		largeArcFlag := 0
		if sweepAngle > math.Pi {
			largeArcFlag = 1
		}

		// White stroke boundary to separate segments
		path := fmt.Sprintf(
			`<path fill-rule="evenodd" stroke="white" stroke-width="2"
				d="M %.1f,%.1f A %.1f,%.1f 0 %d 1 %.1f,%.1f
				   L %.1f,%.1f
				   A %.1f,%.1f 0 %d 0 %.1f,%.1f
				   Z" fill="%s"/>`,
			x1, y1, radius, radius, largeArcFlag, x2, y2,
			x3, y3, innerRadius, innerRadius, largeArcFlag, x4, y4,
			scheme.Color,
		)
		pieElements.WriteString(path + "\n")

		// Legend squares (positioned in right side of canvas)
		legendX := float64(c.Width) * 0.525
		legendY := float64(c.Height) * 0.567 + float64(i*25)
		rect := fmt.Sprintf(
			`<rect x="%.1f" y="%.1f" width="15" height="15" fill="%s"/>`,
			legendX, legendY, scheme.Color,
		)
		legendElements.WriteString(rect + "\n")

		startAngle = endAngle
	}

	return pieElements.String() + legendElements.String()
}