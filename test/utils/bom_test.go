package utils

import (
	"os"
	"reflect"
	"testing"

	"github.com/AeybelV/PartPal/internal/distributors"
	"github.com/AeybelV/PartPal/internal/utils"
)

// Sample CSV data for testing
const sampleCSV = `Part Number,Description,Qty
1234-5678,Resistor,10
9876-5432,Capacitor,5
`

func createTempCSV(t *testing.T, content string) string {
	file, err := os.CreateTemp("", "bom_test_*.csv")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	if _, err := file.WriteString(content); err != nil {
		file.Close() // Ensure file is closed before removing it
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	file.Close()
	return file.Name()
}

func TestBOM_BOMReadCSV(t *testing.T) {
	// Create a temporary CSV file
	tempFile := createTempCSV(t, sampleCSV)
	defer os.Remove(tempFile) // Clean up after the test

	// Call the ReadBOM function
	result, err := utils.ReadBOMCSV(tempFile)
	if err != nil {
		t.Fatalf("ReadBOM returned an error: %v", err)
	}

	// Define the expected result
	expected := []utils.BOMRow{
		{PartNumber: "1234-5678", Description: "Resistor", Quantity: 10},
		{PartNumber: "9876-5432", Description: "Capacitor", Quantity: 5},
	}

	// Check if the result matches the expected output
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("ReadBOM result mismatch. Got: %v, Expected: %v", result, expected)
	}
}

// Test for ReadBOM function
func TestReadBOM(t *testing.T) {
	// Create a temporary CSV file
	tempFile := createTempCSV(t, sampleCSV)
	defer os.Remove(tempFile) // Clean up after the test

	// Call the wrapper function ReadBOM
	result, err := utils.ReadBOM(tempFile, "csv")
	if err != nil {
		t.Fatalf("ReadBOM returned an error: %v", err)
	}

	// Define the expected result
	expected := []distributors.PartInfo{
		{PartNumber: "1234-5678", Description: "Resistor"},
		{PartNumber: "9876-5432", Description: "Capacitor"},
	}

	// Check if the result matches the expected output
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("ReadBOM result mismatch. Got: %v, Expected: %v", result, expected)
	}
}

// Additional test cases for edge cases
func TestReadBOM_EmptyFile(t *testing.T) {
	// Create an empty CSV file
	tempFile := createTempCSV(t, "")
	defer os.Remove(tempFile)

	// Call the wrapper function ReadBOM
	_, err := utils.ReadBOM(tempFile, "csv")
	if err == nil {
		t.Error("Expected error for empty file, but got none")
	}
}

func TestReadBOM_InvalidFileType(t *testing.T) {
	// Create a temporary CSV file
	tempFile := createTempCSV(t, sampleCSV)
	defer os.Remove(tempFile)

	// Call the wrapper function with an invalid file type
	_, err := utils.ReadBOM(tempFile, "invalid")
	if err == nil {
		t.Error("Expected error for unsupported file type, but got none")
	}
}
