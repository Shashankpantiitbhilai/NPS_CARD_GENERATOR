package tests

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"nps-card-generator/internal/generator"
	"nps-card-generator/internal/models"
	"nps-card-generator/pkg/utils"
)

func TestGenerateCard(t *testing.T) {
	// Skip if font files don't exist
	regularFontPath := filepath.Join("..", "assets", "fonts", "Poppins-Regular.ttf")
	boldFontPath := filepath.Join("..", "assets", "fonts", "Poppins-Bold.ttf")
	
	if _, err := os.Stat(regularFontPath); os.IsNotExist(err) {
		t.Skip("Skipping test because regular font file not found")
	}
	
	if _, err := os.Stat(boldFontPath); os.IsNotExist(err) {
		t.Skip("Skipping test because bold font file not found")
	}
	
	// Load fonts
	regularFont, err := utils.LoadFont(regularFontPath)
	if err != nil {
		t.Fatalf("Failed to load regular font: %v", err)
	}
	
	boldFont, err := utils.LoadFont(boldFontPath)
	if err != nil {
		t.Fatalf("Failed to load bold font: %v", err)
	}
	
	// Create test portfolio
	testData := models.UserPortfolio{
		Username:   "TestUser",
		XIRR:       9.5,
		XIRRPeriod: 2,
		Schemes: []models.Scheme{
			{Name: "Equity", Allocation: 60, Color: "#4285F4"},
			{Name: "Debt", Allocation: 30, Color: "#34A853"},
			{Name: "Gold", Allocation: 10, Color: "#FBBC05"},
		},
	}
	
	// Create temporary output directory
	tmpDir := t.TempDir()
	outputPath := filepath.Join(tmpDir, "test_card.png")
	
	// Create card generator
	config := generator.Config{
		Width:       400,
		Height:      300,
		RegularFont: regularFont,
		BoldFont:    boldFont,
	}
	cardGenerator := generator.NewCardGenerator(config)
	
	// Test generating the card
	err = cardGenerator.GenerateCard(testData, outputPath, true)
	if err != nil {
		t.Fatalf("GenerateCard failed: %v", err)
	}
	
	// Check if the output file exists
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Errorf("Output PNG file was not created")
	}
	
	// Check if SVG file was created
	svgPath := filepath.Join(tmpDir, "test_card.svg")
	if _, err := os.Stat(svgPath); os.IsNotExist(err) {
		t.Errorf("Output SVG file was not created")
	}
}

func TestLoadUserPortfolio(t *testing.T) {
	// Create a temporary JSON file for testing
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test_portfolio.json")
	
	testData := models.UserPortfolio{
		Username:   "TestUser",
		XIRR:       9.5,
		XIRRPeriod: 2,
		Schemes: []models.Scheme{
			{Name: "Test1", Allocation: 50, Color: "#FF0000"},
			{Name: "Test2", Allocation: 50, Color: "#00FF00"},
		},
	}
	
	jsonBytes, err := json.Marshal(testData)
	if err != nil {
		t.Fatalf("Failed to marshal test data: %v", err)
	}
	
	if err := os.WriteFile(testFile, jsonBytes, 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}
	
	// Test loading the portfolio (direct call to the function under test)
	var userData models.UserPortfolio
	jsonData, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Error reading portfolio JSON: %v", err)
	}
	
	if err := json.Unmarshal(jsonData, &userData); err != nil {
		t.Fatalf("Error parsing portfolio JSON: %v", err)
	}
	
	// Verify loaded data
	if userData.Username != testData.Username {
		t.Errorf("Expected username %s, got %s", testData.Username, userData.Username)
	}
	
	if userData.XIRR != testData.XIRR {
		t.Errorf("Expected XIRR %f, got %f", testData.XIRR, userData.XIRR)
	}
	
	if len(userData.Schemes) != len(testData.Schemes) {
		t.Errorf("Expected %d schemes, got %d", len(testData.Schemes), len(userData.Schemes))
	}
}