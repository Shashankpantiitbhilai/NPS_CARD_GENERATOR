package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"nps-card-generator/internal/generator"
	"nps-card-generator/internal/models"
	"nps-card-generator/pkg/utils"
)

func main() {
	// Define command-line flags
	jsonPath := flag.String("json", "portfolio.json", "Path to the portfolio JSON file")
	outputPath := flag.String("output", "nps_card.png", "Output PNG file path")
	regularFontPath := flag.String("regular-font", "assets/fonts/Poppins-Regular.ttf", "Path to regular font TTF file")
	boldFontPath := flag.String("bold-font", "assets/fonts/Poppins-Bold.ttf", "Path to bold font TTF file")
	width := flag.Int("width", 400, "Width of the output image")
	height := flag.Int("height", 300, "Height of the output image")
	saveSVG := flag.Bool("save-svg", false, "Save SVG version for debugging")
	flag.Parse()

	// Load user portfolio data
	userData, err := loadUserPortfolio(*jsonPath)
	if err != nil {
		log.Fatalf("Failed to load user portfolio: %v", err)
	}

	// Load fonts
	regularFont, err := utils.LoadFont(*regularFontPath)
	if err != nil {
		log.Fatalf("Failed to load regular font: %v", err)
	}

	boldFont, err := utils.LoadFont(*boldFontPath)
	if err != nil {
		log.Fatalf("Failed to load bold font: %v", err)
	}

	// Create output directory if it doesn't exist
	if err := ensureOutputDir(*outputPath); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	// Create card generator with configuration
	config := generator.Config{
		Width:       *width,
		Height:      *height,
		RegularFont: regularFont,
		BoldFont:    boldFont,
	}
	cardGenerator := generator.NewCardGenerator(config)

	// Generate the card
	if err := cardGenerator.GenerateCard(userData, *outputPath, *saveSVG); err != nil {
		log.Fatalf("Failed to generate card: %v", err)
	}
}

// loadUserPortfolio loads a UserPortfolio from a JSON file
func loadUserPortfolio(filePath string) (models.UserPortfolio, error) {
	var userData models.UserPortfolio
	
	// Read portfolio data from JSON file
	jsonData, err := os.ReadFile(filePath)
	if err != nil {
		return userData, fmt.Errorf("error reading portfolio JSON: %v", err)
	}

	// Parse JSON into our struct
	if err := json.Unmarshal(jsonData, &userData); err != nil {
		return userData, fmt.Errorf("error parsing portfolio JSON: %v", err)
	}
	
	// Validate the portfolio data
	if err := userData.Validate(); err != nil {
		return userData, fmt.Errorf("invalid portfolio data: %v", err)
	}
	
	return userData, nil
}

// ensureOutputDir ensures the directory for the output file exists
func ensureOutputDir(outputPath string) error {
	dir := filepath.Dir(outputPath)
	if dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create output directory: %v", err)
		}
	}
	return nil
}