package generator

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"

	"nps-card-generator/internal/models"
	"nps-card-generator/internal/svg"
)

// CardGenerator handles the generation of portfolio cards
type CardGenerator struct {
	regularFont *truetype.Font
	boldFont    *truetype.Font
	width       int
	height      int
	svgChart    *svg.Chart
}

// Config contains the configuration for the card generator
type Config struct {
	Width         int
	Height        int
	RegularFont   *truetype.Font
	BoldFont      *truetype.Font
}

// NewCardGenerator creates a new CardGenerator with the specified configuration
func NewCardGenerator(config Config) *CardGenerator {
	return &CardGenerator{
		regularFont: config.RegularFont,
		boldFont:    config.BoldFont,
		width:       config.Width,
		height:      config.Height,
		svgChart:    svg.NewChart(config.Width, config.Height),
	}
}

// GenerateCard creates a PNG card from the provided user portfolio data
func (g *CardGenerator) GenerateCard(userData models.UserPortfolio, outputPath string, saveSVG bool) error {
	// Generate SVG data
	svgData := g.svgChart.GenerateSVG(userData)

	// Optionally save the SVG for debugging
	if saveSVG {
		svgPath := strings.TrimSuffix(outputPath, filepath.Ext(outputPath)) + ".svg"
		if err := os.WriteFile(svgPath, []byte(svgData), 0644); err != nil {
			return fmt.Errorf("saving SVG: %v", err)
		}
		fmt.Printf("SVG version saved for debugging: %s\n", svgPath)
	}

	// Parse the SVG with oksvg
	icon, err := oksvg.ReadIconStream(strings.NewReader(svgData))
	if err != nil {
		return fmt.Errorf("parsing SVG: %v", err)
	}
	icon.SetTarget(0, 0, float64(g.width), float64(g.height))

	// Create an RGBA image to draw onto
	img := image.NewRGBA(image.Rect(0, 0, g.width, g.height))
	scanner := rasterx.NewScannerGV(g.width, g.height, img, img.Bounds())
	dasher := rasterx.NewDasher(g.width, g.height, scanner)

	// Draw the donut chart and legend squares
	icon.Draw(dasher, 1.0)

	// Draw text on top of the image
	if err := g.drawLabels(img, userData); err != nil {
		return fmt.Errorf("drawing labels: %v", err)
	}

	// Save final PNG
	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("creating output file: %v", err)
	}
	defer outFile.Close()

	if err := png.Encode(outFile, img); err != nil {
		return fmt.Errorf("encoding PNG: %v", err)
	}
	
	fmt.Printf("Card generated successfully: %s\n", outputPath)
	return nil
}

// drawLabels draws all text elements onto the image
func (g *CardGenerator) drawLabels(img *image.RGBA, userData models.UserPortfolio) error {
	// Draw Title: Smaller, lighter gray (#9E9E9E)
	if err := g.drawTextLine(img, g.regularFont, 14, color.RGBA{0x9E, 0x9E, 0x9E, 0xFF},
		20, 35, fmt.Sprintf("%s's NPS Tier 1 Portfolio", userData.Username)); err != nil {
		return err
	}

	// Draw Subtitle: Larger, darker (#333333)
	if err := g.drawTextLine(img, g.regularFont, 20, color.RGBA{0x33, 0x33, 0x33, 0xFF},
		20, 70, "All Schemes"); err != nil {
		return err
	}

	// Draw XIRR line: Bold font with darker color
	if err := g.drawTextLine(img, g.boldFont, 14, color.RGBA{0x21, 0x21, 0x21, 0xFF},
		20, 100, fmt.Sprintf("XIRR : %.1f%% in last %d years", userData.XIRR, userData.XIRRPeriod)); err != nil {
		return err
	}

	// Draw Legend text: Medium size, dark (#333333)
	legendX := float64(g.width) * 0.525 + 20  // Position to the right of legend squares
	for i, scheme := range userData.Schemes {
		yPos := int(float64(g.height)*0.567) + i*25 + 12
		if err := g.drawTextLine(img, g.regularFont, 14, color.RGBA{0x33, 0x33, 0x33, 0xFF},
			int(legendX), yPos, fmt.Sprintf("%s : %d%%", scheme.Name, scheme.Allocation)); err != nil {
			return err
		}
	}

	return nil
}

// drawTextLine draws a single line of text onto img using the provided TrueType font
func (g *CardGenerator) drawTextLine(img *image.RGBA, fnt *truetype.Font, fontSize float64, col color.Color,
	x, y int, text string) error {
	c := freetype.NewContext()
	c.SetDPI(72)
	c.SetFont(fnt)
	c.SetFontSize(fontSize)
	c.SetClip(img.Bounds())
	c.SetDst(img)
	c.SetSrc(image.NewUniform(col))
	pt := freetype.Pt(x, y)
	_, err := c.DrawString(text, pt)
	return err
}