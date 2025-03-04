package utils

import (
	"fmt"
	"os"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
)

// LoadFont loads a TrueType font from a file path
func LoadFont(fontPath string) (*truetype.Font, error) {
	fontBytes, err := os.ReadFile(fontPath)
	if err != nil {
		return nil, fmt.Errorf("cannot load font from %s: %v", fontPath, err)
	}
	
	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		return nil, fmt.Errorf("cannot parse font from %s: %v", fontPath, err)
	}
	
	return font, nil
}