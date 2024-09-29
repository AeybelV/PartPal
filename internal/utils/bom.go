package utils

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/AeybelV/PartPal/internal/distributors"
)

// A Struct representing a ROW in a BOM

type BOMRow struct {
	PartNumber  string
	Description string
	Quantity    int
}

// This map will define alternative field names for each column we're interested in
var BOMFieldMapping = map[string]string{
	"Part Number": "PartNumber",
	"Component":   "PartNumber",
	"Part":        "PartNumber",
	"Description": "Description",
	"Desc":        "Description",
	"Desc.":       "Description",
	"Qty":         "Quantity",
	"Qty.":        "Quantity",
	"Quantity":    "Quantity",
}

// ReadBOM function to read from a CSV file and return a list of BOMRows
func ReadBOMCSV(filepath string) ([]BOMRow, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true

	// Read the header row to map column names
	header, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read header: %w", err)
	}

	// Create a map of column indices to unified field names
	columnMap := make(map[int]string)
	for i, column := range header {
		normalized := strings.TrimSpace(column)
		if unifiedField, ok := BOMFieldMapping[normalized]; ok {
			columnMap[i] = unifiedField
		}
	}

	var bomRows []BOMRow

	// Process each row and map the fields
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read record: %w", err)
		}

		row := BOMRow{}
		// Use reflection to set struct fields dynamically
		rowValue := reflect.ValueOf(&row).Elem()

		for i, value := range record {
			if fieldName, ok := columnMap[i]; ok {
				// Find the field by name
				field := rowValue.FieldByName(fieldName)
				if !field.IsValid() {
					continue
				}

				// Set the field value depending on its type
				switch field.Kind() {
				case reflect.String:
					field.SetString(value)
				case reflect.Int:
					intVal, err := strconv.Atoi(value)
					if err != nil {
						return nil, fmt.Errorf("failed to convert quantity to int: %w", err)
					}
					field.SetInt(int64(intVal))
				}
			}
		}
		bomRows = append(bomRows, row)
	}

	return bomRows, nil
}

// Wrapper function to read BOM and convert it to PartInfo array
func ReadBOM(filepath, filetype string) ([]distributors.PartInfo, error) {
	var bomRows []BOMRow
	var err error

	// Determine which ReadBOM function to call based on file type
	switch strings.ToLower(filetype) {
	case "csv":
		bomRows, err = ReadBOMCSV(filepath) // Calls the CSV reader function
		if err != nil {
			return nil, fmt.Errorf("failed to read BOM: %w", err)
		}
	default:
		return nil, fmt.Errorf("unsupported file type: %s", filetype)
	}

	// Convert BOMRow array to PartInfo array
	partInfos := make([]distributors.PartInfo, len(bomRows))
	for i, bom := range bomRows {
		partInfos[i] = distributors.PartInfo{
			PartNumber:  bom.PartNumber,
			Description: bom.Description,
			// Add more mappings if needed (e.g., from BOMRow to PartInfo)
		}
	}

	return partInfos, nil
}
